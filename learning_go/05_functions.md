# Functions

## Declaring and Calling Functions
We've already seen functions being declared and used. Every Go program starts from a `main` function and we've been calling `fmt.Println`. Now lets take a look at a function that takes in parameters:
```go
func div(num int, denom int) int {
  if denom == 0 {
    return 0
  }
  return num / denom
}
```
> **_NOTE:_** When you have two or more consecutive input parameters of the same type you can specify the type once for all of them like this: `func div(num, denom int) int {`

### Simulating Named and Optional Parameters
Go *doesn't* have named and optional input parameters. With one exception you must supply all of the parameters for a function. To emulate named and optional parameters, define a struct that has fields that match the desired parameters and pass the struct to your function:
```go
type MyFuncOpts struct {
	FirstName string
	LastName  string
	Age       int
}

func MyFunc(opts MyFuncOpts) error {
	// do something here
}

func main() {
	MyFunc(MyFuncOpts{
		LastName: "Patel",
		Age:      50,
	})
	MyFunc(MyFuncOpts{
		FirstName: "Joe",
		LastName:  "Smith",
	})
}
```
### Variadic Input Parameters and Slices
You may have noticed that `fmt.Println` allows any number of input parameters, it can do this because of Go's support for *variadic parameters*. The variadic parameter *must* be the last parameter in the input list and is indicated with `...` *before* the type. The variable that is created within the function is a slice of the specified type:
```go
func addTo(base int, vals ...int) []int {
	out := make([]int, 0, len(vals))
	for _, v := range vals {
		out = append(out, base+v)
	}
	return out
}

func main() {
	fmt.Println(addTo(3)) // []
	fmt.Println(addTo(3, 2)) // [5]
	fmt.Println(addTo(3, 2, 4, 6, 8)) // [5 7 9 11]
	a := []int{4, 3}
	fmt.Println(addTo(3, a...)) // [7 6]
	fmt.Println(addTo(3, []int{1, 2, 3, 4, 5}...)) // [4 5 6 7 8]
}
```
You can supply a slice as input for the variadic parameter by adding `...` *after* the variable or slice literal.

### Multiple Return Values
Lets add a small feature to our previous division function that uses Go's ability to return multiple return values:
```go
func divAndRemainder(num, denom int) (int, int, error) {
	if denom == 0 {
		return 0, 0, errors.New("cannot divide by zero")
	}
	return num / denom, num % denom, nil
}
```
Another new concept shown here is create and returning an `error`. If you want to learn more about errors check out [section 9](./09_errors.md). For now, just know that you use Go's multiple return value support to return an `error` if something goes wrong in a function. By convention, the `error` is always the last value return from a function.

This is how we would call the updated function:
```go
func main() {
	result, remainder, err := divAndRemainder(5, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(result, remainder)
}
```
### Ignoring Returned Values
If a function returns multiple values but you don't need to read one or more of the values, assign the unused values to the name `_`. For example, if we weren't going to read `remainder`, we would write the assignment as `result, _, err := divAndRemainder(5, 2)`.

You can implicitly ignore *all* of the return values for a function. We've actually been doing this with `fmt.Println`, which returns two values but it's idiomatic to ignore them. In most cases you should be explicit that you are ignoring return values using underscores. 

### Named Return Values
Go allows you to specify names for your return values:
```go
func divAndRemainder(num, denom int) (result int, remainder int, err error) {
	if denom == 0 {
		err = errors.New("cannot divide by zero")
		return result, remainder, err
	}
	result, remainder = num/denom, num%denom
	return result, remainder, err
}
```
When you give names to your return values you are pre-declaring variables that you use within the function to hold the return values. Named return values are initialized to their zero values when created, meaning you can return them without any explicit use or assignment.

While named return values can help clarify your code they can also lead to some problems:
- The first is shadowing, just like any other variable you can shadow a named return value.
- The second is that you don't have to return them:
```go
func divAndRemainder(num, denom int) (result int, remainder int, err error) {
	// assign some values
	result, remainder = 20, 30
	if denom == 0 {
		return 0, 0, errors.New("cannot divide by zero") // these will be the actual values returned
	}
	return num / denom, num % denom, nil // these will be the actual values returned
}
```
The Go compiler inserts code that assigns whatever is returned to the return parameters. The named return parameters give a way to declare an *intent* to use variables to hold return values but don't *require* you to use them.

There is one situation where named return values are essential, we will cover that when talking about `defer` later in this section.

### Blank Returns-Never Use These!
If you have named return values you can just write `return` without specifying values that are returned. This returns the last values assigned to the named return values:
```go
func divAndRemainder(num, denom int) (result, remainder int, err error) {
	if denom == 0 {
		err = errors.New("cannot divide by zero")
		return
	}
	result, remainder = num/denom, num%denom
	return
}
```
Most experienced Go developers consider blank returns a bad idea because they make it harder to understand data flow. Good software is clear and readable. When you use a blank return the reader of your code needs to scan back through the program to find the last value assigned to the return parameters to see what is actually being returned.

## Functions Are Values
The type of a function is built out of the keyword `func` and the types of the parameters and return values. This combination is called the *signature* of the function. Any function that has the exact same number and types of parameters and return values meets the type signature.

