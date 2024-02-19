package main

import (
	"justbeyourselfandenjoy/hash_tables/chaining"
	linearprobing "justbeyourselfandenjoy/hash_tables/linear_probing"
	linearremoving "justbeyourselfandenjoy/hash_tables/linear_removing"
)

func main() {
	chaining.ChainingRun()
	linearprobing.LinearProbingRun()
	linearremoving.LinearProbingRemovingRun()
}
