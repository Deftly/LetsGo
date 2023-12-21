package main

import "fmt"

/*
* for is Go's only looping construct
 */

func main() {
	i := 1
	// Go's version of while
	for i <= 3 {
		fmt.Println(i)
		i = i + 1
	}
	// C style standard loop
	for j := 7; j <= 9; j++ {
		fmt.Println(j)
	}
	// Infinite loop
	for {
		fmt.Println("loop")
		break
	}
	// Loop with continue
	for n := 0; n <= 5; n++ {
		if n%2 == 0 {
			continue
		}
		fmt.Println(n)
	}
}
