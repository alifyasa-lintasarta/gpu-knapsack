package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	GPU struct {
		Number   int              `yaml:"number"`
		Capacity []int            `yaml:"capacity"`
		Mappings map[string][]int `yaml:"mappings"`
	} `yaml:"gpu"`
	Pods map[string]int `yaml:"pods"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <config.yaml>\n", os.Args[0])
	}
	filename := os.Args[1]

	// Read the YAML file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", filename, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Failed to parse YAML: %v", err)
	}

	// Build the gpuRequests slice
	gpuRequests := []string{}
	for gpuType, count := range cfg.Pods {
		for i := 0; i < count; i++ {
			gpuRequests = append(gpuRequests, gpuType)
		}
	}

	// Build itemWeights from requests and mappings
	itemWeights := [][]int{}
	for _, gpu := range gpuRequests {
		weights, ok := cfg.GPU.Mappings[gpu]
		if !ok {
			log.Fatalf("No mapping found for GPU type %s", gpu)
		}
		itemWeights = append(itemWeights, weights)
	}

	assignment := assignItemsToKnapsacks(itemWeights, cfg.GPU.Capacity, cfg.GPU.Number)
	if assignment == nil {
		fmt.Println("No valid assignment found.")
		return
	}

	fmt.Println("Valid assignment found:")
	knapsackToItems := make(map[int][]int)
	for itemIndex, knapsackIndex := range assignment {
		knapsackToItems[knapsackIndex] = append(knapsackToItems[knapsackIndex], itemIndex)
	}

	// Print current assignment
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

	findAllPossibleCombinations(cfg)
}
