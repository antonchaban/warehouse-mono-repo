package algorithm

import (
	"container/heap"
	"sort"
)

// Packer defines the interface for distribution logic.
type Packer interface {
	Distribute(requestID string, warehouses []*Warehouse, items []Product) DistributionPlan
}

// WFDAlgorithm implements the Weighted Worst-Fit Decreasing strategy.
type WFDAlgorithm struct{}

// NewWFDAlgorithm creates a new instance of the algorithm.
func NewWFDAlgorithm() *WFDAlgorithm {
	return &WFDAlgorithm{}
}

// Distribute executes the calculation logic.
func (a *WFDAlgorithm) Distribute(requestID string, warehouses []*Warehouse, items []Product) DistributionPlan {
	plan := DistributionPlan{
		RequestID:        requestID,
		Moves:            make([]Move, 0),
		UnallocatedItems: make([]Product, 0),
	}

	// STEP 1: Sort items (Descending strategy).
	sort.Slice(items, func(i, j int) bool {
		if items[i].VolumeM3 == items[j].VolumeM3 {
			return items[i].Priority > items[j].Priority
		}
		return items[i].VolumeM3 > items[j].VolumeM3
	})

	// STEP 2: Initialize Warehouse Priority Queue (Min-Heap).
	pq := make(WarehousePriorityQueue, len(warehouses))
	for i, w := range warehouses {
		pq[i] = w
	}
	heap.Init(&pq)

	// STEP 3: Allocation Loop.
	for _, item := range items {
		if pq.Len() == 0 {
			plan.UnallocatedItems = append(plan.UnallocatedItems, item)
			continue
		}

		bestWh := heap.Pop(&pq).(*Warehouse)

		if bestWh.CanFit(item) {
			_ = bestWh.Allocate(item)

			plan.Moves = append(plan.Moves, Move{
				ProductID:   item.ID,
				WarehouseID: bestWh.ID,
				VolumeM3:    item.VolumeM3,
			})

			heap.Push(&pq, bestWh)
		} else {
			// Fallback logic
			foundAlternative := false
			poppedBuffer := []*Warehouse{bestWh}

			for pq.Len() > 0 {
				candidate := heap.Pop(&pq).(*Warehouse)
				poppedBuffer = append(poppedBuffer, candidate)

				if candidate.CanFit(item) {
					_ = candidate.Allocate(item)
					plan.Moves = append(plan.Moves, Move{
						ProductID:   item.ID,
						WarehouseID: candidate.ID,
						VolumeM3:    item.VolumeM3,
					})
					foundAlternative = true
					break
				}
			}

			for _, w := range poppedBuffer {
				heap.Push(&pq, w)
			}

			if !foundAlternative {
				plan.UnallocatedItems = append(plan.UnallocatedItems, item)
			}
		}
	}

	return plan
}

// --- Priority Queue Implementation ---

// WarehousePriorityQueue implements heap.Interface.
type WarehousePriorityQueue []*Warehouse

// Len returns the number of elements in the collection.
// Receiver is a pointer for consistency.
func (pq *WarehousePriorityQueue) Len() int { return len(*pq) }

// Less determines the sorting order.
// We want the LEAST utilized warehouse to be popped first (Min-Heap).
func (pq *WarehousePriorityQueue) Less(i, j int) bool {
	return (*pq)[i].UtilizationRatio() < (*pq)[j].UtilizationRatio()
}

// Swap swaps the elements with indexes i and j.
func (pq *WarehousePriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
}

// Push adds x as the last element.
func (pq *WarehousePriorityQueue) Push(x interface{}) {
	item := x.(*Warehouse)
	*pq = append(*pq, item)
}

// Pop removes and returns the last element.
func (pq *WarehousePriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}
