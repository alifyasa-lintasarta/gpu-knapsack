package main

import (
	"fmt"
)

type Event struct {
	Time      int
	Type      string
	ItemIndex int
	Weight    []int
}


func processEvents(items []PodSpec, itemWeights [][]int, gpuCapacity []int, numGPUs int, input *SchedulingInput) bool {
	capacityDimensions := len(gpuCapacity)
	events := input.GPUFamily.buildEventTimeline(items, itemWeights)

	usage := make([][]int, numGPUs)
	for i := range usage {
		usage[i] = make([]int, capacityDimensions)
	}

	assignment := make([]int, len(items))
	for i := range assignment {
		assignment[i] = -1
	}

	fmt.Print("\n")
	fmt.Println("Simulation Starting...")
	fmt.Println("========================")
	fmt.Print("\n")

	for _, event := range events {
		if event.Type == "assign" {
			placed := false
			for gpuIdx := 0; gpuIdx < numGPUs && !placed; gpuIdx++ {
				canFit := true
				for d := 0; d < capacityDimensions; d++ {
					if usage[gpuIdx][d]+event.Weight[d] > gpuCapacity[d] {
						canFit = false
						break
					}
				}

				if canFit {
					for d := 0; d < capacityDimensions; d++ {
						usage[gpuIdx][d] += event.Weight[d]
					}
					assignment[event.ItemIndex] = gpuIdx
					placed = true

					fmt.Printf("Time %d: Added %s to GPU %d\n", event.Time, items[event.ItemIndex].Type, gpuIdx)
					printGPUState(usage, gpuCapacity, numGPUs)
				}
			}

			if !placed {
				fmt.Printf("Time %d: Failed to assign %s - no space available\n", event.Time, items[event.ItemIndex].Type)
				return false
			}

		} else if event.Type == "remove" {
			gpuIndex := assignment[event.ItemIndex]
			if gpuIndex != -1 {
				for d := 0; d < capacityDimensions; d++ {
					usage[gpuIndex][d] -= event.Weight[d]
				}
				assignment[event.ItemIndex] = -1

				fmt.Printf("Time %d: Removed %s from GPU %d\n", event.Time, items[event.ItemIndex].Type, gpuIndex)
				printGPUState(usage, gpuCapacity, numGPUs)
			}
		}
	}

	input.Assignment = assignment
	return true
}

func printGPUState(usage [][]int, capacity []int, numGPUs int) {
	fmt.Print("  GPU Usage: ")
	for gpuIdx := 0; gpuIdx < numGPUs; gpuIdx++ {
		if gpuIdx > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("GPU%d[", gpuIdx)
		for d := 0; d < len(capacity); d++ {
			if d > 0 {
				fmt.Print(",")
			}
			fmt.Printf("%d/%d", usage[gpuIdx][d], capacity[d])
		}
		fmt.Print("]")
	}
	fmt.Println()
}