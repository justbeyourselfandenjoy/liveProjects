package main

import (
	"justbeyourselfandenjoy/hash_tables/chaining"
	linearprobing "justbeyourselfandenjoy/hash_tables/linear_probing"
)

func main() {
	chaining.ChainingRun()
	linearprobing.LinearProbingRun()
}
