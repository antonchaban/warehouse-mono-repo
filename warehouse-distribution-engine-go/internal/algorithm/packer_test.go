package algorithm_test

import (
	"testing"

	"github.com/antonchaban/warehouse-distribution-engine-go/internal/algorithm"
)

func TestDistribute_SingleItemSingleWarehouse_AllocatesSuccessfully(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-1", TotalCapacityM3: 100.0, CurrentStockM3: 0, IncomingM3: 0},
	}
	items := []algorithm.Product{
		{ID: "ITEM-1", VolumeM3: 10.0, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 1 {
		t.Fatalf("expected 1 move, got %d", len(plan.Moves))
	}
	if len(plan.UnallocatedItems) != 0 {
		t.Fatalf("expected 0 unallocated items, got %d", len(plan.UnallocatedItems))
	}
	if plan.Moves[0].ProductID != "ITEM-1" || plan.Moves[0].WarehouseID != "WH-1" {
		t.Errorf("unexpected move: %+v", plan.Moves[0])
	}
}

func TestDistribute_NoWarehouses_AllItemsUnallocated(t *testing.T) {
	warehouses := []*algorithm.Warehouse{}
	items := []algorithm.Product{
		{ID: "ITEM-1", VolumeM3: 10.0, Priority: 1},
		{ID: "ITEM-2", VolumeM3: 20.0, Priority: 2},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 0 {
		t.Fatalf("expected 0 moves, got %d", len(plan.Moves))
	}
	if len(plan.UnallocatedItems) != 2 {
		t.Fatalf("expected 2 unallocated items, got %d", len(plan.UnallocatedItems))
	}
}

func TestDistribute_NoItems_EmptyPlan(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-1", TotalCapacityM3: 100.0},
	}
	items := []algorithm.Product{}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 0 {
		t.Fatalf("expected 0 moves, got %d", len(plan.Moves))
	}
	if len(plan.UnallocatedItems) != 0 {
		t.Fatalf("expected 0 unallocated items, got %d", len(plan.UnallocatedItems))
	}
}

func TestDistribute_ItemTooLarge_RemainsUnallocated(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-1", TotalCapacityM3: 50.0},
		{ID: "WH-2", TotalCapacityM3: 30.0},
	}
	items := []algorithm.Product{
		{ID: "ITEM-1", VolumeM3: 100.0, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 0 {
		t.Fatalf("expected 0 moves, got %d", len(plan.Moves))
	}
	if len(plan.UnallocatedItems) != 1 {
		t.Fatalf("expected 1 unallocated item, got %d", len(plan.UnallocatedItems))
	}
	if plan.UnallocatedItems[0].ID != "ITEM-1" {
		t.Errorf("expected ITEM-1 to be unallocated, got %s", plan.UnallocatedItems[0].ID)
	}
}

func TestDistribute_DistributesAcrossMultipleWarehouses(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-1", TotalCapacityM3: 100.0},
		{ID: "WH-2", TotalCapacityM3: 100.0},
	}
	items := []algorithm.Product{
		{ID: "ITEM-1", VolumeM3: 60.0, Priority: 1},
		{ID: "ITEM-2", VolumeM3: 60.0, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 2 {
		t.Fatalf("expected 2 moves, got %d", len(plan.Moves))
	}
	if len(plan.UnallocatedItems) != 0 {
		t.Fatalf("expected 0 unallocated items, got %d", len(plan.UnallocatedItems))
	}

	warehouseIDs := make(map[string]bool)
	for _, m := range plan.Moves {
		warehouseIDs[m.WarehouseID] = true
	}
	if len(warehouseIDs) != 2 {
		t.Errorf("expected items to be distributed across 2 warehouses, but got %d", len(warehouseIDs))
	}
}

func TestDistribute_LargerItemsAllocatedFirst(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-1", TotalCapacityM3: 100.0},
	}
	items := []algorithm.Product{
		{ID: "SMALL", VolumeM3: 10.0, Priority: 1},
		{ID: "LARGE", VolumeM3: 50.0, Priority: 1},
		{ID: "MEDIUM", VolumeM3: 30.0, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 3 {
		t.Fatalf("expected 3 moves, got %d", len(plan.Moves))
	}
	if plan.Moves[0].ProductID != "LARGE" {
		t.Errorf("expected first move to be LARGE, got %s", plan.Moves[0].ProductID)
	}
	if plan.Moves[1].ProductID != "MEDIUM" {
		t.Errorf("expected second move to be MEDIUM, got %s", plan.Moves[1].ProductID)
	}
	if plan.Moves[2].ProductID != "SMALL" {
		t.Errorf("expected third move to be SMALL, got %s", plan.Moves[2].ProductID)
	}
}

func TestDistribute_SameVolume_HigherPriorityFirst(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-1", TotalCapacityM3: 100.0},
	}
	items := []algorithm.Product{
		{ID: "LOW-PRIO", VolumeM3: 20.0, Priority: 1},
		{ID: "HIGH-PRIO", VolumeM3: 20.0, Priority: 10},
		{ID: "MED-PRIO", VolumeM3: 20.0, Priority: 5},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 3 {
		t.Fatalf("expected 3 moves, got %d", len(plan.Moves))
	}
	if plan.Moves[0].ProductID != "HIGH-PRIO" {
		t.Errorf("expected first move to be HIGH-PRIO, got %s", plan.Moves[0].ProductID)
	}
	if plan.Moves[1].ProductID != "MED-PRIO" {
		t.Errorf("expected second move to be MED-PRIO, got %s", plan.Moves[1].ProductID)
	}
	if plan.Moves[2].ProductID != "LOW-PRIO" {
		t.Errorf("expected third move to be LOW-PRIO, got %s", plan.Moves[2].ProductID)
	}
}

func TestDistribute_PrefersLeastUtilizedWarehouse(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-FULL", TotalCapacityM3: 100.0, CurrentStockM3: 80.0},
		{ID: "WH-EMPTY", TotalCapacityM3: 100.0, CurrentStockM3: 0},
	}
	items := []algorithm.Product{
		{ID: "ITEM-1", VolumeM3: 10.0, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 1 {
		t.Fatalf("expected 1 move, got %d", len(plan.Moves))
	}
	if plan.Moves[0].WarehouseID != "WH-EMPTY" {
		t.Errorf("expected item to go to WH-EMPTY, got %s", plan.Moves[0].WarehouseID)
	}
}

func TestDistribute_AccountsForIncomingStock(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-1", TotalCapacityM3: 100.0, CurrentStockM3: 0, IncomingM3: 95.0},
	}
	items := []algorithm.Product{
		{ID: "ITEM-1", VolumeM3: 10.0, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 0 {
		t.Fatalf("expected 0 moves, got %d", len(plan.Moves))
	}
	if len(plan.UnallocatedItems) != 1 {
		t.Fatalf("expected 1 unallocated item, got %d", len(plan.UnallocatedItems))
	}
}

func TestDistribute_FallbackToAlternativeWarehouse(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-SMALL", TotalCapacityM3: 10.0, CurrentStockM3: 0},
		{ID: "WH-LARGE", TotalCapacityM3: 100.0, CurrentStockM3: 50.0},
	}
	items := []algorithm.Product{
		{ID: "ITEM-1", VolumeM3: 40.0, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 1 {
		t.Fatalf("expected 1 move, got %d", len(plan.Moves))
	}
	if plan.Moves[0].WarehouseID != "WH-LARGE" {
		t.Errorf("expected fallback to WH-LARGE, got %s", plan.Moves[0].WarehouseID)
	}
}

func TestDistribute_ExactFit_AllocatesSuccessfully(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-1", TotalCapacityM3: 50.0, CurrentStockM3: 0},
	}
	items := []algorithm.Product{
		{ID: "ITEM-1", VolumeM3: 50.0, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 1 {
		t.Fatalf("expected 1 move, got %d", len(plan.Moves))
	}
	if len(plan.UnallocatedItems) != 0 {
		t.Fatalf("expected 0 unallocated items, got %d", len(plan.UnallocatedItems))
	}
}

func TestDistribute_FullWarehouse_SkipsIt(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-FULL", TotalCapacityM3: 100.0, CurrentStockM3: 100.0},
		{ID: "WH-AVAILABLE", TotalCapacityM3: 100.0, CurrentStockM3: 0},
	}
	items := []algorithm.Product{
		{ID: "ITEM-1", VolumeM3: 10.0, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 1 {
		t.Fatalf("expected 1 move, got %d", len(plan.Moves))
	}
	if plan.Moves[0].WarehouseID != "WH-AVAILABLE" {
		t.Errorf("expected item in WH-AVAILABLE, got %s", plan.Moves[0].WarehouseID)
	}
}

func TestDistribute_ZeroCapacityWarehouse_TreatedAsFull(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-ZERO", TotalCapacityM3: 0},
		{ID: "WH-NORMAL", TotalCapacityM3: 100.0},
	}
	items := []algorithm.Product{
		{ID: "ITEM-1", VolumeM3: 10.0, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 1 {
		t.Fatalf("expected 1 move, got %d", len(plan.Moves))
	}
	if plan.Moves[0].WarehouseID != "WH-NORMAL" {
		t.Errorf("expected item in WH-NORMAL, got %s", plan.Moves[0].WarehouseID)
	}
}

func TestDistribute_PartialAllocation_SomeItemsUnallocated(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-1", TotalCapacityM3: 50.0},
	}
	items := []algorithm.Product{
		{ID: "ITEM-1", VolumeM3: 30.0, Priority: 1},
		{ID: "ITEM-2", VolumeM3: 30.0, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 1 {
		t.Fatalf("expected 1 move, got %d", len(plan.Moves))
	}
	if len(plan.UnallocatedItems) != 1 {
		t.Fatalf("expected 1 unallocated item, got %d", len(plan.UnallocatedItems))
	}
}

func TestDistribute_ManySmallItems_AllAllocated(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-1", TotalCapacityM3: 100.0},
	}
	items := make([]algorithm.Product, 100)
	for i := 0; i < 100; i++ {
		items[i] = algorithm.Product{ID: "ITEM", VolumeM3: 1.0, Priority: 1}
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 100 {
		t.Fatalf("expected 100 moves, got %d", len(plan.Moves))
	}
	if len(plan.UnallocatedItems) != 0 {
		t.Fatalf("expected 0 unallocated items, got %d", len(plan.UnallocatedItems))
	}
}

func TestDistribute_BalancesLoadAcrossWarehouses(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-1", TotalCapacityM3: 100.0},
		{ID: "WH-2", TotalCapacityM3: 100.0},
		{ID: "WH-3", TotalCapacityM3: 100.0},
	}
	items := make([]algorithm.Product, 30)
	for i := 0; i < 30; i++ {
		items[i] = algorithm.Product{ID: "ITEM", VolumeM3: 10.0, Priority: 1}
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 30 {
		t.Fatalf("expected 30 moves, got %d", len(plan.Moves))
	}

	warehouseLoad := make(map[string]float64)
	for _, m := range plan.Moves {
		warehouseLoad[m.WarehouseID] += m.VolumeM3
	}

	for whID, load := range warehouseLoad {
		if load != 100.0 {
			t.Errorf("expected warehouse %s to have load 100.0, got %f", whID, load)
		}
	}
}

func TestDistribute_MoveContainsCorrectVolume(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-1", TotalCapacityM3: 100.0},
	}
	items := []algorithm.Product{
		{ID: "ITEM-1", VolumeM3: 25.5, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 1 {
		t.Fatalf("expected 1 move, got %d", len(plan.Moves))
	}
	if plan.Moves[0].VolumeM3 != 25.5 {
		t.Errorf("expected VolumeM3 to be 25.5, got %f", plan.Moves[0].VolumeM3)
	}
}

func TestDistribute_ZeroVolumeItem_AllocatesSuccessfully(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-1", TotalCapacityM3: 100.0},
	}
	items := []algorithm.Product{
		{ID: "ITEM-1", VolumeM3: 0.0, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 1 {
		t.Fatalf("expected 1 move, got %d", len(plan.Moves))
	}
}

func TestDistribute_NegativePriority_StillSortsCorrectly(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-1", TotalCapacityM3: 100.0},
	}
	items := []algorithm.Product{
		{ID: "NEG", VolumeM3: 20.0, Priority: -5},
		{ID: "ZERO", VolumeM3: 20.0, Priority: 0},
		{ID: "POS", VolumeM3: 20.0, Priority: 5},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 3 {
		t.Fatalf("expected 3 moves, got %d", len(plan.Moves))
	}
	if plan.Moves[0].ProductID != "POS" {
		t.Errorf("expected first move to be POS, got %s", plan.Moves[0].ProductID)
	}
	if plan.Moves[1].ProductID != "ZERO" {
		t.Errorf("expected second move to be ZERO, got %s", plan.Moves[1].ProductID)
	}
	if plan.Moves[2].ProductID != "NEG" {
		t.Errorf("expected third move to be NEG, got %s", plan.Moves[2].ProductID)
	}
}

func TestDistribute_AllWarehousesFull_AllItemsUnallocated(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-1", TotalCapacityM3: 100.0, CurrentStockM3: 100.0},
		{ID: "WH-2", TotalCapacityM3: 50.0, CurrentStockM3: 50.0},
	}
	items := []algorithm.Product{
		{ID: "ITEM-1", VolumeM3: 10.0, Priority: 1},
		{ID: "ITEM-2", VolumeM3: 5.0, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 0 {
		t.Fatalf("expected 0 moves, got %d", len(plan.Moves))
	}
	if len(plan.UnallocatedItems) != 2 {
		t.Fatalf("expected 2 unallocated items, got %d", len(plan.UnallocatedItems))
	}
}

func TestDistribute_LargeItemWithSmallItemAfterFallback_BothAllocated(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-SMALL", TotalCapacityM3: 20.0, CurrentStockM3: 0},
		{ID: "WH-LARGE", TotalCapacityM3: 100.0, CurrentStockM3: 0},
	}
	items := []algorithm.Product{
		{ID: "BIG", VolumeM3: 80.0, Priority: 1},
		{ID: "SMALL", VolumeM3: 10.0, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 2 {
		t.Fatalf("expected 2 moves, got %d", len(plan.Moves))
	}
	if len(plan.UnallocatedItems) != 0 {
		t.Fatalf("expected 0 unallocated items, got %d", len(plan.UnallocatedItems))
	}
}

func TestDistribute_MultipleItemsExceedCapacity_OnlyFittingItemsAllocated(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-1", TotalCapacityM3: 100.0},
	}
	items := []algorithm.Product{
		{ID: "A", VolumeM3: 40.0, Priority: 1},
		{ID: "B", VolumeM3: 40.0, Priority: 1},
		{ID: "C", VolumeM3: 40.0, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 2 {
		t.Fatalf("expected 2 moves, got %d", len(plan.Moves))
	}
	if len(plan.UnallocatedItems) != 1 {
		t.Fatalf("expected 1 unallocated item, got %d", len(plan.UnallocatedItems))
	}
}

func TestDistribute_DifferentCapacityWarehouses_BalancesByUtilization(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-SMALL", TotalCapacityM3: 50.0},
		{ID: "WH-LARGE", TotalCapacityM3: 100.0},
	}
	items := []algorithm.Product{
		{ID: "ITEM-1", VolumeM3: 20.0, Priority: 1},
		{ID: "ITEM-2", VolumeM3: 20.0, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 2 {
		t.Fatalf("expected 2 moves, got %d", len(plan.Moves))
	}

	warehouseItems := make(map[string]int)
	for _, m := range plan.Moves {
		warehouseItems[m.WarehouseID]++
	}
	if warehouseItems["WH-SMALL"] != 1 || warehouseItems["WH-LARGE"] != 1 {
		t.Errorf("expected balanced distribution, got %v", warehouseItems)
	}
}

func TestDistribute_WarehouseWithOnlyIncomingStock_CalculatesCapacityCorrectly(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-1", TotalCapacityM3: 100.0, CurrentStockM3: 0, IncomingM3: 50.0},
	}
	items := []algorithm.Product{
		{ID: "ITEM-1", VolumeM3: 40.0, Priority: 1},
		{ID: "ITEM-2", VolumeM3: 20.0, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 1 {
		t.Fatalf("expected 1 move, got %d", len(plan.Moves))
	}
	if plan.Moves[0].ProductID != "ITEM-1" {
		t.Errorf("expected ITEM-1 to be allocated, got %s", plan.Moves[0].ProductID)
	}
	if len(plan.UnallocatedItems) != 1 {
		t.Fatalf("expected 1 unallocated item, got %d", len(plan.UnallocatedItems))
	}
}

func TestDistribute_VolumeGreaterThanPriority_VolumeIsSortedFirst(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-1", TotalCapacityM3: 100.0},
	}
	items := []algorithm.Product{
		{ID: "HIGH-PRIO-SMALL", VolumeM3: 10.0, Priority: 100},
		{ID: "LOW-PRIO-LARGE", VolumeM3: 50.0, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 2 {
		t.Fatalf("expected 2 moves, got %d", len(plan.Moves))
	}
	if plan.Moves[0].ProductID != "LOW-PRIO-LARGE" {
		t.Errorf("expected larger item first regardless of priority, got %s", plan.Moves[0].ProductID)
	}
}

func TestDistribute_SequentialAllocationUpdatesWarehouseState(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-1", TotalCapacityM3: 100.0},
		{ID: "WH-2", TotalCapacityM3: 100.0},
	}
	items := []algorithm.Product{
		{ID: "A", VolumeM3: 60.0, Priority: 1},
		{ID: "B", VolumeM3: 60.0, Priority: 1},
		{ID: "C", VolumeM3: 30.0, Priority: 1},
		{ID: "D", VolumeM3: 30.0, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 4 {
		t.Fatalf("expected 4 moves, got %d", len(plan.Moves))
	}
	if len(plan.UnallocatedItems) != 0 {
		t.Fatalf("expected 0 unallocated items, got %d", len(plan.UnallocatedItems))
	}

	warehouseLoad := make(map[string]float64)
	for _, m := range plan.Moves {
		warehouseLoad[m.WarehouseID] += m.VolumeM3
	}
	if warehouseLoad["WH-1"] != 90.0 || warehouseLoad["WH-2"] != 90.0 {
		t.Errorf("expected balanced loads of 90.0 each, got WH-1: %f, WH-2: %f", warehouseLoad["WH-1"], warehouseLoad["WH-2"])
	}
}

func TestDistribute_SingleWarehouseMultipleItems_CorrectOrder(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-1", TotalCapacityM3: 200.0},
	}
	items := []algorithm.Product{
		{ID: "C", VolumeM3: 10.0, Priority: 3},
		{ID: "A", VolumeM3: 30.0, Priority: 1},
		{ID: "B", VolumeM3: 20.0, Priority: 2},
		{ID: "D", VolumeM3: 10.0, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 4 {
		t.Fatalf("expected 4 moves, got %d", len(plan.Moves))
	}
	if plan.Moves[0].ProductID != "A" {
		t.Errorf("expected first move to be A (largest), got %s", plan.Moves[0].ProductID)
	}
	if plan.Moves[1].ProductID != "B" {
		t.Errorf("expected second move to be B, got %s", plan.Moves[1].ProductID)
	}
	if plan.Moves[2].ProductID != "C" {
		t.Errorf("expected third move to be C (same volume, higher priority), got %s", plan.Moves[2].ProductID)
	}
	if plan.Moves[3].ProductID != "D" {
		t.Errorf("expected fourth move to be D (same volume, lower priority), got %s", plan.Moves[3].ProductID)
	}
}

func TestDistribute_VerySmallVolumes_HandledCorrectly(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-1", TotalCapacityM3: 1.0},
	}
	items := []algorithm.Product{
		{ID: "TINY-1", VolumeM3: 0.001, Priority: 1},
		{ID: "TINY-2", VolumeM3: 0.001, Priority: 1},
		{ID: "TINY-3", VolumeM3: 0.001, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 3 {
		t.Fatalf("expected 3 moves, got %d", len(plan.Moves))
	}
	if len(plan.UnallocatedItems) != 0 {
		t.Fatalf("expected 0 unallocated items, got %d", len(plan.UnallocatedItems))
	}
}

func TestDistribute_BothStockAndIncoming_CombinedForUtilization(t *testing.T) {
	warehouses := []*algorithm.Warehouse{
		{ID: "WH-1", TotalCapacityM3: 100.0, CurrentStockM3: 30.0, IncomingM3: 30.0},
		{ID: "WH-2", TotalCapacityM3: 100.0, CurrentStockM3: 50.0, IncomingM3: 0},
	}
	items := []algorithm.Product{
		{ID: "ITEM-1", VolumeM3: 35.0, Priority: 1},
	}

	algo := algorithm.NewWFDAlgorithm()
	plan := algo.Distribute(warehouses, items)

	if len(plan.Moves) != 1 {
		t.Fatalf("expected 1 move, got %d", len(plan.Moves))
	}
	if plan.Moves[0].WarehouseID != "WH-2" {
		t.Errorf("expected item in WH-2 (50%% utilized vs 60%%), got %s", plan.Moves[0].WarehouseID)
	}
}
