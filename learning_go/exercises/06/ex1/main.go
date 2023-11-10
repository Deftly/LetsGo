package main

import "fmt"

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

func MakePerson(firstName, lastName string, age int) Person {
	return Person{
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
	}
}

func MakePersonPointer(firstName, lastName string, age int) *Person {
	return &Person{
		FirstName: firstName,
		LastName:  lastName,
		Age:       age,
	}
}

// compile with -gcflags="-m"
func main() {
	p := MakePerson("Nitin", "Kunaparaju", 16)
	fmt.Println(p)
	p2 := MakePersonPointer("Anish", "Kunaparaju", 13)
	fmt.Println(p2)
}
