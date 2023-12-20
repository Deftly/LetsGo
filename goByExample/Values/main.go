package main

import "fmt"

func main() {
	// Strings can be concatenated with +
	fmt.Println("go" + "lang") // golang

	// Integers and Floats
	fmt.Println("1+1 = ", 1+1)         // 1+1 = 2
	fmt.Println("7.0/3.0 = ", 7.0+3.0) // 7.0/3.0 = 2.3333333333333335

	// Booleans with normal boolean operators
	fmt.Println(true && false) // false
	fmt.Println(true || false) // true
	fmt.Println(!true)         // false
}
