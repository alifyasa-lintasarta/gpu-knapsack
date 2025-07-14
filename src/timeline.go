package main

import (
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

func sortItemsByTime(items []AssignmentItem, weights [][]int) []Item {
	result := make([]Item, len(items))
	for i, item := range items {
		result[i] = Item{i, weights[i], item.AssignmentTime, item.RemoveTime}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Time < result[j].Time
	})
	return result
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
		
		// Add removal event if remove_time is specified
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

func tryTimelineAssignment(items []AssignmentItem, itemWeights [][]int, knapsackCapacity []int, numKnapsacks int) []int {
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
				}
			}
			
			if !placed {
				return nil // No space available
			}
			
		} else if event.Type == "remove" {
			// Remove item from its assigned GPU
			gpuIndex := assignment[event.ItemIndex]
			if gpuIndex != -1 {
				for d := 0; d < numDimensions; d++ {
					usage[gpuIndex][d] -= event.Weight[d]
				}
			}
		}
	}
	
	return assignment
}