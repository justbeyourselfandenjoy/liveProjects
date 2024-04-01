package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

const numItems = 20 // A reasonable value for exhaustive search.

const minValue = 1
const maxValue = 10
const minWeight = 4
const maxWeight = 10

var allowedWeight int

type Item struct {
	id, blockedBy int
	blockList     []int //Other items that this one blocks.
	value, weight int
	isSelected    bool
}

// Make some random items.
func makeItems(numItems, minValue, maxValue, minWeight, maxWeight int) []Item {
	// Initialize a pseudorandom number generator.
	//random := rand.New(rand.NewSource(time.Now().UnixNano())) // Initialize with a changing seed
	random := rand.New(rand.NewSource(1337)) // Initialize with a fixed seed

	items := make([]Item, numItems)
	for i := 0; i < numItems; i++ {
		items[i] = Item{
			i, -1, nil,
			random.Intn(maxValue-minValue+1) + minValue,
			random.Intn(maxWeight-minWeight+1) + minWeight,
			false}
	}
	return items
}

// Return a copy of the items slice.
func copyItems(items []Item) []Item {
	newItems := make([]Item, len(items))
	copy(newItems, items)
	return newItems
}

// Return the total value of the items.
// If addAll is false, only add up the selected items.
func sumValues(items []Item, addAll bool) int {
	total := 0
	for i := 0; i < len(items); i++ {
		if addAll || items[i].isSelected {
			total += items[i].value
		}
	}
	return total
}

// Return the total weight of the items.
// If addAll is false, only add up the selected items.
func sumWeights(items []Item, addAll bool) int {
	total := 0
	for i := 0; i < len(items); i++ {
		if addAll || items[i].isSelected {
			total += items[i].weight
		}
	}
	return total
}

// Return the value of this solution.
// If the solution is too heavy, return -1 so we prefer an empty solution.
func solutionValue(items []Item, allowedWeight int) int {
	// If the solution's total weight > allowedWeight,
	// return 0 so we won't use this solution.
	if sumWeights(items, false) > allowedWeight {
		return -1
	}

	// Return the sum of the selected values.
	return sumValues(items, false)
}

// Print the selected items.
func printSelected(items []Item) {
	numPrinted := 0
	for i, item := range items {
		if item.isSelected {
			fmt.Printf("%d(%d, %d) ", i, item.value, item.weight)
		}
		numPrinted += 1
		if numPrinted > 100 {
			fmt.Println("...")
			return
		}
	}
	fmt.Println()
}

// Run the algorithm. Display the elapsed time and solution.
func runAlgorithm(alg func([]Item, int) ([]Item, int, int), items []Item, allowedWeight int) {
	// Copy the items so the run isn't influenced by a previous run.
	testItems := copyItems(items)

	start := time.Now()

	// Run the algorithm.
	solution, totalValue, functionCalls := alg(testItems, allowedWeight)

	elapsed := time.Since(start)

	fmt.Printf("Elapsed: %f\n", elapsed.Seconds())
	printSelected(solution)
	fmt.Printf("Value: %d, Weight: %d, Calls: %d\n",
		totalValue, sumWeights(solution, false), functionCalls)
	fmt.Println()
}

// Recursively assign values in or out of the solution.
// Return the best assignment, value of that assignment,
// and the number of function calls we made.
func exhaustiveSearch(items []Item, allowedWeight int) ([]Item, int, int) {
	return doExhaustiveSearch(items, allowedWeight, 0)
}

func doExhaustiveSearch(items []Item, allowedWeight, nextIndex int) ([]Item, int, int) {

	if nextIndex >= len(items) {
		itemsCopy := copyItems(items)
		return itemsCopy, solutionValue(itemsCopy, allowedWeight), 1
	}

	items[nextIndex].isSelected = true
	solutionWItem, totalValueWItem, callsNumberWItem := doExhaustiveSearch(items, allowedWeight, nextIndex+1)

	items[nextIndex].isSelected = false
	solutionWOItem, totalValueWOItem, callsNumberWOItem := doExhaustiveSearch(items, allowedWeight, nextIndex+1)

	if totalValueWItem >= totalValueWOItem {
		return solutionWItem, totalValueWItem, callsNumberWItem + callsNumberWOItem + 1
	}

	return solutionWOItem, totalValueWOItem, callsNumberWItem + callsNumberWOItem + 1
}

