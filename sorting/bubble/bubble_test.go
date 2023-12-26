package bubble

import (
	"slices"
	"testing"
)

var MakeRandomSlice = makeRandomSlice
var BubbleSort = bubbleSort

func TestMakeRandomSlice(t *testing.T) {
	slice := makeRandomSlice(100, 200)
	if len(slice) != 100 {
		t.Fatalf("Wrong slice size: %v instead of %v", len(slice), 100)
	}
	v := slices.Max(slice)
	if v > 200 {
		t.Fatalf("Wrong slice max value: %v", v)
	}
}

func TestBubbleSort(t *testing.T) {
	slice_origin_1 := makeRandomSlice(100, 200)
	slice_origin_2 := make([]int, 100)
	copy(slice_origin_2, slice_origin_1)
	bubbleSort(slice_origin_1)
	slices.Sort(slice_origin_2)
	if slices.Compare(slice_origin_1, slice_origin_2) != 0 {
		t.Fatalf("Broken sorting detected: %v %v", slice_origin_1, slice_origin_2)
	}
}
