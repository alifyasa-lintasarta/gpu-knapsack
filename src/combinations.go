package main

import (
	"fmt"
	"sort"
)

func extractPodTypes(cfg Config) []string {
	podTypes := make([]string, 0, len(cfg.GPU.Mappings))
	for podType := range cfg.GPU.Mappings {
		podTypes = append(podTypes, podType)
	}
	return podTypes
}

func filterMaximalCombinations(cfg Config, allFeasible []map[string]int) []map[string]int {
	var maximal []map[string]int
	for _, combination := range allFeasible {
		if isMaximalFeasible(cfg, combination) {
			maximal = append(maximal, combination)
		}
	}
	return maximal
}

func printMaximalCombinations(maximal []map[string]int) {
	sort.Slice(maximal, func(i, j int) bool {
		return combinationToString(maximal[i]) < combinationToString(maximal[j])
	})

	fmt.Println("\nAdditional Pod Combinations:")
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

func findAllPossibleCombinations(cfg Config) []map[string]int {
	podTypes := extractPodTypes(cfg)

	var allFeasible []map[string]int
	maxPerType := 20

	generateCombinations(cfg, podTypes, make(map[string]int), 0, maxPerType, &allFeasible)

	maximal := filterMaximalCombinations(cfg, allFeasible)
	return maximal
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
