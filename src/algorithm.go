package main

import (
	"sort"
)

// Item keeps original index and precomputed total weight
type Item struct {
	Index  int
	Weight []int
	Total  int
}

// Sort items in descending order of total resource usage
func sortItemsByWeight(items [][]int) []Item {
	result := make([]Item, len(items))
	for i, w := range items {
		total := 0
		for _, v := range w {
			total += v
		}
		result[i] = Item{i, w, total}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Total > result[j].Total
	})
	return result
}

// Greedy First-Fit Decreasing heuristic
func firstFitDecreasing(sortedItems []Item, knapsackCapacity []int, numKnapsacks int) []int {
	numDimensions := len(knapsackCapacity)
	usage := make([][]int, numKnapsacks)
	for i := range usage {
		usage[i] = make([]int, numDimensions)
	}
	assignment := make([]int, len(sortedItems))
	for i := range assignment {
		assignment[i] = -1
	}

	for _, item := range sortedItems {
		placed := false
		for k := 0; k < numKnapsacks && !placed; k++ {
			canFit := true
			for d := 0; d < numDimensions; d++ {
				if usage[k][d]+item.Weight[d] > knapsackCapacity[d] {
					canFit = false
					break
				}
			}
			if canFit {
				for d := 0; d < numDimensions; d++ {
					usage[k][d] += item.Weight[d]
				}
				assignment[item.Index] = k
				placed = true
			}
		}
		if !placed {
			return nil
		}
	}
	return assignment
}

func assignItemsToKnapsacks(itemWeights [][]int, knapsackCapacity []int, numKnapsacks int) []int {
	numItems := len(itemWeights)
	numDimensions := len(knapsackCapacity)

	// Early rejection: total demand exceeds total capacity
	totalDemand := make([]int, numDimensions)
	for _, item := range itemWeights {
		for d := 0; d < numDimensions; d++ {
			totalDemand[d] += item[d]
		}
	}
	for d := 0; d < numDimensions; d++ {
		if totalDemand[d] > knapsackCapacity[d]*numKnapsacks {
			return nil
		}
	}

	// Early acceptance: all items fit individually and there are enough knapsacks
	if numItems <= numKnapsacks {
		for _, item := range itemWeights {
			for d := 0; d < numDimensions; d++ {
				if item[d] > knapsackCapacity[d] {
					return nil
				}
			}
		}
		result := make([]int, numItems)
		for i := 0; i < numItems; i++ {
			result[i] = i
		}
		return result
	}

	// Sort items by decreasing weight
	sortedItems := sortItemsByWeight(itemWeights)

	// Fast Greedy Assignment
	if greedy := firstFitDecreasing(sortedItems, knapsackCapacity, numKnapsacks); greedy != nil {
		return greedy
	}

	// Full backtracking fallback
	usagePerKnapsack := make([][]int, numKnapsacks)
	for i := range usagePerKnapsack {
		usagePerKnapsack[i] = make([]int, numDimensions)
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
