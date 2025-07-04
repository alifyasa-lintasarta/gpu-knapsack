package main

import (
	"fmt"
	"sort"
	"strings"
)

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

func buildTestConfiguration(cfg Config, additional map[string]int) map[string]int {
	testPods := make(map[string]int)
	for k, v := range cfg.Pods {
		testPods[k] = v
	}
	for k, v := range additional {
		testPods[k] = testPods[k] + v
	}
	return testPods
}

func buildRequestsFromPods(testPods map[string]int) []string {
	totalPods := 0
	for _, count := range testPods {
		totalPods += count
	}

	requests := make([]string, 0, totalPods)
	for podType, count := range testPods {
		for i := 0; i < count; i++ {
			requests = append(requests, podType)
		}
	}
	return requests
}

func testAssignment(cfg Config, testPods map[string]int) bool {
	requests := buildRequestsFromPods(testPods)

	itemWeights := make([][]int, 0, len(requests))
	for _, gpu := range requests {
		weights, ok := cfg.GPU.Mappings[gpu]
		if !ok {
			return false
		}
		itemWeights = append(itemWeights, weights)
	}

	assignment := assignItemsToKnapsacks(itemWeights, cfg.GPU.Capacity, cfg.GPU.Number)
	return assignment != nil
}

func canAssignWithAdditional(cfg Config, additional map[string]int) bool {
	cacheKey := combinationToString(additional)
	if cached, exists := assignmentCache[cacheKey]; exists {
		return cached
	}

	testPods := buildTestConfiguration(cfg, additional)
	result := testAssignment(cfg, testPods)
	assignmentCache[cacheKey] = result
	return result
}