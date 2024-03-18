package primality

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const numTests = 20

// Initialize a pseudorandom number generator.
var random = rand.New(rand.NewSource(time.Now().UnixNano())) // Initialize with a changing seed

// Return a pseudo random number in the range [min, max).
func randRange(min int, max int) int {
	return min + random.Intn(max-min)
}

func fastExpMod(num, pow, mod int) int {
	if pow < 0 {
		num = 1 / num
		pow = -pow
	}
	if pow == 0 {
		return 1
	}
	result := 1
	for pow > 1 {
		if pow%2 == 1 {
			result = (result * num) % mod
			pow--
		}
		num = (num * num) % mod
		pow /= 2
	}
	return (num * result) % mod

}

// Perform tests to see if a number is (probably) prime.
func isProbablyPrime(p int, numTests int) bool {
	for i := 0; i < numTests; i++ {
		result := fastExpMod(randRange(1, p), p-1, p)
		if result < 0 {
			fmt.Println("result is too high, skipping")
		}
		if result != 1 {
			return false
		}

	}
	return true
}

// Probabilistically find a prime number within the range [min, max).
func findPrime(min, max, numTests int) int {
	for {
		p := randRange(min, max)
		if p%2 != 0 && isProbablyPrime(p, numTests) {
			return p
		}
	}
}

func testKnownValues() {
	primes := []int{
		10009, 11113, 11699, 12809, 14149,
		15643, 17107, 17881, 19301, 19793,
	}
	composites := []int{
		10323, 11397, 12212, 13503, 14599,
		16113, 17547, 17549, 18893, 19999,
	}

	fmt.Printf("Probability: %f%%\n\n", (1.0-1.0/math.Pow(2, numTests))*100.0)

	fmt.Println("Primes:")
	for _, number := range primes {
		if isProbablyPrime(number, 10) {
			fmt.Println(number, " Prime")
		} else {
			fmt.Println(number, " Composite")
		}
	}
	fmt.Println()

	fmt.Println("Composites:")
	for _, number := range composites {
		if isProbablyPrime(number, 10) {
			fmt.Println(number, " Prime")
		} else {
			fmt.Println(number, " Composite")
		}
	}
}

func PrimalityTestRun() {
	fmt.Println("Running PrimaryTestRun()")

	fmt.Println(time.Now())

	// Test some known primes and composites.
	testKnownValues()

	// Generate random primes.
	for {
		// Get the number of digits.
		var numDigits int
		fmt.Printf("\n# Digits: ")
		fmt.Scan(&numDigits)
		if numDigits < 1 {
			break
		}

		// Calculate minimum and maximum values.
		min := int(math.Pow(10.0, float64(numDigits-1)))
		max := 10 * min
		if min == 1 {
			min = 2
		} // 1 is not prime.

		// Find a prime.
		fmt.Printf("Prime: %d\n", findPrime(min, max, numTests))
	}
}
