package main

import (
	"fmt"
	"time"
)

func main() {
	// Basic switch
	i := 2
	fmt.Print("Write ", i, " as ")
	switch i {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	case 3:
		fmt.Println("three")
	}

	// You can separate multiple expression for the same case statement using commas.
	// There is also an optional default case
	switch time.Now().Weekday() {
	case time.Saturday, time.Sunday:
		fmt.Println("It's the weekend")
	default:
		fmt.Println("It's a weekday")
	}

	// Blank switch statements are another way to express if/else logic and
	// case expressions can be non-constants
	t := time.Now()
	switch {
	case t.Hour() < 12:
		fmt.Println("It's before noon")
	default:
		fmt.Println("It's after noon")
	}

	// A type switch compares types instead of values. Can be used to find out
	// the type of an interface value.
	whatType := func(i interface{}) {
		switch t := i.(type) {
		case bool:
			fmt.Println("The type is bool")
		case int:
			fmt.Println("The type is int")
		default:
			fmt.Printf("This type is unknown %T\n", t)
		}
	}

	whatType(true)
	whatType(1)
	whatType("test")
}
