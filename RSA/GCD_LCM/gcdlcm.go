package gcdlcm

import (
	"fmt"
	"strconv"
)

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

func GCDLCMRun() {
	fmt.Println("Running GCDLCMRun()")

	for {
		sa, sb := "", ""
		fmt.Printf("A: ")
		fmt.Scanln(&sa)
		if len(sa) == 0 {
			break
		}
		a, err := strconv.Atoi(sa)
		if err != nil {
			break
		}

		fmt.Printf("B: ")
		fmt.Scanln(&sb)
		if len(sb) == 0 {
			break
		}
		b, err := strconv.Atoi(sb)
		if err != nil {
			break
		}

		fmt.Printf("A: %d, B: %d, gcd(%d, %d): %d, lcm(%d, %d): %d\n", a, b, a, b, gcd(a, b), a, b, lcm(a, b))
	}
}
