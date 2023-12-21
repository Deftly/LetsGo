package main

import (
	"fmt"
	"math"
)

/*
* Go supports constants of character, string, boolean, and numeric values.
* Constants in go are more limited than in other languages. Think of them
* as a way of giving a name to a literal
 */

// Typed constant
const s string = "constant"

func main() {
	fmt.Println(s)

	// Untyped constant
	const n = 500000000

	// Constant expressions perform arithmetic with arbitrary precision
	const d = 3e20 / n
	fmt.Println(d)

	// A numeric constant has no type until it's given one, such as by an explicity
	// conversion
	fmt.Println(int64(d))

	// A number can be given a type by using it in a context that requires one,
	// such as a variable assignment or a function call
	fmt.Println(math.Sin(n))
}
