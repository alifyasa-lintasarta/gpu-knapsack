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
	Items []ConfigItem `yaml:"items"`
}

type ConfigItem struct {
	Type string `yaml:"type"`
	Time int    `yaml:"time"`
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
	fmt.Printf("Items: %d\n", len(cfg.Items))
	for _, item := range cfg.Items {
		fmt.Printf("  %s (t=%d)\n", item.Type, item.Time)
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

func buildAllPods(cfg Config) []AssignmentItem {
	allPods := make([]AssignmentItem, len(cfg.Items))
	for i, item := range cfg.Items {
		allPods[i] = AssignmentItem{
			Type:           item.Type,
			AssignmentTime: item.Time,
		}
	}
	return allPods
}

func main() {
	filename := parseArgs()
	cfg := loadConfig(filename)

	printConfig(cfg)

	// Build unified list of all pods with timestamps
	allPods := buildAllPods(cfg)

	// Create assignment input
	input := &AssignmentInput{
		Items:            allPods,
		KnapsackCapacity: cfg.GPU.Capacity,
		NumKnapsacks:     cfg.GPU.Number,
		Mappings:         cfg.GPU.Mappings,
	}

	// Assign all pods using timestamp-based algorithm
	success, err := AssignItems(input)
	if err != nil {
		log.Fatalf("Assignment error: %v", err)
	}
	if !success {
		fmt.Println("No valid assignment found.")
		return
	}

	// Print results
	printNewAssignment(input)
}

func printNewAssignment(input *AssignmentInput) {
	fmt.Println("GPU Assignment:")
	
	// Group items by knapsack
	knapsackToItems := make(map[int][]int)
	for itemIndex, knapsackIndex := range input.Assignment {
		knapsackToItems[knapsackIndex] = append(knapsackToItems[knapsackIndex], itemIndex)
	}
	
	for k := 0; k < input.NumKnapsacks; k++ {
		items := knapsackToItems[k]
		fmt.Printf("GPU %d: ", k)
		
		if len(items) == 0 {
			fmt.Print("(empty)")
		} else {
			for i, itemIndex := range items {
				if i > 0 {
					fmt.Print(", ")
				}
				podType := input.Items[itemIndex].Type
				podTime := input.Items[itemIndex].AssignmentTime
				fmt.Printf("%s (t=%d)", podType, podTime)
			}
		}
		fmt.Println()
	}
}
