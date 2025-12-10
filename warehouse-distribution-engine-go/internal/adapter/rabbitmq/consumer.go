package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// DistributionService defines the contract we need from the app layer.
type DistributionService interface {
	CalculateDistribution(ctx context.Context, requestID string, supplyID int64) error
}

// Consumer handles RabbitMQ connection and message processing.
type Consumer struct {
	connURL string
	queue   string
	service DistributionService
	logger  *slog.Logger
}

func NewConsumer(url, queueName string, svc DistributionService, l *slog.Logger) *Consumer {
	return &Consumer{
		connURL: url,
		queue:   queueName,
		service: svc,
		logger:  l,
	}
}

// MessageBody updated according to Spec v3: includes audit info about who initiated the calculation.
type MessageBody struct {
	RequestID           string `json:"request_id"`
	SupplyID            int64  `json:"supply_id"`
	SourceWarehouseID   int64  `json:"source_warehouse_id"`
	InitiatedByUserId   int64  `json:"initiated_by_user_id"`
	InitiatedByUsername string `json:"initiated_by_username"`
}

func (c *Consumer) Start(ctx context.Context) error {
	c.logger.Info("connecting to rabbitmq", "url", c.connURL)

	var conn *amqp.Connection
	var err error

	// RETRY LOGIC
	for i := 0; i < 15; i++ {
		conn, err = amqp.Dial(c.connURL)
		if err == nil {
			c.logger.Info("successfully connected to rabbitmq")
			break
		}

		c.logger.Warn("failed to connect to rabbitmq, retrying in 2s...",
			"attempt", i+1,
			"error", err)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(2 * time.Second):
			continue
		}
	}

	if err != nil {
		return fmt.Errorf("failed to connect to rabbitmq after retries: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	// Declare Queue (Idempotent)
	_, err = ch.QueueDeclare(
		c.queue, // name
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// 3. Start Consuming
	msgs, err := ch.Consume(
		c.queue, // queue
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	c.logger.Info("waiting for calculation requests", "queue", c.queue)

	// 4. Message Loop
	for {
		select {
		case <-ctx.Done():
			c.logger.Info("stopping consumer...")
			return nil
		case d, ok := <-msgs:
			if !ok {
				return fmt.Errorf("rabbitmq channel closed")
			}
			c.handleMessage(ctx, d)
		}
	}
}

func (c *Consumer) handleMessage(ctx context.Context, d amqp.Delivery) {
	var body MessageBody
	if err := json.Unmarshal(d.Body, &body); err != nil {
		c.logger.Error("critical error: failed to unmarshal message", "error", err)
		_ = d.Nack(false, false)
		return
	}

	// AUDIT LOGGING
	c.logger.Info("received calculation request",
		"request_id", body.RequestID,
		"user_id", body.InitiatedByUserId,
		"username", body.InitiatedByUsername,
		"username", body.InitiatedByUsername, // Log WHO did it
		"attempt", d.Redelivered,
	)

	calcCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	err := c.service.CalculateDistribution(calcCtx, body.RequestID, body.SupplyID)

	if err != nil {
		c.logger.Error("calculation failed", "request_id", body.RequestID, "error", err)
		time.Sleep(5 * time.Second)
		_ = d.Nack(false, true)
	} else {
		if err := d.Ack(false); err != nil {
			c.logger.Error("failed to ack message", "error", err)
		}
	}
}
