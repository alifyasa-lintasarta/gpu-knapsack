package main

import "fmt"

type PodItem struct {
	Type           string
	AssignmentTime int
	RemoveTime     *int
}

type SchedulingInput struct {
	Items       []PodItem
	GPUCapacity []int
	NumGPUs     int
	Mappings    map[string][]int
	Assignment  []int
}

func validateSchedulingInput(input *SchedulingInput) error {
	if len(input.Items) == 0 {
		return fmt.Errorf("no items to assign")
	}
	if input.NumGPUs <= 0 {
		return fmt.Errorf("number of GPUs must be positive")
	}
	if len(input.GPUCapacity) == 0 {
		return fmt.Errorf("GPU capacity must be specified")
	}
	if input.Mappings == nil {
		return fmt.Errorf("mappings cannot be nil")
	}

	for _, item := range input.Items {
		if _, exists := input.Mappings[item.Type]; !exists {
			return fmt.Errorf("no mapping found for item type %s", item.Type)
		}
	}

	return nil
}

func createPodItems(cfg Config) []PodItem {
	podItems := make([]PodItem, len(cfg.Pods))
	for i, podSpec := range cfg.Pods {
		podItems[i] = PodItem{
			Type:           podSpec.Type,
			AssignmentTime: podSpec.AddTime,
			RemoveTime:     podSpec.RemoveTime,
		}
	}
	return podItems
}