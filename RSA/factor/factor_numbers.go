package factor

import (
	"fmt"
	"time"
)

var primes []int

func findFactors(num int) []int {
	result := make([]int, 0)

	if num < 0 {
		result = append(result, -1)
		num = -num
	}

	for num%2 == 0 {
		result = append(result, 2)
		num /= 2
	}
	factor := 3
	for factor*factor <= num {
		if num%factor == 0 {
			result = append(result, factor)
			num /= factor
		} else {
			factor += 2
		}
	}
	if num > 1 {
		result = append(result, num)
	}
	return result
}

func findFactorsSieve(num int) []int {
	result := make([]int, 0)

	if num < 0 {
		result = append(result, -1)
		num = -num
	}

	for num%2 == 0 {
		result = append(result, 2)
		num /= 2
	}

	for _, factor := range primes {
		for num%factor == 0 {
			result = append(result, factor)
			num /= factor
			if num == 1 {
				break
			}
		}
		if factor*factor > num {
			result = append(result, num)
			break
		}
	}
	return result
}

func initPrimes(max int) {
	fmt.Printf("Initializing the primes up to %v ... ", max)
	start := time.Now()
	result := eulersSieve(max)
	primes = sieveToPrimes(result)
	elapsed := time.Since(start)
	fmt.Printf("[%f seconds]\n", elapsed.Seconds())
}

func multiplySlice(slice []int) int {
	if len(slice) == 0 {
		return 0
	}
	result := 1
	for i := range slice {
		result *= slice[i]
	}
	return result
}

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

func FactorNumbersRun() {
	fmt.Println("Running FactorNumbersRun()")

	var num int
	fmt.Printf("Max: ")
	fmt.Scan(&num)

	// Find the factors the slow way.
	start := time.Now()
	factors := findFactors(num)
	elapsed := time.Since(start)
	fmt.Printf("findFactors:       %f seconds\n", elapsed.Seconds())
	fmt.Println(multiplySlice(factors))
	fmt.Println(factors)
	fmt.Println()

	initPrimes(1600000000)
	// Use the Euler's sieve to find the factors.
	start = time.Now()
	factors = findFactorsSieve(num)
	elapsed = time.Since(start)
	fmt.Printf("findFactorsSieve: %f seconds\n", elapsed.Seconds())
	fmt.Println(multiplySlice(factors))
	fmt.Println(factors)
	fmt.Println()
}
