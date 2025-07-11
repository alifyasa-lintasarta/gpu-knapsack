package main

import (
	"fmt"
)

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
