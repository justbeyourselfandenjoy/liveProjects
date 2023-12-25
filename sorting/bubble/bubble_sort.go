package bubble

import (
	"fmt"
	"math/rand"
	"time"
)

func makeRandomSlice(numItems, max int) []int {
	slice := make([]int, numItems)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range slice {
		slice[i] = random.Intn(max)
	}
	return slice
}

func printSlice(slice []int, numItems int) {
	fmt.Printf("%v", slice[:min(len(slice), numItems)])
}

func checkSorted(slice []int) {
	prev := slice[0]
	for i := range slice {
		if slice[i] < prev {
			fmt.Println("The slice is NOT sorted!")
			return
		} else {
			prev = slice[i]
		}
	}
	fmt.Println("The slice is sorted")
}

func bubbleSort(slice []int) {
	n := len(slice)
	for {
		swapped := false
		for i := 1; i < n; i++ {
			if slice[i-1] > slice[i] {
				slice[i-1], slice[i] = slice[i], slice[i-1]
				swapped = true
			}
		}
		if !swapped {
			break
		}
	}
}

func BubbleSortRun() {
	// Get the number of items and maximum item value.
	var numItems, max int
	fmt.Printf("# Items: ")
	fmt.Scanln(&numItems)
	fmt.Printf("Max: ")
	fmt.Scanln(&max)

	// Make and display an unsorted slice.
	slice := makeRandomSlice(numItems, max)
	printSlice(slice, 40)
	fmt.Println()

	// Sort and display the result.
	bubbleSort(slice)
	printSlice(slice, 40)

	// Verify that it's sorted.
	checkSorted(slice)
}
