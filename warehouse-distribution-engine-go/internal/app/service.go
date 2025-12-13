package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/antonchaban/warehouse-distribution-engine-go/internal/algorithm"
)

// StateProvider defines the contract for retrieving the current world state.
// This interface allows us to mock the data source (database) in tests.
type StateProvider interface {
	FetchWorldState(ctx context.Context) ([]*algorithm.Warehouse, error)
	FetchPendingItems(ctx context.Context, supplyID int64) ([]algorithm.Product, error)
}

// ResultSender defines the contract for sending the calculation result.
// Usually implemented by a gRPC client or a message queue producer.
type ResultSender interface {
	SendPlan(ctx context.Context, plan algorithm.DistributionPlan) error
}

// Service is the main orchestrator of the business logic.
// It connects the Infrastructure (Adapter) with the Core Domain (Algorithm).
type Service struct {
	provider StateProvider
	algo     algorithm.Packer
	sender   ResultSender
	logger   *slog.Logger
}

// NewService creates a new instance of the application service with required dependencies.
func NewService(
	p StateProvider,
	a algorithm.Packer,
	s ResultSender,
	l *slog.Logger,
) *Service {
	return &Service{
		provider: p,
		algo:     a,
		sender:   s,
		logger:   l,
	}
}

func (s *Service) CalculateDistribution(ctx context.Context, requestID string, supplyID int64) error {
	s.logger.Info("starting distribution calculation",
		"request_id", requestID,
		"supply_id", supplyID,
	)

	warehouses, err := s.provider.FetchWorldState(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch world state: %w", err)
	}

	items, err := s.provider.FetchPendingItems(ctx, supplyID)
	if err != nil {
		return fmt.Errorf("failed to fetch pending items: %w", err)
	}

	s.logger.Info("data loaded",
		"warehouses_count", len(warehouses),
		"items_count", len(items))

	plan := s.algo.Distribute(warehouses, items)

	s.logger.Info("calculation completed",
		"moves_generated", len(plan.Moves),
		"unallocated_count", len(plan.UnallocatedItems))

	if err := s.sender.SendPlan(ctx, plan); err != nil {
		return fmt.Errorf("failed to send distribution plan: %w", err)
	}

	s.logger.Info("distribution plan sent successfully", "request_id", requestID)
	return nil
}