func branchAndBound(items []Item, allowedWeight int) ([]Item, int, int) {
	bestValue := 0
	currentValue := 0
	currentWeight := 0
	remainingValue := sumValues(items, true)

	return doBranchAndBound(items, allowedWeight, 0, bestValue, currentValue, currentWeight, remainingValue)
}

func doBranchAndBound(items []Item, allowedWeight, nextIndex, bestValue, currentValue, currentWeight, remainingValue int) ([]Item, int, int) {

	if nextIndex >= len(items) {
		return copyItems(items), currentValue, 1
	}

	if currentValue+remainingValue <= bestValue {
		return nil, currentValue, 1
	}

	solutionWItem, totalValueWItem, callsNumberWItem := []Item(nil), 0, 1
	if currentWeight+items[nextIndex].weight <= allowedWeight {
		items[nextIndex].isSelected = true
		solutionWItem, totalValueWItem, callsNumberWItem = doBranchAndBound(items, allowedWeight, nextIndex+1, bestValue, currentValue+items[nextIndex].value, currentWeight+items[nextIndex].weight, remainingValue-items[nextIndex].value)
		if totalValueWItem > bestValue {
			bestValue = totalValueWItem
		}
	}

	solutionWOItem, totalValueWOItem, callsNumberWOItem := []Item(nil), 0, 1
	if currentValue+remainingValue-items[nextIndex].value > bestValue {
		items[nextIndex].isSelected = false
		solutionWOItem, totalValueWOItem, callsNumberWOItem = doBranchAndBound(items, allowedWeight, nextIndex+1, bestValue, currentValue, currentWeight, remainingValue-items[nextIndex].value)
	}

	if totalValueWItem >= totalValueWOItem {
		return solutionWItem, totalValueWItem, callsNumberWItem + callsNumberWOItem + 1
	}

	return solutionWOItem, totalValueWOItem, callsNumberWItem + callsNumberWOItem + 1
}

// Build the items' block lists.
func makeBlockLists(items []Item) {
	for i := range items {
		items[i].blockList = []int{}
		for j := range items {
			if i != j {
				if items[i].value >= items[j].value && items[i].weight <= items[j].weight {
					items[i].blockList = append(items[i].blockList, items[j].id)
				}
			}
		}
	}
}

// Block items on this item's blocks list.
func blockItems(source Item, items []Item) {
	for _, otherId := range source.blockList {
		if items[otherId].blockedBy < 0 {
			items[otherId].blockedBy = source.id
		}
	}
}

// Unblock items on this item's blocks list.
func unblockItems(source Item, items []Item) {
	for _, otherId := range source.blockList {
		if items[otherId].blockedBy == source.id {
			items[otherId].blockedBy = -1
		}
	}
}

func rodsTechnique(items []Item, allowedWeight int) ([]Item, int, int) {
	bestValue := 0
	currentValue := 0
	currentWeight := 0
	remainingValue := sumValues(items, true)
	makeBlockLists(items)

	return doRodsTechnique(items, allowedWeight, 0, bestValue, currentValue, currentWeight, remainingValue)
}

