package main

import "fmt"

func main() {
	filename := parseArgs()
	cfg := loadConfig(filename)

	printConfig(cfg)

	podItems := createPodItems(cfg)

	input := &SchedulingInput{
		Items:       podItems,
		GPUCapacity: cfg.GPU.Capacity,
		NumGPUs:     cfg.GPU.Number,
		Mappings:    cfg.GPU.Mappings,
	}

	success := runSimulation(input)
	if !success {
		fmt.Println("No valid assignment found.")
		return
	}

	printFinalSummary(input)
}
