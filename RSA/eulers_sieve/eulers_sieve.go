package eulerssieve

import (
	"fmt"
	"time"
)

// Build a sieve of Eratosthenes.
func eulersSieve(max int) []bool {
	if max < 2 {
		return []bool{false}
	}
	result := make([]bool, max+1)
	result[2] = true
	for i := 3; i <= max; i += 2 {
		result[i] = true
	}
	for p := 3; p <= max; p += 2 {
		if result[p] {
			maxQ := max / p
			if maxQ%2 == 0 {
				maxQ--
			}
			for q := maxQ; q >= p; q -= 2 {
				if result[q] {
					result[p*q] = false
				}
			}
		}
	}
	return result
}

// Print out the primes in the sieve.
func printSieve(sieve []bool) {
	for i := range sieve {
		if i > 1 && sieve[i] {
			fmt.Printf("%v ", i)
		}
	}
	fmt.Println()
}

// Convert the sieve into a slice holding prime numbers.
func sieveToPrimes(sieve []bool) []int {
	result := make([]int, 0)
	for i := range sieve {
		if i > 1 && sieve[i] {
			result = append(result, i)
		}
	}
	return result
}

func EulersSieveRun() {
	fmt.Println("Running EulersSieveRun()")

	var max int
	fmt.Printf("Max: ")
	fmt.Scan(&max)

	start := time.Now()
	sieve := eulersSieve(max)

	elapsed := time.Since(start)
	fmt.Printf("Elapsed: %f seconds\n", elapsed.Seconds())

	if max <= 1000 {
		printSieve(sieve)

		primes := sieveToPrimes(sieve)
		fmt.Println(primes)
	}
}
