package bubble

import (
	my_slice_utils "justbeyourselfandenjoy/sorting/utils"
	"slices"
	"testing"
)

func TestBubbleSort(t *testing.T) {
	slice_origin_1 := my_slice_utils.MakeRandomSlice(100, 200)
	slice_origin_2 := make([]int, 100)
	copy(slice_origin_2, slice_origin_1)
	bubbleSort(slice_origin_1)
	slices.Sort(slice_origin_2)
	if slices.Compare(slice_origin_1, slice_origin_2) != 0 {
		t.Fatalf("Broken sorting detected: %v %v", slice_origin_1, slice_origin_2)
	}
}
