package binary

import (
	"fmt"
	my_slice_helpers "justbeyourselfandenjoy/sorting/helpers"
	"justbeyourselfandenjoy/sorting/quick"
	"math"
	"strconv"
)

func binarySearch(slice []int, target int) (index, numTests int) {
	index, numTests = -1, 0
	L, R := 0, len(slice)-1
	for L <= R {
		numTests++
		m := (int)(math.Floor(float64((L + R) / 2)))
		if slice[m] < target {
			L = m + 1
		} else if slice[m] > target {
			R = m - 1
		} else {
			index = m
			return
		}
	}
	return
}

func BinarySearchRun() {
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
	quick.QuickSort(slice)
	my_slice_helpers.PrintSlice(slice, 40)
	fmt.Println()

	for {
		var targetString string
		fmt.Printf("Target search: ")
		fmt.Scanln(&targetString)

		// If the target string is blank, break out of the loop.
		if len(targetString) == 0 {
			break
		}

		// Convert to int and find it.
		target, _ := strconv.Atoi(targetString)
		fmt.Printf("Target: %v\n", targetString)

		ind, numTests := binarySearch(slice, target)
		if ind >= 0 {
			fmt.Printf("values[%v] = %v, %v tests\n", ind, target, numTests)
		} else {
			fmt.Printf("Target: %v not found, %v tests\n", target, numTests)
		}
	}
}
