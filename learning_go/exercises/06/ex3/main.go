package main

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

func main() {
	var people []Person
	// Pre-allocating space for the slice makes the program much faster append
	// reduces the number of gc cycles
	// people := make([]Person, 0, 10_000_000)
	for i := 0; i < 10_000_000; i++ {
		people = append(people, MakePerson("Fred", "Williamson", 25))
	}
}
