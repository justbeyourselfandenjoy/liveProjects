package fastexp

import (
	"fmt"
	"math"
	"strconv"
)

func fastExp(num, pow int) int {
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
			result *= num
			pow--
		}
		num *= num
		pow /= 2
	}
	return num * result
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

func FastExtRun() {
	fmt.Println("Running FastExtRun()")

	for {
		input := ""

		fmt.Printf("num: ")
		fmt.Scanln(&input)
		if len(input) == 0 {
			break
		}
		num, err := strconv.Atoi(input)
		if err != nil || num < 1 {
			break
		}

		fmt.Printf("pow: ")
		fmt.Scanln(&input)
		if len(input) == 0 {
			break
		}
		pow, err := strconv.Atoi(input)
		if err != nil || pow < 1 {
			break
		}

		fmt.Printf("mod: ")
		fmt.Scanln(&input)
		if len(input) == 0 {
			break
		}
		mod, err := strconv.Atoi(input)
		if err != nil || mod < 1 {
			break
		}

		fmt.Println(" num  pow  mod   (num ^ pow) (num ^ pow % mod)")
		if fastExp(num, pow) != int(math.Pow(float64(num), float64(pow))) {
			fmt.Printf("fastExp(num, pow) != math.Pow(num, pow) for num=%v and pow = %v\n", num, pow)
			break
		}

		if fastExpMod(num, pow, mod) != int(math.Pow(float64(num), float64(pow)))%mod {
			fmt.Printf("fastExp(num, pow, mod) != math.Pow(num, pow)%%mod for num=%v, pow = %v, and mod = %v\n", num, pow, mod)
			break
		}

		fmt.Printf("%4v %4v %4v %13v %17v\n", num, pow, mod, fastExp(num, pow), fastExpMod(num, pow, mod))
	}
}
