package main

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

// Global memoization caches
var (
	assignmentCache = make(map[string]bool)
	knapsackCache   = make(map[string][]int)
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
	// Create cache key for knapsack solver
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
	cacheKey := fmt.Sprintf("%x", h.Sum(nil))

	if cached, exists := knapsackCache[cacheKey]; exists {
		return cached
	}

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

	var result []int
	if backtrack(0) {
		result = itemAssignment
	} else {
		result = nil
	}

	// Cache the result
	knapsackCache[cacheKey] = result
	return result
}

func canAssignWithAdditional(cfg Config, additional map[string]int) bool {
	// Create cache key from the combination
	cacheKey := combinationToString(additional)
	if cached, exists := assignmentCache[cacheKey]; exists {
		return cached
	}

	// Create new pod configuration
	testPods := make(map[string]int)
	for k, v := range cfg.Pods {
		testPods[k] = v
	}
	for k, v := range additional {
		testPods[k] = testPods[k] + v
	}

	// Pre-calculate total pods for allocation
	totalPods := 0
	for _, count := range testPods {
		totalPods += count
	}

	// Build request list with pre-allocated capacity
	requests := make([]string, 0, totalPods)
	for podType, count := range testPods {
		for i := 0; i < count; i++ {
			requests = append(requests, podType)
		}
	}

	// Build item weights with pre-allocated capacity
	itemWeights := make([][]int, 0, len(requests))
	for _, gpu := range requests {
		weights, ok := cfg.GPU.Mappings[gpu]
		if !ok {
			assignmentCache[cacheKey] = false
			return false
		}
		itemWeights = append(itemWeights, weights)
	}

	// Test if assignment is possible
	assignment := assignItemsToKnapsacks(itemWeights, cfg.GPU.Capacity, cfg.GPU.Number)
	result := assignment != nil
	assignmentCache[cacheKey] = result
	return result
}

// Helper function to create a canonical string representation of a combination
func combinationToString(combination map[string]int) string {
	var builder strings.Builder
	var parts []string
	for podType, count := range combination {
		if count > 0 {
			parts = append(parts, fmt.Sprintf("%s:%d", podType, count))
		}
	}
	sort.Strings(parts)
	for i, part := range parts {
		if i > 0 {
			builder.WriteByte(',')
		}
		builder.WriteString(part)
	}
	return builder.String()
}

// findAllPossibleCombinations finds all maximal feasible combinations of additional pods
func findAllPossibleCombinations(cfg Config) {
	// Get all pod types
	podTypes := make([]string, 0, len(cfg.GPU.Mappings))
	for podType := range cfg.GPU.Mappings {
		podTypes = append(podTypes, podType)
	}

	// Generate all combinations up to a reasonable limit
	var allFeasible []map[string]int
	maxPerType := 20 // Reasonable limit per pod type

	// Use iterative approach to generate all combinations
	generateCombinations(cfg, podTypes, make(map[string]int), 0, maxPerType, &allFeasible)

	// Filter for maximal combinations (can't add any more pods)
	var maximal []map[string]int
	for _, combination := range allFeasible {
		if isMaximalFeasible(cfg, combination) {
			maximal = append(maximal, combination)
		}
	}

	// Sort and print maximal combinations
	sort.Slice(maximal, func(i, j int) bool {
		return combinationToString(maximal[i]) < combinationToString(maximal[j])
	})

	fmt.Println("\nMaximal additional pod combinations you can add:")
	for i, combination := range maximal {
		fmt.Printf("%d. ", i+1)
		first := true
		totalPods := 0
		for pType, count := range combination {
			if count > 0 {
				if !first {
					fmt.Print(", ")
				}
				fmt.Printf("%s: %d", pType, count)
				totalPods += count
				first = false
			}
		}
		if totalPods == 0 {
			fmt.Print("No additional pods can be added")
		}
		fmt.Println()
	}
}

// generateCombinations generates all possible combinations recursively with early pruning
func generateCombinations(cfg Config, podTypes []string, current map[string]int, typeIndex int, maxPerType int, results *[]map[string]int) {
	if typeIndex >= len(podTypes) {
		// Test if this combination is feasible (can be added to current system)
		if canAssignWithAdditional(cfg, current) {
			// Make a copy of current combination
			combination := make(map[string]int)
			for k, v := range current {
				combination[k] = v
			}
			*results = append(*results, combination)
		}
		return
	}

	podType := podTypes[typeIndex]

	// Try all counts from 0 to maxPerType for this pod type
	for count := 0; count <= maxPerType; count++ {
		if count > 0 {
			current[podType] = count
		}

		// Early pruning: if current combination already exceeds capacity, skip this branch
		if count > 0 && !canAssignWithAdditional(cfg, current) {
			// If adding this count makes it infeasible, no point trying higher counts
			if count > 0 {
				delete(current, podType)
			}
			break
		}

		generateCombinations(cfg, podTypes, current, typeIndex+1, maxPerType, results)

		if count > 0 {
			delete(current, podType)
		}
	}
}

// isMaximalFeasible checks if a combination is maximal (can't add any more pods)
func isMaximalFeasible(cfg Config, combination map[string]int) bool {
	// Get all pod types
	podTypes := make([]string, 0, len(cfg.GPU.Mappings))
	for podType := range cfg.GPU.Mappings {
		podTypes = append(podTypes, podType)
	}

	// Try to add one more pod of each type
	for _, podType := range podTypes {
		testCombination := make(map[string]int)
		for k, v := range combination {
			testCombination[k] = v
		}
		testCombination[podType] = testCombination[podType] + 1

		// If we can still add more, then current combination is not maximal
		if canAssignWithAdditional(cfg, testCombination) {
			return false
		}
	}

	return true
}
