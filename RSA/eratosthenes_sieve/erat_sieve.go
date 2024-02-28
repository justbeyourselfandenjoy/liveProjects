package eratosthenessieve

import (
	"fmt"
	"time"
)

// Build a sieve of Eratosthenes.
func sieveOfEratosthenes(max int) []bool {
	if max < 2 {
		return []bool{false}
	}
	result := make([]bool, max+1)
	for i := range result {
		result[i] = true
	}
	for i := 2; i <= max; i++ {
		if result[i] {
			for j := i + 1; j <= max; j++ {
				if result[j] && j%i == 0 {
					result[j] = false
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

func EratosthenesSieveRun() {
	fmt.Println("Running EratosthenesSieveRun()")

	var max int
	fmt.Printf("Max: ")
	fmt.Scan(&max)

	start := time.Now()
	sieve := sieveOfEratosthenes(max)

	elapsed := time.Since(start)
	fmt.Printf("Elapsed: %f seconds\n", elapsed.Seconds())

	if max <= 1000 {
		printSieve(sieve)

		primes := sieveToPrimes(sieve)
		fmt.Println(primes)
	}
}
