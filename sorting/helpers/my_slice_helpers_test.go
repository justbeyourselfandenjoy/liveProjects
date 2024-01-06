package my_slice_helpers

import (
	"slices"
	"testing"
)

func TestMakeRandomSlice(t *testing.T) {
	slice := MakeRandomSlice(100, 200)
	if len(slice) != 100 {
		t.Fatalf("Wrong slice size: %v instead of %v", len(slice), 100)
	}
	v := slices.Max(slice)
	if v > 200 {
		t.Fatalf("Wrong slice max value: %v", v)
	}
}
