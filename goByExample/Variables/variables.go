package main

import "fmt"

func main() {
	// var a = "initial" var declares 1 or more variables

	// You can declare multiple variables at once
	var b, c int = 1, 2
	fmt.Println(b, c)

	// var d = true Go will infer the type of initialzed variables.

	// Variables declared without a corresponding initialization are zero-valued.
	var e int
	fmt.Println(e)

	// := syntax is shorthand for declaring an initilizing a variable
	// The following is eqivalent to var f string = "apple"
	// This syntax can only be used inside of functions.
	f := "apple"
	fmt.Println(f)
}
