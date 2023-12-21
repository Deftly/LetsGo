package main

import "fmt"

/*
* Array types are a numbered sequenc of elements of a specific length.
* This means [3]int is a different type than [4]int. This makes arrays rigid
* and far less common in use than slices.
 */

func main() {
	var a [5]int
	fmt.Println("emp:", a)

	a[4] = 100
	fmt.Println("set:", a)
	fmt.Println("get:", a[4])

	fmt.Println(len(a))

	b := [5]int{1, 2, 3, 4, 5}
	fmt.Println("dcl:", b)

	// Array types are one-dimensional, but you can compose types to mimic
	// multi-dimensional structures.
	var twoD [2][3]int
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			twoD[i][j] = i + j
		}
	}

	// Arrays appear in the form [v1 v2 v3 ...] when printed with fmt.Println.
	fmt.Println("2d: ", twoD)
}
