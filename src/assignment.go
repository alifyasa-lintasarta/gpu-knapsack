package main

import (
	"crypto/sha256"
	"fmt"
)

// Global memoization caches
var (
	knapsackCache = make(map[string][]int)
)

func generateCacheKey(itemWeights [][]int, knapsackCapacity []int, numKnapsacks int, initialUsage [][]int) string {
	h := sha256.New()
	for _, item := range itemWeights {
		for _, w := range item {
			h.Write([]byte(fmt.Sprintf("%d,", w)))
		}
		h.Write([]byte(";"))
	}
	for _, cap := range knapsackCapacity {
		h.Write([]byte(fmt.Sprintf("%d,", cap)))
	}
	h.Write([]byte(fmt.Sprintf("n%d", numKnapsacks)))

	for _, usage := range initialUsage {
		for _, u := range usage {
			h.Write([]byte(fmt.Sprintf("u%d,", u)))
		}
		h.Write([]byte("|"))
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func validateInputs(itemWeights [][]int, knapsackCapacity []int, numKnapsacks int, initialUsage [][]int) []int {
	numItems := len(itemWeights)
	numDimensions := len(knapsackCapacity)

	// Calculate remaining capacity after initial usage
	remainingCapacity := make([]int, numDimensions)
	for d := 0; d < numDimensions; d++ {
		remainingCapacity[d] = knapsackCapacity[d] * numKnapsacks
		for k := 0; k < numKnapsacks && k < len(initialUsage); k++ {
			if d < len(initialUsage[k]) {
				remainingCapacity[d] -= initialUsage[k][d]
			}
		}
	}

	// Early rejection: total demand exceeds remaining capacity
	totalDemand := make([]int, numDimensions)
	for _, item := range itemWeights {
		for d := 0; d < numDimensions; d++ {
			totalDemand[d] += item[d]
		}
	}
	for d := 0; d < numDimensions; d++ {
		if totalDemand[d] > remainingCapacity[d] {
			return nil
		}
	}

	// Find available knapsacks for early acceptance check
	availableKnapsacks := 0
	for k := 0; k < numKnapsacks; k++ {
		hasCapacity := true
		for d := 0; d < numDimensions; d++ {
			currentUsage := 0
			if k < len(initialUsage) && d < len(initialUsage[k]) {
				currentUsage = initialUsage[k][d]
			}
			if currentUsage >= knapsackCapacity[d] {
				hasCapacity = false
				break
			}
		}
		if hasCapacity {
			availableKnapsacks++
		}
	}

	// Early acceptance: all items fit individually and there are enough available knapsacks
	if numItems <= availableKnapsacks {
		for _, item := range itemWeights {
			for d := 0; d < numDimensions; d++ {
				if item[d] > knapsackCapacity[d] {
					return nil
				}
			}
		}

		// Try to assign items to available knapsacks
		result := make([]int, numItems)
		assignedCount := 0
		for k := 0; k < numKnapsacks && assignedCount < numItems; k++ {
			canFit := true
			for d := 0; d < numDimensions; d++ {
				currentUsage := 0
				if k < len(initialUsage) && d < len(initialUsage[k]) {
					currentUsage = initialUsage[k][d]
				}
				if currentUsage >= knapsackCapacity[d] {
					canFit = false
					break
				}
			}
			if canFit {
				result[assignedCount] = k
				assignedCount++
			}
		}
		if assignedCount == numItems {
			return result
		}
	}

	return []int{} // Empty slice indicates validation passed but no early result
}

func assignItemsToKnapsacksWithInitial(itemWeights [][]int, knapsackCapacity []int, numKnapsacks int, initialUsage [][]int) []int {
	cacheKey := generateCacheKey(itemWeights, knapsackCapacity, numKnapsacks, initialUsage)

	if cached, exists := knapsackCache[cacheKey]; exists {
		return cached
	}

	validationResult := validateInputs(itemWeights, knapsackCapacity, numKnapsacks, initialUsage)
	if validationResult == nil {
		return nil
	}
	if len(validationResult) > 0 {
		knapsackCache[cacheKey] = validationResult
		return validationResult
	}

	if greedy := tryGreedyAssignmentWithInitial(itemWeights, knapsackCapacity, numKnapsacks, initialUsage); greedy != nil {
		knapsackCache[cacheKey] = greedy
		return greedy
	}

	result := tryBacktrackingAssignmentWithInitial(itemWeights, knapsackCapacity, numKnapsacks, initialUsage)
	knapsackCache[cacheKey] = result
	return result
}
