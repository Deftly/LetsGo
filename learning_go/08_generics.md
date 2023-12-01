# Generics

## Generics Reduce Repetitive Code and Increase Type Safety
Go is a statically typed language, meaning that the types of variables and parameters are checked when the code is compiled. Built-in types(maps, slice, channels) and functions(such as `len`, `cap`, or `make`) are able to accept and return values of different concrete types, but until Go 1.18, user-defined Go types and functions could not.

So far we've seen functions that take in parameters whose values are specified when the function is called. In the following example the function has two `int` parameters and returns two `int` values:
```go
func divAndRemainder(num, denom int) (int, int, error) {
	if denom == 0 {
		return 0, 0, errors.New("cannot divide by zero")
	}
	return num / denom, num % denom, nil
}
```
Similarly, when creating structs the type for the fields are specified when the struct is declared:
```go
type Node struct {
  val int
  next *Node
}
```
There are situations where it's useful to write functions or structs where the specific *type* of a parameter or field is left unspecified until it is used.

In an [earlier section](./07_types_methods_and_interfaces.md#code-your-methods-for-nil-instances) we looked at a binary tree for ints. If you want a binary tree for strings or float64s and you wanted type saftey you have a few options. The first is writing a custom tree for each type, but that much duplicated code is verbose and error-prone.

Without generics the only way to avoid duplicated code would be to modify our tree implementation so that it uses an interface to specify how to order values.

## Introducing Generics in Go

## Generic Functions Abstract Algorithms

## Generics and Interfaces

## Use Type Terms to Specify Operators

## Type Inference and Generics

## Type Elements Limit Constants

## Combining Generic Functions with Generic Data Structures

## More on comparable

## Things That Are Left Out

## Idiomatic Go and Generics

## Adding Generics to The Standard Library

## Future Features Unlocked

## Exercises

## Wrapping Up
