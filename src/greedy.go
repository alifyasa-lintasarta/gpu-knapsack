package main

func firstFitDecreasingWithInitial(sortedItems []Item, knapsackCapacity []int, numKnapsacks int, initialUsage [][]int) []int {
	numDimensions := len(knapsackCapacity)
	usage := make([][]int, numKnapsacks)
	for i := range usage {
		usage[i] = make([]int, numDimensions)
		// Initialize with existing usage
		for d := 0; d < numDimensions; d++ {
			if i < len(initialUsage) && d < len(initialUsage[i]) {
				usage[i][d] = initialUsage[i][d]
			}
		}
	}
	assignment := make([]int, len(sortedItems))
	for i := range assignment {
		assignment[i] = -1
	}

	for _, item := range sortedItems {
		placed := false
		for k := 0; k < numKnapsacks && !placed; k++ {
			canFit := true
			for d := 0; d < numDimensions; d++ {
				if usage[k][d]+item.Weight[d] > knapsackCapacity[d] {
					canFit = false
					break
				}
			}
			if canFit {
				for d := 0; d < numDimensions; d++ {
					usage[k][d] += item.Weight[d]
				}
				assignment[item.Index] = k
				placed = true
			}
		}
		if !placed {
			return nil
		}
	}
	return assignment
}

func tryGreedyAssignmentWithInitial(itemWeights [][]int, knapsackCapacity []int, numKnapsacks int, initialUsage [][]int) []int {
	sortedItems := sortItemsByWeight(itemWeights)
	return firstFitDecreasingWithInitial(sortedItems, knapsackCapacity, numKnapsacks, initialUsage)
}
