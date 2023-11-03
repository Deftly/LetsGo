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


## Functions Are Values

### Function Type Declarations

### Anonymous Functions

## Closures

### Passing Functions as Parameters

### Returning Functions from Functions

## defer

## Go Is Call By Value

## Exercises

## Wrapping Up
