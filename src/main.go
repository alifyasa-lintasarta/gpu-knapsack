package main

import "fmt"

func main() {
	filename := parseArgs()
	cfg := loadConfig(filename)

	printConfig(cfg)

	podItems := createPodItems(cfg)
	
	itemWeights := make([][]int, len(podItems))
	for i, item := range podItems {
		weights, exists := cfg.GPU.GPUFamily.Mappings[item.Type]
		if !exists {
			fmt.Printf("No mapping found for item type %s\n", item.Type)
			return
		}
		itemWeights[i] = weights
	}

	events := buildEventTimeline(podItems, itemWeights)
	result := cfg.GPU.GPUFamily.calculateFit(cfg.GPU.Number, cfg.GPU.Quota, events, podItems)
	
	fmt.Printf("\nCalculateFit Result:\n")
	fmt.Printf("Success: %t\n", result.Success)
	fmt.Printf("Quota Usage: %v\n", result.QuotaUsage)
	
	if result.Success {
		fmt.Printf("Assignment: %v\n", result.Assignment)
		fmt.Printf("Final Layout: %v\n", result.Layout)
	} else {
		fmt.Println("No valid assignment found.")
	}
}