func doRodsTechnique(items []Item, allowedWeight, nextIndex, bestValue, currentValue, currentWeight, remainingValue int) ([]Item, int, int) {

	if nextIndex >= len(items) {
		return copyItems(items), currentValue, 1
	}

	if currentValue+remainingValue <= bestValue {
		return nil, currentValue, 1
	}

	solutionWItem, totalValueWItem, callsNumberWItem := []Item(nil), 0, 1
	if items[nextIndex].blockedBy < 0 {
		if currentWeight+items[nextIndex].weight <= allowedWeight {
			items[nextIndex].isSelected = true
			solutionWItem, totalValueWItem, callsNumberWItem = doRodsTechnique(items, allowedWeight, nextIndex+1, bestValue, currentValue+items[nextIndex].value, currentWeight+items[nextIndex].weight, remainingValue-items[nextIndex].value)
			if totalValueWItem > bestValue {
				bestValue = totalValueWItem
			}
		}
	}

	solutionWOItem, totalValueWOItem, callsNumberWOItem := []Item(nil), 0, 1
	if currentValue+remainingValue-items[nextIndex].value > bestValue {
		blockItems(items[nextIndex], items)
		items[nextIndex].isSelected = false
		solutionWOItem, totalValueWOItem, callsNumberWOItem = doRodsTechnique(items, allowedWeight, nextIndex+1, bestValue, currentValue, currentWeight, remainingValue-items[nextIndex].value)
		unblockItems(items[nextIndex], items)
	}

	if totalValueWItem >= totalValueWOItem {
		return solutionWItem, totalValueWItem, callsNumberWItem + callsNumberWOItem + 1
	}

	return solutionWOItem, totalValueWOItem, callsNumberWItem + callsNumberWOItem + 1
}

// Use Rod's technique sorted to find a solution.
// Return the best assignment, value of that assignment,
// and the number of function calls we made.
func rodsTechniqueSorted(items []Item, allowedWeight int) ([]Item, int, int) {
	makeBlockLists(items)

	// Sort so items with longer blocked lists come first.
	sort.Slice(items, func(i, j int) bool {
		return len(items[i].blockList) > len(items[j].blockList)
	})

	// Reset the items' IDs.
	for i := range items {
		items[i].id = i
	}

	// Rebuild the blocked lists with the new indices.
	makeBlockLists(items)

	bestValue := 0
	currentValue := 0
	currentWeight := 0
	remainingValue := sumValues(items, true)

	return doRodsTechnique(items, allowedWeight, 0,
		bestValue, currentValue, currentWeight, remainingValue)
}

func main() {
	//items := makeTestItems()
	items := makeItems(numItems, minValue, maxValue, minWeight, maxWeight)
	allowedWeight = sumWeights(items, true) / 2

	// Display basic parameters.
	fmt.Println("*** Parameters ***")
	fmt.Printf("# items: %d\n", numItems)
	fmt.Printf("Total value: %d\n", sumValues(items, true))
	fmt.Printf("Total weight: %d\n", sumWeights(items, true))
	fmt.Printf("Allowed weight: %d\n", allowedWeight)
	fmt.Println()

	// Exhaustive search
	if numItems > 25 { // Only run exhaustive search if numItems <= 25.
		fmt.Println("Too many items for exhaustive search\n")
	} else {
		fmt.Println("*** Exhaustive Search ***")
		runAlgorithm(exhaustiveSearch, items, allowedWeight)
	}

	// Branch and bound
	if numItems > 45 { // Only run branch and bound if numItems <= 45.
		fmt.Println("Too many items for branch and bound\n")
	} else {
		fmt.Println("*** Branch and Bound ***")
		runAlgorithm(branchAndBound, items, allowedWeight)
	}

	// Rod's technique
	if numItems > 85 { // Only use Rod's technique if numItems <= 85.
		fmt.Println("Too many items for Rod's technique\n")
	} else {
		fmt.Println("*** Rod's technique ***")
		runAlgorithm(rodsTechnique, items, allowedWeight)
	}

	// Rod's technique sorted
	if numItems > 350 { // Only use Rod's technique sorted if numItems <= 350.
		fmt.Println("Too many items for Rod's technique sorted\n")
	} else {
		fmt.Println("*** Rod's technique sorted ***")
		runAlgorithm(rodsTechniqueSorted, items, allowedWeight)
	}
}
