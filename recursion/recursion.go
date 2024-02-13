package main

import (
	"justbeyourselfandenjoy/recursion/factorial"
	"justbeyourselfandenjoy/recursion/fibonacci"
	knightstour "justbeyourselfandenjoy/recursion/knights_tour"
	nqueens "justbeyourselfandenjoy/recursion/n_queens"
)

func main() {
	factorial.RunFactorial()
	fibonacci.FibonacciRun()
	fibonacci.DynamicFibonacciRun()
	knightstour.KnightsTourRun()
	nqueens.NQueensRun()
}
