package quick

import (
	"fmt"
	my_slice_helpers "justbeyourselfandenjoy/sorting/helpers"
)

func partition(slice []int) int {
	hi, i := len(slice)-1, -1
	pivot := slice[hi]

	for j := range slice {
		if j == hi {
			break
		}
		if slice[j] <= pivot {
			i++
			slice[i], slice[j] = slice[j], slice[i]
		}
	}
	i++
	slice[i], slice[hi] = slice[hi], slice[i]
	return i
}

func quickSort(slice []int) {
	if len(slice) < 2 {
		return
	}

	p := partition(slice)
	quickSort(slice[:p])
	quickSort(slice[p+1:])
}

func QuickSortRun() {
	// Get the number of items and maximum item value.
	var numItems, max int
	fmt.Printf("# Items: ")
	fmt.Scanln(&numItems)
	fmt.Printf("Max: ")
	fmt.Scanln(&max)

	// Make and display the unsorted slice.
	slice := my_slice_helpers.MakeRandomSlice(numItems, max)
	my_slice_helpers.PrintSlice(slice, 40)
	fmt.Println()

	// Sort and display the result.
	quickSort(slice)
	my_slice_helpers.PrintSlice(slice, 40)
	fmt.Println()

	// Verify that it's sorted.
	my_slice_helpers.CheckSorted(slice)
}
