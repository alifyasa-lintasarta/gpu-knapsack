package main

import "fmt"


type SchedulingInput struct {
	Items      []PodSpec
	NumGPUs    int
	Quota      map[string]int
	GPUFamily  GPUFamily
	Assignment []int
}

func validateSchedulingInput(input *SchedulingInput) error {
	if len(input.Items) == 0 {
		return fmt.Errorf("no items to assign")
	}
	if input.NumGPUs <= 0 {
		return fmt.Errorf("number of GPUs must be positive")
	}
	if len(input.GPUFamily.Capacity) == 0 {
		return fmt.Errorf("GPU capacity must be specified")
	}
	if input.GPUFamily.Mappings == nil {
		return fmt.Errorf("mappings cannot be nil")
	}

	for _, item := range input.Items {
		if _, exists := input.GPUFamily.Mappings[item.Type]; !exists {
			return fmt.Errorf("no mapping found for item type %s", item.Type)
		}
	}

	return nil
}

