package main

import (
	"sort"
)

// Item keeps original index and timestamp for ordering
type Item struct {
	Index  int
	Weight []int
	Time   int
}

// Sort items by assignment time (earliest first)
func sortItemsByTime(items []AssignmentItem, weights [][]int) []Item {
	result := make([]Item, len(items))
	for i, item := range items {
		result[i] = Item{i, weights[i], item.AssignmentTime}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Time < result[j].Time
	})
	return result
}

func tryBacktrackingAssignment(items []AssignmentItem, itemWeights [][]int, knapsackCapacity []int, numKnapsacks int) []int {
	numItems := len(items)
	numDimensions := len(knapsackCapacity)
	sortedItems := sortItemsByTime(items, itemWeights)

	usagePerKnapsack := make([][]int, numKnapsacks)
	for i := range usagePerKnapsack {
		usagePerKnapsack[i] = make([]int, numDimensions)
	}
	itemAssignment := make([]int, numItems)
	for i := range itemAssignment {
		itemAssignment[i] = -1
	}

	var backtrack func(int) bool
	backtrack = func(itemIndex int) bool {
		if itemIndex == len(sortedItems) {
			return true
		}

		item := sortedItems[itemIndex]

		// Try knapsacks in order (deterministic)
		for k := 0; k < numKnapsacks; k++ {
			canFit := true
			for d := 0; d < numDimensions; d++ {
				if usagePerKnapsack[k][d]+item.Weight[d] > knapsackCapacity[d] {
					canFit = false
					break
				}
			}
			if !canFit {
				continue
			}

			for d := 0; d < numDimensions; d++ {
				usagePerKnapsack[k][d] += item.Weight[d]
			}
			itemAssignment[item.Index] = k

			if backtrack(itemIndex + 1) {
				return true
			}

			for d := 0; d < numDimensions; d++ {
				usagePerKnapsack[k][d] -= item.Weight[d]
			}
			itemAssignment[item.Index] = -1
		}
		return false
	}

	if backtrack(0) {
		return itemAssignment
	}
	return nil
}

// Legacy sorting function for backward compatibility
func sortItemsByWeightLegacy(items [][]int) []Item {
	result := make([]Item, len(items))
	for i, w := range items {
		total := 0
		for _, v := range w {
			total += v
		}
		result[i] = Item{i, w, total} // Using total as time for legacy
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Time > result[j].Time // Descending by total weight
	})
	return result
}

func tryBacktrackingAssignmentWithInitial(itemWeights [][]int, knapsackCapacity []int, numKnapsacks int, initialUsage [][]int) []int {
	numItems := len(itemWeights)
	numDimensions := len(knapsackCapacity)
	
	// Keep old weight-based sorting for backward compatibility
	sortedItems := sortItemsByWeightLegacy(itemWeights)

	usagePerKnapsack := make([][]int, numKnapsacks)
	for i := range usagePerKnapsack {
		usagePerKnapsack[i] = make([]int, numDimensions)
		// Initialize with existing usage
		for d := 0; d < numDimensions; d++ {
			if i < len(initialUsage) && d < len(initialUsage[i]) {
				usagePerKnapsack[i][d] = initialUsage[i][d]
			}
		}
	}
	itemAssignment := make([]int, numItems)
	for i := range itemAssignment {
		itemAssignment[i] = -1
	}

	var sumUsage = func(u []int) int {
		s := 0
		for _, x := range u {
			s += x
		}
		return s
	}

	var backtrack func(int) bool
	backtrack = func(itemIndex int) bool {
		if itemIndex == len(sortedItems) {
			return true
		}

		// Early termination: check if remaining items can fit in remaining capacity
		remainingCapacity := make([]int, numDimensions)
		for k := 0; k < numKnapsacks; k++ {
			for d := 0; d < numDimensions; d++ {
				remainingCapacity[d] += knapsackCapacity[d] - usagePerKnapsack[k][d]
			}
		}

		remainingDemand := make([]int, numDimensions)
		for i := itemIndex; i < len(sortedItems); i++ {
			item := sortedItems[i]
			for d := 0; d < numDimensions; d++ {
				remainingDemand[d] += item.Weight[d]
			}
		}

		for d := 0; d < numDimensions; d++ {
			if remainingDemand[d] > remainingCapacity[d] {
				return false
			}
		}

		triedEmpty := false
		item := sortedItems[itemIndex]

		knapsackOrder := make([]int, numKnapsacks)
		for i := 0; i < numKnapsacks; i++ {
			knapsackOrder[i] = i
		}
		// Optional: prioritize knapsacks with lowest current usage
		sort.Slice(knapsackOrder, func(i, j int) bool {
			return sumUsage(usagePerKnapsack[knapsackOrder[i]]) < sumUsage(usagePerKnapsack[knapsackOrder[j]])
		})

		for _, k := range knapsackOrder {
			canFit := true
			for d := 0; d < numDimensions; d++ {
				if usagePerKnapsack[k][d]+item.Weight[d] > knapsackCapacity[d] {
					canFit = false
					break
				}
			}
			if !canFit {
				continue
			}

			wasEmpty := true
			for d := 0; d < numDimensions; d++ {
				if usagePerKnapsack[k][d] != 0 {
					wasEmpty = false
					break
				}
			}
			if wasEmpty && triedEmpty {
				continue
			}
			if wasEmpty {
				triedEmpty = true
			}

			for d := 0; d < numDimensions; d++ {
				usagePerKnapsack[k][d] += item.Weight[d]
			}
			itemAssignment[item.Index] = k

			if backtrack(itemIndex + 1) {
				return true
			}

			for d := 0; d < numDimensions; d++ {
				usagePerKnapsack[k][d] -= item.Weight[d]
			}
			itemAssignment[item.Index] = -1
		}
		return false
	}

	if backtrack(0) {
		return itemAssignment
	}
	return nil
}