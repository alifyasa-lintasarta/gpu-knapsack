package main

import (
	"fmt"
	"sort"
)

type Item struct {
	Index      int
	Weight     []int
	Time       int
	RemoveTime *int
}

type Event struct {
	Time      int
	Type      string // "assign" or "remove"
	ItemIndex int
	Weight    []int
}

func buildEventTimeline(items []AssignmentItem, itemWeights [][]int) []Event {
	var events []Event

	for i, item := range items {
		// Add assignment event
		events = append(events, Event{
			Time:      item.AssignmentTime,
			Type:      "assign",
			ItemIndex: i,
			Weight:    itemWeights[i],
		})

		// Add removal event if removeTime is specified
		if item.RemoveTime != nil {
			events = append(events, Event{
				Time:      *item.RemoveTime,
				Type:      "remove",
				ItemIndex: i,
				Weight:    itemWeights[i],
			})
		}
	}

	// Sort events by time, with removals before assignments at the same time
	sort.Slice(events, func(i, j int) bool {
		if events[i].Time == events[j].Time {
			return events[i].Type == "remove" && events[j].Type == "assign"
		}
		return events[i].Time < events[j].Time
	})

	return events
}

func runEventLoop(items []AssignmentItem, itemWeights [][]int, knapsackCapacity []int, numKnapsacks int, input *AssignmentInput) bool {
	numDimensions := len(knapsackCapacity)
	events := buildEventTimeline(items, itemWeights)

	// Track current usage for each GPU
	usage := make([][]int, numKnapsacks)
	for i := range usage {
		usage[i] = make([]int, numDimensions)
	}

	assignment := make([]int, len(items))
	for i := range assignment {
		assignment[i] = -1
	}

	fmt.Print("\n")
	fmt.Println("Simulation Starting...")
	fmt.Println("========================")
	fmt.Print("\n")

	// Process events chronologically
	for _, event := range events {
		if event.Type == "assign" {
			// Find first GPU with enough space (first-fit)
			placed := false
			for k := 0; k < numKnapsacks && !placed; k++ {
				canFit := true
				for d := 0; d < numDimensions; d++ {
					if usage[k][d]+event.Weight[d] > knapsackCapacity[d] {
						canFit = false
						break
					}
				}

				if canFit {
					// Assign to this GPU
					for d := 0; d < numDimensions; d++ {
						usage[k][d] += event.Weight[d]
					}
					assignment[event.ItemIndex] = k
					placed = true

					// Print state change
					fmt.Printf("Time %d: Added %s to GPU %d\n", event.Time, items[event.ItemIndex].Type, k)
					printGPUState(usage, knapsackCapacity, numKnapsacks)
				}
			}

			if !placed {
				fmt.Printf("Time %d: Failed to assign %s - no space available\n", event.Time, items[event.ItemIndex].Type)
				return false // No space available
			}

		} else if event.Type == "remove" {
			// Remove item from its assigned GPU
			gpuIndex := assignment[event.ItemIndex]
			if gpuIndex != -1 {
				for d := 0; d < numDimensions; d++ {
					usage[gpuIndex][d] -= event.Weight[d]
				}
				assignment[event.ItemIndex] = -1

				// Print state change
				fmt.Printf("Time %d: Removed %s from GPU %d\n", event.Time, items[event.ItemIndex].Type, gpuIndex)
				printGPUState(usage, knapsackCapacity, numKnapsacks)
			}
		}
	}

	// Store assignment in input for final summary
	input.Assignment = assignment
	return true
}

func printGPUState(usage [][]int, capacity []int, numKnapsacks int) {
	fmt.Print("  GPU Usage: ")
	for k := 0; k < numKnapsacks; k++ {
		if k > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("GPU%d[", k)
		for d := 0; d < len(capacity); d++ {
			if d > 0 {
				fmt.Print(",")
			}
			fmt.Printf("%d/%d", usage[k][d], capacity[d])
		}
		fmt.Print("]")
	}
	fmt.Println()
}
