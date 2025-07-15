package main

import "fmt"

func main() {
	filename := parseArgs()
	cfg := loadConfig(filename)

	printConfig(cfg)

	result := cfg.GPU.GPUFamily.calculateFit(cfg.GPU.Number, cfg.GPU.Quota, cfg.Pods)
	
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
