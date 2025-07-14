package main

import "fmt"

func runSimulation(input *SchedulingInput) bool {
	if input == nil {
		return false
	}

	if err := validateSchedulingInput(input); err != nil {
		return false
	}

	itemWeights := make([][]int, len(input.Items))
	for i, item := range input.Items {
		weights, exists := input.Mappings[item.Type]
		if !exists {
			return false
		}
		itemWeights[i] = weights
	}

	return processEvents(input.Items, itemWeights, input.GPUCapacity, input.NumGPUs, input)
}

func printConfig(cfg Config) {
	fmt.Printf("GPUs: %d\n", cfg.GPU.Number)
	fmt.Printf("GPU Capacities: %v\n", cfg.GPU.Capacity)
	fmt.Printf("Events: %d\n", len(cfg.Pods))
	for _, event := range cfg.Pods {
		if event.RemoveTime != nil {
			fmt.Printf("  %s (addTime=%d, removeTime=%d)\n", event.Type, event.AddTime, *event.RemoveTime)
		} else {
			fmt.Printf("  %s (addTime=%d)\n", event.Type, event.AddTime)
		}
	}
	fmt.Println()
}

func printFinalSummary(input *SchedulingInput) {
	fmt.Println("\nFinal GPU Assignment:")

	gpuToItems := make(map[int][]int)
	for itemIndex, gpuIndex := range input.Assignment {
		gpuToItems[gpuIndex] = append(gpuToItems[gpuIndex], itemIndex)
	}

	for gpuIdx := 0; gpuIdx < input.NumGPUs; gpuIdx++ {
		items := gpuToItems[gpuIdx]
		fmt.Printf("GPU %d: ", gpuIdx)

		if len(items) == 0 {
			fmt.Print("(empty)")
		} else {
			for i, itemIndex := range items {
				if i > 0 {
					fmt.Print(", ")
				}
				itemType := input.Items[itemIndex].Type
				fmt.Print(itemType)
			}
		}
		fmt.Println()
	}
}
