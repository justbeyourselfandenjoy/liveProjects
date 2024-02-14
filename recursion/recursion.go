package main

import (
	"justbeyourselfandenjoy/recursion/factorial"
	"justbeyourselfandenjoy/recursion/fibonacci"
	knightstour "justbeyourselfandenjoy/recursion/knights_tour"
	nqueens "justbeyourselfandenjoy/recursion/n_queens"
	towerofhanoi "justbeyourselfandenjoy/recursion/tower_of_hanoi"
)

func main() {
	factorial.RunFactorial()
	fibonacci.FibonacciRun()
	fibonacci.DynamicFibonacciRun()
	knightstour.KnightsTourRun()
	nqueens.NQueensRun()
	towerofhanoi.TowerOfHanoiRun()
}
