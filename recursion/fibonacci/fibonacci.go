package fibonacci

import (
	"fmt"
	"strconv"
)

func fibonacci(n int64) int64 {
	if n <= 1 {
		return n
	}
	return fibonacci(n-1) + fibonacci(n-2)
}

func FibonacciRun() {
	for {
		// Get n as a string.
		var nString string
		fmt.Printf("N: ")
		fmt.Scanln(&nString)

		// If the n string is blank, break out of the loop.
		if len(nString) == 0 {
			break
		}

		// Convert to int and calculate the Fibonacci number.
		n, _ := strconv.ParseInt(nString, 10, 64)
		fmt.Printf("fibonacci(%d) = %d\n", n, fibonacci(n))
	}
}
