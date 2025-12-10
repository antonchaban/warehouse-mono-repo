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

	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()
	logger.Info("connected to database")

	grpcClient, err := grpc.NewClient(cfg.JavaServiceAddr)
	if err != nil {
		logger.Error("failed to create grpc client", "error", err)
	} else {
		defer func(grpcClient *grpc.Client) {
			err := grpcClient.Close()
			if err != nil {
				logger.Error("failed to close grpc client", "error", err)
			}
		}(grpcClient)
		logger.Info("connected to java service", "addr", cfg.JavaServiceAddr)
	}

	repo := postgres.NewRepository(dbPool)
	algo := algorithm.NewWFDAlgorithm()

	var sender app.ResultSender

	if cfg.MockOutput {
		logger.Info("using MOCK sender (results will be logged, not sent via network)")
		sender = &noopSender{l: logger}
	} else {
		grpcClient, err := grpc.NewClient(cfg.JavaServiceAddr)
		if err != nil {
			logger.Error("failed to create grpc client", "error", err)
		}

		if grpcClient != nil {
			sender = grpcClient
		} else {
			logger.Warn("gRPC client creation failed, falling back to noopSender")
			sender = &noopSender{l: logger}
		}
	}

	distributionService := app.NewService(repo, algo, sender, logger)

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

type noopSender struct {
	l *slog.Logger
}

func (n *noopSender) SendPlan(ctx context.Context, plan algorithm.DistributionPlan) error {
	n.l.Info("mock sender: plan calculated", "moves", len(plan.Moves))
	return nil
}
