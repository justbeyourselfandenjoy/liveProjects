package counting

import (
	"fmt"
	my_slice_helpers "justbeyourselfandenjoy/sorting/helpers"
	"math/rand"
	"slices"
	"time"
)

type Customer struct {
	id           string
	numPurchases int
}

func makeRandomSliceStruct(numItems, max int) []Customer {
	slice := make([]Customer, numItems)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range slice {
		slice[i].id = "C" + fmt.Sprint(i)
		slice[i].numPurchases = random.Intn(max)
	}
	return slice
}

func printSliceStruct(slice []Customer, numItems int) {
	fmt.Printf("%v", slice[:min(len(slice), numItems)])
}

func checkSortedStruct(slice []Customer) {
	prev := slice[0]
	for i := range slice {
		if slice[i].numPurchases < prev.numPurchases {
			fmt.Println("The slice is NOT sorted!")
			return
		} else {
			prev = slice[i]
		}
	}
	fmt.Println("The slice is sorted")
}
func countingSortStruct(slice []Customer) []Customer {
	max_val := 0
	for i := range slice {
		if slice[i].numPurchases > max_val {
			max_val = slice[i].numPurchases
		}
	}
	C := make([]int, max_val+1)
	for _, v := range slice {
		C[v.numPurchases]++
	}
	for i := 1; i < len(C); i++ {
		C[i] += C[i-1]
	}
	B := make([]Customer, len(slice))
	for i := len(slice) - 1; i >= 0; i-- {
		C[slice[i].numPurchases]--
		B[C[slice[i].numPurchases]] = slice[i]
	}
	return B
}

func countingSort(slice []int) []int {
	C := make([]int, slices.Max(slice)+1)
	for _, v := range slice {
		C[v]++
	}
	for i := 1; i < len(C); i++ {
		C[i] += C[i-1]
	}
	B := make([]int, len(slice))
	for i := len(slice) - 1; i >= 0; i-- {
		C[slice[i]]--
		B[C[slice[i]]] = slice[i]
	}
	return B
}

func CountingSortRun() {
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
	sorted := countingSort(slice)
	my_slice_helpers.PrintSlice(sorted, 40)
	fmt.Println()

	// Verify that it's sorted.
	my_slice_helpers.CheckSorted(sorted)
}

func CountingSortStructRun() {
	// Get the number of items and maximum item value.
	var numItems, max int
	fmt.Printf("# Items: ")
	fmt.Scanln(&numItems)
	fmt.Printf("Max: ")
	fmt.Scanln(&max)

	// Make and display the unsorted slice.
	slice := makeRandomSliceStruct(numItems, max)
	printSliceStruct(slice, 40)
	fmt.Println()

	// Sort and display the result.
	sorted := countingSortStruct(slice)
	printSliceStruct(sorted, 40)
	fmt.Println()

	// Verify that it's sorted.
	checkSortedStruct(sorted)
}
