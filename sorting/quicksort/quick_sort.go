package quicksort

import (
	"fmt"
	my_slice_utils "justbeyourselfandenjoy/sorting/utils"
)

func partition(slice []int) int {
	lo := 0
	hi := len(slice) - 1
	pivot := slice[hi]
	i := lo - 1

	for j := range slice {
		if slice[j] <= pivot {
			i++
			slice[i], slice[j] = slice[j], slice[i]
		}
	}

	i++
	slice[i], slice[hi] = slice[hi], slice[i]
	return i
}

func quicksort(slice []int) {
	if len(slice) < 2 {
		return
	}

	p := partition(slice)
	quicksort(slice[:p-1])
	quicksort(slice[p+1:])
}

func QuickSortRun() {
	// Get the number of items and maximum item value.
	var numItems, max int
	fmt.Printf("# Items: ")
	fmt.Scanln(&numItems)
	fmt.Printf("Max: ")
	fmt.Scanln(&max)

	// Make and display the unsorted slice.
	slice := my_slice_utils.MakeRandomSlice(numItems, max)
	my_slice_utils.PrintSlice(slice, 40)
	fmt.Println()

	// Sort and display the result.
	quicksort(slice)
	my_slice_utils.PrintSlice(slice, 40)
	fmt.Println()

	// Verify that it's sorted.
	my_slice_utils.CheckSorted(slice)
}
