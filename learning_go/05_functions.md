# Functions

<!--toc:start-->
- [Functions](#functions)
  - [Declaring and Calling Functions](#declaring-and-calling-functions)
    - [Simulating Named and Optional Parameters](#simulating-named-and-optional-parameters)
    - [Variadic Input Parameters and Slices](#variadic-input-parameters-and-slices)
    - [Multiple Return Values](#multiple-return-values)
    - [Ignoring Returned Values](#ignoring-returned-values)
    - [Named Return Values](#named-return-values)
    - [Blank Returns-Never Use These!](#blank-returns-never-use-these)
  - [Functions Are Values](#functions-are-values)
    - [Function Type Declarations](#function-type-declarations)
    - [Anonymous Functions](#anonymous-functions)
  - [Closures](#closures)
    - [Passing Functions as Parameters](#passing-functions-as-parameters)
    - [Returning Functions from Functions](#returning-functions-from-functions)
  - [defer](#defer)
  - [Go Is Call By Value](#go-is-call-by-value)
  - [Exercises](#exercises)
  - [Wrapping Up](#wrapping-up)
<!--toc:end-->

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
Creating a closure that references local variables and passing that closure to another function has several implications and is a very useful pattern that appears several times in the standard library. 

Now for some examples. The `sort` package in the standard library has a function called `sort.Slice`. It takes in any slice and a function that is used to sort the slice that's passed in.

> **_NOTE:_** The `sort.Slice` function predates the addition of Generics in Go. Because of this it has to do some internal magic to work with any kind of slice. We'll look at how it does this in [section 16](./16_here_be_dragons_reflect_unsafe_and_cgo.md)

We'll use closures to sort the same data in different ways:
```go
type Person struct {
  FirstName string
  LastName  string
  Age       int
}

people := []Person{
  {"Pat", "Patterson", 37},
  {"Tracy", "Carter", 23},
  {"Fred", "Fredson", 18},
}
fmt.Println(people)

// Sort slice by last name
sort.Slice(people, func(i, j int) bool {
  return people[i].LastName < people[j].LastName
})
fmt.Println(people)

// Sort slice by age
sort.Slice(people, func(i, j int) bool {
  return people[i].Age < people[j].Age
})
fmt.Println(people)
```
The closure passed to `sort.Slice` has two parameters, `i` and `j` but within the closure `people` is used to sort using its different fields. Running this code gives us the following output:
```
[{Pat Patterson 37} {Tracy Carter 23} {Fred Fredson 18}]
[{Tracy Carter 23} {Fred Fredson 18} {Pat Patterson 37}]
[{Fred Fredson 18} {Tracy Carter 23} {Pat Patterson 37}]
```
The `people` slice is changed by the call to `sort.Slice`, we'll cover this briefly later in [this section](#go-is-call-by-value) and in more detail in the [next section](./06_pointers.md).

### Returning Functions from Functions
We saw that we can use a closure to pass some function state to another function, we can also return a closure from a function:
```go
func makeMult(base int) func(int) int {
	return func(factor int) int {
		return base * factor
	}
}

func main() {
	twoBase := makeMult(2)
	threeBase := makeMult(3)
	for i := 0; i < 3; i++ {
		fmt.Println(twoBase(i), threeBase(i))
	}
}
```
This gives us the following output:
```
0 0
2 3
4 6
```
> **_NOTE:** You may sometimes hear the term *higher-order functions*, especially when talking about functional programming languages. That's a fancy way of saying that a function has a function as an input parameter or a return value.

## defer
Programs often create temporary resources, like files or network connections, that need to be cleaned up. This cleanup has to happen no matter how many exit points a function has, or whether a function completed successfully or not. In Go, cleanup code is attached to the function with `defer`.

Let's see how `defer` works by writing a simple version of `cat`, the Unix utility for printing the contents of a file:
```go
func main() {
	if len(os.Args) < 2 {
		log.Fatal("no file specified")
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	data := make([]byte, 2048)
	for {
		count, err := f.Read(data)
		os.Stdout.Write(data[:count])
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
	}
}
```
Once we know there is a valid file handle we need to close it after we just it. To ensure the cleanup code runs, we use `defer`, followed by a function or method call, in this case we use the `Close` method on the file variable. Normally a function call runs immediately, but `defer` delays the invocation until the surrounding function exits.

You can `defer` multiple functions in a Go function, they will run in last-in-first-out order. The code within `defer` function run *after* the return statement. You can supply a function with input parameters to a `defer`. The input parameters are evaluated immediately and their values are stored until the function runs. Here is an example:
```go
func deferExample() int {
	a := 10
	defer func(val int) {
		fmt.Println("first:", val)
	}(a)

	a = 20
	defer func(val int) {
		fmt.Println("second:", val)
	}(a)

	a = 30
	fmt.Println("exiting:", a)
	return a
}
```
The output of this code is:
```
exiting: 30
second: 20
first: 10
```
Named return values allow a deferred function examine and modify the return values of its surrounding function. This allows your code to take actions based on an error. We'll talk more about a pattern that uses `defer` to add contextual information to an error returned from a function in [section 9](./09_errors.md). The next example shows a way to use named return values and `defer` to handle database transaction cleanup:
```go
func DoSomeInserts(ctx context.Context, db *sql.DB, value1, value2 string) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			err = tx.Commit()
		}
		if err != nil {
			tx.Rollback()
		}
	}()
	_, err = tx.ExecContext(ctx, "INSERT INTO FOO (val) values $1", value1)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO FOO (val) values $2", value2)
	if err != nil {
		return err
	}
	return nil
}
```
In this example function we create a transaction to do a series of database inserts. If any of them fail we want to roll back(not modify the database). If all of the transactions succeed we want to commit(store the database changes). We use a closure with `defer` to check if `err` has a value. If it doesn't, we run a `tx.Commit()`, which could also generate an error. If it does, the value `err` is modified. If any database interaction returns an error, we call `tx.Rollback()`.

A common pattern in Go is for a function that allocates a resource to also return a closure that cleans up the resource. Here are some modifications to our simple `cat` program that uses this pattern:
```go
// Helper function that opens a file and returns a closure
func getFile(name string) (*os.File, func(), error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, nil, err
	}
	return file, func() {
		file.Close()
	}, nil
}
```
Now in `main` we can use `getFile` like so:
```go
	f, closer, err := getFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
  defer closer()
```
Because Go doesn't allow unused variables, returning the `closer` from the function means the program will not compile if the function is not called.

## Go Is Call By Value
Go is a *call by value* language which means that when you supply a variable for a parameter to a function, Go *always* makes a copy of the value of the variable. Let's explore this with some examples:
```go
type person struct {
	name string
	age  int
}

func modifyFails(i int, s string, p person) {
	i = i * 2
	s = "Goodbye"
	p.name = "Bob"
}

func main() {
	p := person{}
	i := 2
	s := "Hello"
	modifyFails(i, s, p)
	fmt.Println(i, s, p) // 2 Hello { 0}
}
```
If we run this code we can see that the function won't change the values of the parameters passed into it. This doesn't just apply to primitive types, as we can see we passed in a struct and it showed the same behavior.

The behavior is different for maps and slices:
```go
func modMap(m map[int]string) {
	m[2] = "hello"
	m[3] = "goodbye"
	delete(m, 1)
}

func modSlice(s []int) {
	for i, v := range s {
		s[i] = v * 2
	}
	s = append(s, 10)
}

func main() {
	m := map[int]string{
		1: "first",
		2: "second",
	}
	modMap(m)
	fmt.Println(m) // map[2:hello 3:goodbye]

	s := []int{1, 2, 3}
	modSlice(s)
	fmt.Println(s) // [2 4 6]
}
```
Any changes made to a map parameter are reflected in the variable passed into the function. For a slice we can modify any element in the slice, but you can't lengthen the slice. This true for maps and slices that are passed directly into functions as well as map and slice fields in structs.

The reason why maps and slices behave differently is because maps and slices are both implemented with pointers, which we will be going into detail on in the [next section](./06_pointers.md).

> **_NOTE:_** Every type in Go is a value type. It's just that sometimes the value is a pointer.

Call by value is one reason why Go's limited support for constants is only a minor handicap. Since variables are passed by value you can be sure that calling a function doesn't modify the variable whose value was passed in(unless the variable is of a pointer type). It makes it easier to understand the flow of data through your program when functions don't modify their input parameters and instead return new values. 

However, there are cases where you need to pass something mutable to a function, that's when you need a pointer.

## Exercises
1. The simple calculator program doesn't handle one error case: division by zero. Change the function signature for the math operations to return both an `int` and an `error`. In the `div` function, if the divisor is `0`, return `errors.New("division by zero")` for the error. In all other cases, return `nil`. Adjust the `main` function to check for this error.
2. Write a function called `fileLen` that has an input parameter of type `string` and return an `int` and an `error`. The function takes in a file name and returns the number of bytes in the file. If there is an error reading the file, return the error. Use `defer` to make sure the file is closed properly.
3. Write a function called `prefixer` that has an input parameter of type `string` and returns a function that has an input parameter of type `string` and returns a `string`. The returned function should prefix its input with the string passed into `prefixer`. Use the following `main` function to test `prefixer`:
```go
func main() {
  helloPrefix := prefixer("Hello")
  fmt.Println(helloPrefix("Bob")) // should print Hello Bob
  fmt.Println(helloPrefix("Maria")) // should print Hello Maria
}
```
> [Solutions](./exercises/05/)

## Wrapping Up
This section covered functions in Go and their unique features. In the [next section](./06_pointers.md) we'll cover pointers and how to take advantage of them to write efficient programs.
