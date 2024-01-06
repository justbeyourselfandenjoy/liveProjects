package main

import (
	"justbeyourselfandenjoy/sorting/bubble"
	"justbeyourselfandenjoy/sorting/counting"
	"justbeyourselfandenjoy/sorting/quick"
)

func main() {
	bubble.BubbleSortRun()
	quick.QuickSortRun()
	counting.CountingSortRun()
	counting.CountingSortStructRun()
}
