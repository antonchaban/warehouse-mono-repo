package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/antonchaban/warehouse-distribution-engine-go/config"
	"github.com/antonchaban/warehouse-distribution-engine-go/internal/adapter/grpc"
	"github.com/antonchaban/warehouse-distribution-engine-go/internal/adapter/postgres"
	"github.com/antonchaban/warehouse-distribution-engine-go/internal/adapter/rabbitmq"
	"github.com/antonchaban/warehouse-distribution-engine-go/internal/algorithm"
	"github.com/antonchaban/warehouse-distribution-engine-go/internal/app"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("initializing distribution engine...")

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Database Connection
	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()
	logger.Info("connected to database")

	// 2. Init Dependencies
	repo := postgres.NewRepository(dbPool)
	algo := algorithm.NewWFDAlgorithm()
	var sender app.ResultSender

	// 3. gRPC Client Logic
	if cfg.MockOutput {
		logger.Info("using MOCK sender (results will be logged, not sent via network)")
		sender = &noopSender{l: logger}
	} else {
		// Створюємо клієнт тільки якщо НЕ мок
		grpcClient, err := grpc.NewClient(cfg.JavaServiceAddr)
		if err != nil {
			logger.Error("failed to create grpc client, falling back to mock", "error", err)
			sender = &noopSender{l: logger}
		} else {
			// Закриваємо з'єднання при виході
			defer func() {
				if err := grpcClient.Close(); err != nil {
					logger.Error("failed to close grpc client", "error", err)
				}
			}()
			logger.Info("connected to java service", "addr", cfg.JavaServiceAddr)
			sender = grpcClient
		}
	}

	// 4. Create Service
	distributionService := app.NewService(repo, algo, sender, logger)

	// 5. RabbitMQ Consumer
	consumer := rabbitmq.NewConsumer(
		cfg.RabbitMQURL,
		cfg.QueueName,
		distributionService,
		logger,
	)

	go func() {
		logger.Info("starting rabbitmq consumer", "queue", cfg.QueueName)
		if err := consumer.Start(ctx); err != nil {
			logger.Error("consumer stopped", "error", err)
			cancel()
		}
	}()

	logger.Info("service is running")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Info("shutdown signal received")
	cancel()
	time.Sleep(1 * time.Second)
	logger.Info("service stopped")
}

// --- MOCK IMPLEMENTATION ---

type noopSender struct {
	l *slog.Logger
}

func (n *noopSender) SendPlan(ctx context.Context, plan algorithm.DistributionPlan, sourceID int64, supplyID int64) error {
	n.l.Info("mock sender: plan calculated",
		"moves", len(plan.Moves),
		"source_id", sourceID,
		"supply_id", supplyID,
	)
	return nil
}
