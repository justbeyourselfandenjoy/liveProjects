package counting

import (
	"fmt"
	my_slice_utils "justbeyourselfandenjoy/sorting/utils"
)

func countingSort(slice []int, max int) []int {
	return make([]int, 0)
}

func CountingSortRun() {
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
	sorted := countingSort(slice, max)
	my_slice_utils.PrintSlice(sorted, 40)
	fmt.Println()

	// Verify that it's sorted.
	my_slice_utils.CheckSorted(slice)
}
