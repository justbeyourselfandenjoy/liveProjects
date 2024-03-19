package rsa

import (
	"fmt"
	"math/rand"
	"time"
)

// Initialize a pseudorandom number generator.
var random = rand.New(rand.NewSource(time.Now().UnixNano())) // Initialize with a changing seed

func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	if b == 0 {
		return a
	}
	if a == 0 {
		return b
	}
	r := a % b
	for r > 0 {
		b, r = r, b%r
	}
	return b

}

func lcm(a, b int) int {
	return a * b / gcd(a, b)
}

// Calculate the totient function λ(n)
// where n = p * q and p and q are prime.
func totient(p, q int) int {
	return lcm(p-1, q-1)
}

// Pick a random exponent e in the range [1, λn]
// such that gcd(e, λn) = 1.
func randomExponent(λn int) int {
	for {
		e := randRange(3, λn)
		if gcd(e, λn) == 1 {
			return e
		}
	}
}

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

// Calculate the inverse of a in the modulus.
// See https://en.wikipedia.org/wiki/Extended_Euclidean_algorithm#Modular_integers
// Look at:
//
//	Section "Computing multiplicative inverses in modular structures"
//	Subsection "Modular integers"
func inverseMod(a, modulus int) int {
	t := 0
	newt := 1
	r := modulus
	newr := a

	for newr != 0 {
		quotient := r / newr
		t, newt = newt, t-quotient*newt
		r, newr = newr, r-quotient*newr
	}

	if r > 1 {
		return -1
	}

	if t < 0 {
		t = t + modulus
	}
	return t
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

func RSARun() {
	fmt.Println("Running RSARun()")
	// Pick two random primes p and q.
	const numTests = 100
	p := findPrime(10000, 50000, numTests)
	q := findPrime(10000, 50000, numTests)

	// Calculate the public key modulus n.
	n := p * q

	// Calculate Carmichael's totient function λ(n).
	λn := totient(p, q)

	// Pick a random public key exponent e in the range [3, λn)
	// where gcd(e, λn) = 1.
	e := randomExponent(λn)

	// Find the inverse of e mod λn.
	d := inverseMod(e, λn)

	// Print out the important values.
	fmt.Printf("*** Public ***\n")
	fmt.Printf("Public key modulus:    %d\n", n)
	fmt.Printf("Public key exponent e: %d\n", e)
	fmt.Printf("\n*** Private ***\n")
	fmt.Printf("Primes:    %d, %d\n", p, q)
	fmt.Printf("λ(n):      %d\n", λn)
	fmt.Printf("d:         %d\n", d)

	for {
		var m int
		fmt.Printf("\nMessage:    ")
		fmt.Scan(&m)
		if m < 1 {
			break
		}

		ciphertext := fastExpMod(m, e, n)
		fmt.Printf("Ciphertext: %d\n", ciphertext)

		plaintext := fastExpMod(ciphertext, d, n)
		fmt.Printf("Plaintext:  %d\n", plaintext)
	}
}