Since functions are values you can declare a function variable:
```go
var myFuncVariable func(string) int
```
`myFuncVariable` can be assigned any function that has a single parameter of type `string` and returns a single value of type `int`:
```go
func f1(a string) int {
	return len(a)
}

func f2(a string) int {
	total := 0
	for _, v := range a {
		total += int(v)
	}
	return total
}

func main() {
	var myFuncVariable func(string) int
	myFuncVariable = f1
	result := myFuncVariable("Hello") // 5
	fmt.Println(result)

	myFuncVariable = f2
	result = myFuncVariable("Hello") // 500
	fmt.Println(result)
}
```
The default zero value for a function variable is `nil`. Attempting to run a function variable with a `nil` value results in a panic.

Having functions as values can let us do some clever things. In the next example we'll build a simple calculator using functions as values in a map:
```go
// Create set of functions with the same signature
func add(i, j int) int { return i + j }

func sub(i, j int) int { return i - j }

func mul(i, j int) int { return i * j }

func div(i, j int) int { return i / j }

// Create a map to associate a math operator with each function
var opMap = map[string]func(int, int) int{
	"+": add,
	"-": sub,
	"*": mul,
	"/": div,
}

func main() {
	expressions := [][]string{
		{"2", "+", "3"},
		{"2", "-", "3"},
		{"2", "*", "3"},
		{"2", "/", "3"},
		{"2", "%", "3"},
		{"two", "+", "three"},
		{"5"},
	}
	for _, expression := range expressions {
		if len(expression) != 3 {
			fmt.Println("invalid expression:", expression)
			continue
		}
		p1, err := strconv.Atoi(expression[0])
		if err != nil {
			fmt.Println(err)
			continue
		}
		op := expression[1]
		opFunc, ok := opMap[op]
		if !ok {
			fmt.Println("unsupported operator:", op)
			continue
		}
		p2, err := strconv.Atoi(expression[2])
		if err != nil {
			fmt.Println(err)
			continue
		}
		result := opFunc(p1, p2)
		fmt.Println(result)
	}
}
```
> **_NOTE:_** Don't write fragile programs. The core logic is very short, just 6 lines inside the `for` loop. The rest of the algorithm is error checking and data validation. Skipping these steps produces unstable, unmaintainable code. Error handling is what separates the professionals from the amateurs.

### Function Type Declarations
Just like how you use the `type` keyword to define a `struct` you can use it to define a function type too:
```go
type opFuncType func(int, int) int
```
We can use this to rewrite the `opMap` declaration in our previous example: 
```go
var opMap = map[string]opFuncType{
  // same as before
}
```
One advantage of declaring a function type is that it acts as documentation. It can be helpful to give something a name if you are going to refer to it multiple times. We'll see another reason in a [later section](./07_types_methods_and_interfaces.md#function-types-are-a-bridge-to-interfaces)

### Anonymous Functions
Not only can you assign functions to variables, you can define new functions within a function and assign them to variables:
```go
func main() {
	f := func(j int) {
		fmt.Println("printing", j, "from inside an anonymous function")
	}
	for i := 0; i < 5; i++ {
		f(i)
	}
}
```
You don't have to assign an anonymous function to a variable. You can write them inline and call them immediately:
```go
func main() {
	for i := 0; i < 5; i++ {
		func(j int) {
			fmt.Println("printing", j, "from inside of anonymous function")
		}(i)
	}
}
```
There are two situations where declaring anonymous functions without assigning them to variables is useful:`defer` statements and launching goroutines. We'll talk about `defer` [later in this section](#defer) and goroutines will be covered in [section 12](./12_concurrency_in_go.md).

Unlike a normal function definition, you can assign a new value to a package level anonymous function:
```go
var (
	add = func(i, j int) int { return i + j }
	sub = func(i, j int) int { return i - j }
	mul = func(i, j int) int { return i * j }
	div = func(i, j int) int { return i / j }
)

func main() {
	x := add(2, 3)
	fmt.Println(x) // 5
	changeAdd()
	y := add(2, 3)
	fmt.Println(y) // 8
}

func changeAdd() {
	add = func(i, j int) int { return i + j + j }
}
```
Before using a package-level anonymous function, be very sure you need this capability. You should always try to keep package-level state immutable to make data flow easier to understand. If a function's meaning changes while a program is running, it becomes difficult to understand not just how data flows, but how it is processed.

## Closures
*Closures* is a computer science term that means that functions declared inside of functions are able to access and modify variables declared in the outer function. Here's an example:
```go
func main() {
	a := 20
	f := func() {
		fmt.Println(a) // 20
		a = 30
	}
	f()
	fmt.Println(a) // 30
}
```
What benefit do we get from creating mini-functions within larger functions?

One thing that closures allow is to limit a function's scope. If a function is only going to be called from one other function we can use an inner function to "hide" the called function. This reduces the number of declarations at the package level. 

Closures really become interesting when they are passed to other functions or returned from a function, allowing you to take variables within your function and use those values *outside* of your function.

### Passing Functions as Parameters


### Returning Functions from Functions

## defer

## Go Is Call By Value

## Exercises

## Wrapping Up
