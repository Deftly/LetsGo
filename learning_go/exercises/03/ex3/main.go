package main

import "fmt"

type Employee struct {
	firstName string
	lastName  string
	id        int
}

func main() {
	sidd := Employee{
		"Siddharth",
		"Buddharaju",
		24,
	}
	chris := Employee{
		firstName: "Chris",
		id:        11,
		lastName:  "Hyde",
	}
	josh := Employee{}
	josh.firstName = "Josh"
	josh.lastName = "Ferguson"
	josh.id = 30
	fmt.Println(sidd)
	fmt.Println(chris)
	fmt.Println(josh)
}
