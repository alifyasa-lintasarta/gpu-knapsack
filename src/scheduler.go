package main

import (
	"fmt"
)

// New data structures for timestamp-based assignment
type AssignmentItem struct {
	Type           string
	AssignmentTime int
	RemoveTime     *int
}

type AssignmentInput struct {
	Items            []AssignmentItem
	KnapsackCapacity []int
	NumKnapsacks     int
	Mappings         map[string][]int
	Assignment       []int // populated by function
}

func validateAssignmentInput(input *AssignmentInput) error {
	if len(input.Items) == 0 {
		return fmt.Errorf("no items to assign")
	}
	if input.NumKnapsacks <= 0 {
		return fmt.Errorf("number of knapsacks must be positive")
	}
	if len(input.KnapsackCapacity) == 0 {
		return fmt.Errorf("knapsack capacity must be specified")
	}
	if input.Mappings == nil {
		return fmt.Errorf("mappings cannot be nil")
	}

	// Validate all item types have mappings
	for _, item := range input.Items {
		if _, exists := input.Mappings[item.Type]; !exists {
			return fmt.Errorf("no mapping found for item type %s", item.Type)
		}
	}

	return nil
}

func AssignItems(input *AssignmentInput) (bool, error) {
	if input == nil {
		return false, fmt.Errorf("input cannot be nil")
	}

	// Validate input
	if err := validateAssignmentInput(input); err != nil {
		return false, err
	}

	// Build item weights from mappings
	itemWeights := make([][]int, len(input.Items))
	for i, item := range input.Items {
		weights, exists := input.Mappings[item.Type]
		if !exists {
			return false, fmt.Errorf("no mapping found for item type %s", item.Type)
		}
		itemWeights[i] = weights
	}

	// Use timeline-based assignment (handles remove_time properly)
	if assignment := tryTimelineAssignment(input.Items, itemWeights, input.KnapsackCapacity, input.NumKnapsacks); assignment != nil {
		input.Assignment = assignment
		return true, nil
	}

	return false, nil
}
