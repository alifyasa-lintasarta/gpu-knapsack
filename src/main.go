package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	GPU struct {
		Number       int              `yaml:"number"`
		Capacity     []int            `yaml:"capacity"`
		Mappings     map[string][]int `yaml:"mappings"`
		InitialState map[int][]string `yaml:"initialState,omitempty"`
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

func printAssignmentWithInitial(knapsackToItems map[int][]int, allPods []string, initialState map[int][]string) {
	fmt.Println("GPU Assignment:")

	// Count initial pods per GPU
	initialCounts := make(map[int]int)
	for gpuIndex, pods := range initialState {
		initialCounts[gpuIndex] = len(pods)
	}

	for k := 0; k < len(knapsackToItems); k++ {
		items := knapsackToItems[k]
		fmt.Printf("GPU %d: ", k)

		itemCount := 0
		initialCount := initialCounts[k]

		for i, itemIndex := range items {
			if i > 0 {
				fmt.Print(", ")
			}

			podName := allPods[itemIndex]
			if itemCount < initialCount {
				// This is an existing pod from initial state
				fmt.Printf("%s (existing)", podName)
			} else {
				// This is a newly assigned pod
				fmt.Printf("%s (new)", podName)
			}
			itemCount++
		}

		if len(items) == 0 {
			fmt.Print("(empty)")
		}
		fmt.Println()
	}
}

func main() {
	filename := parseArgs()
	cfg := loadConfig(filename)

	// Validate initial state if present
	if err := validateInitialState(cfg); err != nil {
		log.Fatalf("Invalid initial state: %v", err)
	}

	printConfig(cfg)

	// Compute initial usage from configuration
	initialUsage := computeInitialUsage(cfg.GPU.InitialState, cfg.GPU.Mappings, cfg.GPU.Number, len(cfg.GPU.Capacity))

	// Build requests for new pods to be assigned
	gpuRequests := buildGPURequests(cfg.Pods)
	itemWeights := buildItemWeights(gpuRequests, cfg.GPU.Mappings)

	// Assign new pods using the initial state
	assignment := assignItemsToKnapsacksWithInitial(itemWeights, cfg.GPU.Capacity, cfg.GPU.Number, initialUsage)
	if assignment == nil {
		fmt.Println("No valid assignment found for new pods.")
		return
	}

	// Merge initial state with new assignments for display
	allAssignments, allPods := mergeAssignments(cfg.GPU.InitialState, assignment, gpuRequests)
	knapsackToItems := groupItemsByKnapsack(allAssignments)
	printAssignmentWithInitial(knapsackToItems, allPods, cfg.GPU.InitialState)
}
