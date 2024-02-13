package factorial

import (
	"fmt"
	"math/big"
)

func factorial(n int64) int64 {
	if n <= 1 {
		return 1
	}
	return n * factorial(n-1)
}

func iterativeFactorial(n int64) int64 {
	result := int64(1)
	for i := int64(2); i <= int64(n); i++ {
		result *= i
	}

	return result
}

func bigFactorial(n *big.Int) *big.Int {
	one := big.NewInt(1)
	if n.Cmp(one) <= 0 {
		return big.NewInt(1)
	}

	var nMinus1 big.Int
	nMinus1.Sub(n, one)

	var result big.Int
	result.Mul(n, bigFactorial(&nMinus1))
	return &result
}

func RunFactorial() {
	fmt.Println("Running RunFactorial()")

	var n int64
	for n = 0; n <= 21; n++ {
		fmt.Printf("%3d! = %20d\n", n, factorial(n))
	}
	for n = 0; n <= 21; n++ {
		fmt.Printf("%3d! = %20d\n", n, bigFactorial(big.NewInt(n)))
	}
	fmt.Println()
}
