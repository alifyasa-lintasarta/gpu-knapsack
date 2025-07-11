package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	GPU struct {
		Number   int              `yaml:"number"`
		Capacity []int            `yaml:"capacity"`
		Mappings map[string][]int `yaml:"mappings"`
	} `yaml:"gpu"`
	Pods map[string]int `yaml:"pods"`
}

func parseArgs() string {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <config.yaml>\n", os.Args[0])
	}
	return os.Args[1]
}

func loadConfig(filename string) Config {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", filename, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Failed to parse YAML: %v", err)
	}
	return cfg
}

func buildGPURequests(pods map[string]int) []string {
	totalPods := 0
	for _, count := range pods {
		totalPods += count
	}

	gpuRequests := make([]string, 0, totalPods)
	for gpuType, count := range pods {
		for i := 0; i < count; i++ {
			gpuRequests = append(gpuRequests, gpuType)
		}
	}
	return gpuRequests
}

func buildItemWeights(gpuRequests []string, mappings map[string][]int) [][]int {
	itemWeights := make([][]int, 0, len(gpuRequests))
	for _, gpu := range gpuRequests {
		weights, ok := mappings[gpu]
		if !ok {
			log.Fatalf("No mapping found for GPU type %s", gpu)
		}
		itemWeights = append(itemWeights, weights)
	}
	return itemWeights
}

func groupItemsByKnapsack(assignment []int) map[int][]int {
	knapsackToItems := make(map[int][]int)
	for itemIndex, knapsackIndex := range assignment {
		knapsackToItems[knapsackIndex] = append(knapsackToItems[knapsackIndex], itemIndex)
	}
	return knapsackToItems
}

func printConfig(cfg Config) {
	fmt.Printf("GPUs: %d\n", cfg.GPU.Number)
	fmt.Printf("GPU Capacities: %v\n", cfg.GPU.Capacity)
	fmt.Println("Requested Pods:")
	for podType, count := range cfg.Pods {
		fmt.Printf("  %s: %d\n", podType, count)
	}
	fmt.Println()
}

func printAssignment(knapsackToItems map[int][]int, gpuRequests []string) {
	fmt.Println("GPU Assignment:")
	for k := 0; k < len(knapsackToItems); k++ {
		items := knapsackToItems[k]
		fmt.Printf("GPU %d: ", k)
		for i, itemIndex := range items {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(gpuRequests[itemIndex])
		}
		fmt.Println()
	}
}

func main() {
	filename := parseArgs()
	cfg := loadConfig(filename)
	printConfig(cfg)
	gpuRequests := buildGPURequests(cfg.Pods)
	itemWeights := buildItemWeights(gpuRequests, cfg.GPU.Mappings)

	assignment := assignItemsToKnapsacks(itemWeights, cfg.GPU.Capacity, cfg.GPU.Number)
	if assignment == nil {
		fmt.Println("No valid assignment found.")
		return
	}

	knapsackToItems := groupItemsByKnapsack(assignment)
	printAssignment(knapsackToItems, gpuRequests)

	maximalCombinations := findAllPossibleCombinations(cfg)
	printMaximalCombinations(maximalCombinations)
}
