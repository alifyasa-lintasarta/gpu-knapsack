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
	// For backward compatibility when no initial state is considered
	return testAssignmentWithInitial(cfg, testPods)
}

func testAssignmentWithInitial(cfg Config, testPods map[string]int) bool {
	requests := buildRequestsFromPods(testPods)

	itemWeights := make([][]int, 0, len(requests))
	for _, gpu := range requests {
		weights, ok := cfg.GPU.Mappings[gpu]
		if !ok {
			return false
		}
		itemWeights = append(itemWeights, weights)
	}

	// Compute initial usage from configuration
	initialUsage := computeInitialUsage(cfg.GPU.InitialState, cfg.GPU.Mappings, cfg.GPU.Number, len(cfg.GPU.Capacity))
	
	assignment := assignItemsToKnapsacksWithInitial(itemWeights, cfg.GPU.Capacity, cfg.GPU.Number, initialUsage)
	return assignment != nil
}

func canAssignWithAdditional(cfg Config, additional map[string]int) bool {
	cacheKey := combinationToString(additional)
	if cached, exists := assignmentCache[cacheKey]; exists {
		return cached
	}

	testPods := buildTestConfiguration(cfg, additional)
	result := testAssignmentWithInitial(cfg, testPods)
	assignmentCache[cacheKey] = result
	return result
}

func validateInitialState(cfg Config) error {
	for gpuIndex, pods := range cfg.GPU.InitialState {
		if gpuIndex < 0 || gpuIndex >= cfg.GPU.Number {
			return fmt.Errorf("invalid GPU index %d: must be between 0 and %d", gpuIndex, cfg.GPU.Number-1)
		}

		usage := make([]int, len(cfg.GPU.Capacity))
		for _, podType := range pods {
			weights, exists := cfg.GPU.Mappings[podType]
			if !exists {
				return fmt.Errorf("unknown pod type '%s' in initial state for GPU %d", podType, gpuIndex)
			}
			
			for d := 0; d < len(cfg.GPU.Capacity); d++ {
				usage[d] += weights[d]
				if usage[d] > cfg.GPU.Capacity[d] {
					return fmt.Errorf("initial state for GPU %d exceeds capacity in dimension %d: usage=%d, capacity=%d", 
						gpuIndex, d, usage[d], cfg.GPU.Capacity[d])
				}
			}
		}
	}
	return nil
}

func computeInitialUsage(initialState map[int][]string, mappings map[string][]int, numGPUs, numDimensions int) [][]int {
	usage := make([][]int, numGPUs)
	for i := range usage {
		usage[i] = make([]int, numDimensions)
	}

	for gpuIndex, pods := range initialState {
		if gpuIndex >= 0 && gpuIndex < numGPUs {
			for _, podType := range pods {
				if weights, exists := mappings[podType]; exists {
					for d := 0; d < numDimensions && d < len(weights); d++ {
						usage[gpuIndex][d] += weights[d]
					}
				}
			}
		}
	}

	return usage
}

func mergeAssignments(initialState map[int][]string, newAssignment []int, newPods []string) ([]int, []string) {
	allAssignments := make([]int, 0)
	allPods := make([]string, 0)

	for gpuIndex, pods := range initialState {
		for _, podType := range pods {
			allAssignments = append(allAssignments, gpuIndex)
			allPods = append(allPods, podType)
		}
	}

	for i, gpuIndex := range newAssignment {
		if i < len(newPods) {
			allAssignments = append(allAssignments, gpuIndex)
			allPods = append(allPods, newPods[i])
		}
	}

	return allAssignments, allPods
}
