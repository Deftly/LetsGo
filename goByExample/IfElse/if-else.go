package main

import "fmt"

func main() {
	if 7%2 == 0 {
		fmt.Println("7 is even")
	} else {
		fmt.Println("7 is odd")
	}

	// You can have an if statement without an else
	if 8%4 == 0 {
		fmt.Println("8 is divisible by 4")
	}

	if 7%2 == 0 || 8%2 == 0 {
		fmt.Println("either 8 or 7 is even")
	}

	// A statemetn can precede conditionals; any variables decalred in this
	// statemetn are available in the current and all subsequent branches.
	if num := 9; num < 0 {
		fmt.Println(num, "is negative")
	} else if num < 10 {
		fmt.Println(num, "has 1 digit")
	} else {
		fmt.Println(num, "has multiple digits")
	}
}
