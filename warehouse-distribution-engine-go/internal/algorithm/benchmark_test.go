package algorithm

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

// Baseline Algorithm (First-Fit)
func DistributeFirstFit(requestID string, warehouses []*Warehouse, items []Product) DistributionPlan {
	plan := DistributionPlan{RequestID: requestID, Moves: []Move{}}

	whCopies := make([]*Warehouse, len(warehouses))
	for i, w := range warehouses {
		whCopies[i] = &Warehouse{
			ID:              w.ID,
			TotalCapacityM3: w.TotalCapacityM3,
			CurrentStockM3:  w.CurrentStockM3,
			IncomingM3:      w.IncomingM3,
		}
	}

	for _, item := range items {
		allocated := false
		for _, wh := range whCopies {
			if wh.CanFit(item) {
				_ = wh.Allocate(item)
				plan.Moves = append(plan.Moves, Move{WarehouseID: wh.ID, ProductID: item.ID, VolumeM3: item.VolumeM3})
				allocated = true
				break
			}
		}
		if !allocated {
			plan.UnallocatedItems = append(plan.UnallocatedItems, item)
		}
	}
	return plan
}

func generateScenario() ([]*Warehouse, []Product) {
	warehouses := make([]*Warehouse, 10)
	for i := 0; i < 10; i++ {
		warehouses[i] = &Warehouse{
			ID:              fmt.Sprintf("WH-%d", i),
			TotalCapacityM3: 1000.0,
		}
	}

	items := make([]Product, 2500)
	for i := 0; i < 2500; i++ {
		items[i] = Product{
			ID:       fmt.Sprintf("ITEM-%d", i),
			VolumeM3: 1.0 + rand.Float64()*2.0,
		}
	}
	return warehouses, items
}

// Cv calculation
func calculateCV(warehouses []*Warehouse, plan DistributionPlan) float64 {
	usageMap := make(map[string]float64)
	for _, m := range plan.Moves {
		usageMap[m.WarehouseID] += m.VolumeM3
	}

	var loadRatios []float64
	var sumRatios float64

	for _, w := range warehouses {
		used := usageMap[w.ID]
		ratio := used / w.TotalCapacityM3
		loadRatios = append(loadRatios, ratio)
		sumRatios += ratio
	}

	mean := sumRatios / float64(len(warehouses))

	if mean == 0 {
		return 0
	}

	var varianceSum float64
	for _, r := range loadRatios {
		varianceSum += math.Pow(r-mean, 2)
	}

	stdDev := math.Sqrt(varianceSum / float64(len(warehouses)))

	return stdDev / mean
}

func TestCompareAlgorithms(t *testing.T) {
	warehouses, items := generateScenario()
	requestID := "benchmark-test"

	// A. First-Fit
	baselinePlan := DistributeFirstFit(requestID, warehouses, items)
	baselineCV := calculateCV(warehouses, baselinePlan)

	// B. WFD
	whForWFD := make([]*Warehouse, len(warehouses))
	for i, w := range warehouses {
		whForWFD[i] = &Warehouse{ID: w.ID, TotalCapacityM3: w.TotalCapacityM3}
	}

	packer := NewWFDAlgorithm()
	wfdPlan := packer.Distribute(requestID, whForWFD, items)
	wfdCV := calculateCV(whForWFD, wfdPlan)

	fmt.Println("\n=== EXPERIMENT RESULTS (Load Balancing) ===")
	fmt.Printf("Input data: %d warehouses, %d items (Load ~50%%)\n", len(warehouses), len(items))
	fmt.Println("-------------------------------------------------------")
	fmt.Printf("| %-15s | %-12s | %-12s |\n", "Algorithm", "CV (Imbalance)", "Unallocated")
	fmt.Println("-------------------------------------------------------")
	fmt.Printf("| %-15s | %-12.4f | %-12d |\n", "First-Fit", baselineCV, len(baselinePlan.UnallocatedItems))
	fmt.Printf("| %-15s | %-12.4f | %-12d |\n", "WFD (Proposed)", wfdCV, len(wfdPlan.UnallocatedItems))
	fmt.Println("-------------------------------------------------------")

	if wfdCV < baselineCV {
		improvement := (baselineCV - wfdCV) / baselineCV * 100
		t.Logf("Success! WFD distributed load %.2f%% more evenly.", improvement)
	} else {
		t.Log("WFD did not show better balancing.")
	}
}
