package main

import (
	"log"
	"os"
	"sort"

	"gopkg.in/yaml.v2"
)

type GPUFamily struct {
	Capacity []int            `yaml:"capacity"`
	Mappings map[string][]int `yaml:"mappings"`
}

type FitResult struct {
	Success    bool              
	Assignment []int             
	Layout     [][]int           
	QuotaUsage map[string]int    
}

type Config struct {
	GPU struct {
		Number int              `yaml:"number"`
		Quota  map[string]int   `yaml:"quota"`
		GPUFamily               `yaml:",inline"`
	} `yaml:"gpu"`
	Pods []PodSpec `yaml:"pods"`
}

type PodSpec struct {
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

func (gf *GPUFamily) buildEventTimeline(pods []PodSpec, itemWeights [][]int) []Event {
	var events []Event

	for i, pod := range pods {
		events = append(events, Event{
			Time:      pod.AddTime,
			Type:      "assign",
			ItemIndex: i,
			Weight:    itemWeights[i],
		})

		if pod.RemoveTime != nil {
			events = append(events, Event{
				Time:      *pod.RemoveTime,
				Type:      "remove",
				ItemIndex: i,
				Weight:    itemWeights[i],
			})
		}
	}

	sort.Slice(events, func(i, j int) bool {
		if events[i].Time == events[j].Time {
			return events[i].Type == "remove" && events[j].Type == "assign"
		}
		return events[i].Time < events[j].Time
	})

	return events
}

func (gf *GPUFamily) calculateFit(numGPUs int, quota map[string]int, pods []PodSpec) FitResult {
	itemWeights := make([][]int, len(pods))
	for i, pod := range pods {
		weights, exists := gf.Mappings[pod.Type]
		if !exists {
			return FitResult{
				Success:    false,
				Assignment: make([]int, len(pods)),
				Layout:     make([][]int, numGPUs),
				QuotaUsage: make(map[string]int),
			}
		}
		itemWeights[i] = weights
	}
	
	podEvents := gf.buildEventTimeline(pods, itemWeights)
	
	quotaUsage := make(map[string]int)
	
	capacityDimensions := len(gf.Capacity)
	usage := make([][]int, numGPUs)
	for i := range usage {
		usage[i] = make([]int, capacityDimensions)
	}
	
	assignment := make([]int, len(pods))
	for i := range assignment {
		assignment[i] = -1
	}
	
	for _, event := range podEvents {
		if event.Type == "assign" {
			podType := pods[event.ItemIndex].Type
			
			if quota[podType] > quotaUsage[podType] {
				quotaUsage[podType]++
				assignment[event.ItemIndex] = -1
				continue
			}
			
			placed := false
			for gpuIdx := 0; gpuIdx < numGPUs && !placed; gpuIdx++ {
				canFit := true
				for d := 0; d < capacityDimensions; d++ {
					if usage[gpuIdx][d]+event.Weight[d] > gf.Capacity[d] {
						canFit = false
						break
					}
				}
				
				if canFit {
					for d := 0; d < capacityDimensions; d++ {
						usage[gpuIdx][d] += event.Weight[d]
					}
					assignment[event.ItemIndex] = gpuIdx
					placed = true
				}
			}
			
			if !placed {
				return FitResult{
					Success:    false,
					Assignment: assignment,
					Layout:     usage,
					QuotaUsage: quotaUsage,
				}
			}
			
		} else if event.Type == "remove" {
			gpuIndex := assignment[event.ItemIndex]
			if gpuIndex != -1 {
				for d := 0; d < capacityDimensions; d++ {
					usage[gpuIndex][d] -= event.Weight[d]
				}
				assignment[event.ItemIndex] = -1
			}
		}
	}
	
	return FitResult{
		Success:    true,
		Assignment: assignment,
		Layout:     usage,
		QuotaUsage: quotaUsage,
	}
}