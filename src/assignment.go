package main

import (
	"crypto/sha256"
	"fmt"
)

// Global memoization caches
var (
	assignmentCache = make(map[string]bool)
	knapsackCache   = make(map[string][]int)
)

func generateCacheKey(itemWeights [][]int, knapsackCapacity []int, numKnapsacks int) string {
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
	return fmt.Sprintf("%x", h.Sum(nil))
}

func validateInputs(itemWeights [][]int, knapsackCapacity []int, numKnapsacks int) []int {
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

	return []int{} // Empty slice indicates validation passed but no early result
}

func assignItemsToKnapsacks(itemWeights [][]int, knapsackCapacity []int, numKnapsacks int) []int {
	cacheKey := generateCacheKey(itemWeights, knapsackCapacity, numKnapsacks)

	if cached, exists := knapsackCache[cacheKey]; exists {
		return cached
	}

	validationResult := validateInputs(itemWeights, knapsackCapacity, numKnapsacks)
	if validationResult == nil {
		return nil
	}
	if len(validationResult) > 0 {
		knapsackCache[cacheKey] = validationResult
		return validationResult
	}

	if greedy := tryGreedyAssignment(itemWeights, knapsackCapacity, numKnapsacks); greedy != nil {
		knapsackCache[cacheKey] = greedy
		return greedy
	}

	result := tryBacktrackingAssignment(itemWeights, knapsackCapacity, numKnapsacks)
	knapsackCache[cacheKey] = result
	return result
}
