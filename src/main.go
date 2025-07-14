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
	Pods []ConfigItem `yaml:"pods"`
}

type ConfigItem struct {
	Type       string `yaml:"type"`
	AddTime    int    `yaml:"addTime"`
	RemoveTime *int   `yaml:"removeTime,omitempty"`
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

func printConfig(cfg Config) {
	fmt.Printf("GPUs: %d\n", cfg.GPU.Number)
	fmt.Printf("GPU Capacities: %v\n", cfg.GPU.Capacity)
	fmt.Printf("Events: %d\n", len(cfg.Pods))
	for _, event := range cfg.Pods {
		fmt.Printf("  %s (addTime=%d)\n", event.Type, event.AddTime)
	}
	fmt.Println()
}

func buildAllPods(cfg Config) []AssignmentItem {
	allPods := make([]AssignmentItem, len(cfg.Pods))
	for i, event := range cfg.Pods {
		allPods[i] = AssignmentItem{
			Type:           event.Type,
			AssignmentTime: event.AddTime,
			RemoveTime:     event.RemoveTime,
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

	// Run event-driven simulation
	success := runSimulation(input)
	if !success {
		fmt.Println("No valid assignment found.")
		return
	}

	// Print final summary with just types
	printFinalSummary(input)
}

func runSimulation(input *AssignmentInput) bool {
	if input == nil {
		return false
	}

	// Validate input
	if err := validateAssignmentInput(input); err != nil {
		return false
	}

	// Build item weights from mappings
	itemWeights := make([][]int, len(input.Items))
	for i, item := range input.Items {
		weights, exists := input.Mappings[item.Type]
		if !exists {
			return false
		}
		itemWeights[i] = weights
	}

	// Run the simulation with state tracking
	return runEventLoop(input.Items, itemWeights, input.KnapsackCapacity, input.NumKnapsacks, input)
}

func printFinalSummary(input *AssignmentInput) {
	fmt.Println("\nFinal GPU Assignment:")

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
				fmt.Print(podType)
			}
		}
		fmt.Println()
	}
}
