package bubble

import (
	"fmt"
	my_slice_helpers "justbeyourselfandenjoy/sorting/helpers"
)

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

	slice := my_slice_helpers.MakeRandomSlice(numItems, max)
	my_slice_helpers.PrintSlice(slice, 40)
	fmt.Println()

	// Sort and display the result.
	bubbleSort(slice)
	my_slice_helpers.PrintSlice(slice, 40)
	fmt.Println()

	// Verify that it's sorted.
	my_slice_helpers.CheckSorted(slice)
}
