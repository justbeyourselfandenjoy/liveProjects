package my_slice_helpers

import (
	"fmt"
	"math/rand"
	"time"
)

var MakeRandomSlice = makeRandomSlice
var PrintSlice = printSlice
var CheckSorted = checkSorted

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
