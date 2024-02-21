package main

import (
	"justbeyourselfandenjoy/hash_tables/chaining"
	doublehashing "justbeyourselfandenjoy/hash_tables/double_hashing"
	linearprobing "justbeyourselfandenjoy/hash_tables/linear_probing"
	linearremoving "justbeyourselfandenjoy/hash_tables/linear_removing"
	quadraticprobing "justbeyourselfandenjoy/hash_tables/quadratic_probing"
)

func main() {
	chaining.ChainingRun()
	linearprobing.LinearProbingRun()
	linearremoving.LinearProbingRemovingRun()
	quadraticprobing.QuadraticProbingRemovingRun()
	doublehashing.DoubleHashingRun()
}
