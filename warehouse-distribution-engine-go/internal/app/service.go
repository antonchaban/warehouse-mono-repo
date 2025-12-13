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
	GetWarehouseIDBySupplyID(ctx context.Context, supplyID int64) (int64, error)
}

// ResultSender defines the contract for sending the calculation result.
type ResultSender interface {
	SendPlan(ctx context.Context, plan algorithm.DistributionPlan, sourceID int64) error
}

// Service is the main orchestrator of the business logic.
// It connects the Infrastructure (Adapter) with the Core Domain (Algorithm).
type Service struct {
	provider StateProvider
	algo     algorithm.Packer
	sender   ResultSender
	logger   *slog.Logger
}

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

	// 1. Get warehouse ID from database
	sourceID, err := s.provider.GetWarehouseIDBySupplyID(ctx, supplyID)
	if err != nil {
		return fmt.Errorf("failed to determine source warehouse: %w", err)
	}
	s.logger.Info("identified source warehouse", "source_id", sourceID)

	// 2. Get world state
	warehouses, err := s.provider.FetchWorldState(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch world state: %w", err)
	}

	// 3. Get pending items
	items, err := s.provider.FetchPendingItems(ctx, supplyID)
	if err != nil {
		return fmt.Errorf("failed to fetch pending items: %w", err)
	}

	s.logger.Info("data loaded", "warehouses", len(warehouses), "items", len(items))

	// 4. Calculate distribution
	plan := s.algo.Distribute(requestID, warehouses, items)

	// 5. Send result
	if err := s.sender.SendPlan(ctx, plan, sourceID); err != nil {
		return fmt.Errorf("failed to send distribution plan: %w", err)
	}

	s.logger.Info("plan sent", "request_id", requestID)
	return nil
}
