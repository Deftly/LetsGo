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

Without generics the only way to avoid duplicated code would be to modify our tree implementation so that it uses an interface to specify how to order values. The interface could like something like this:
```go
type Orderable interface {
	// Order returns:
	// a value < 0 when the Orderable is less than the supplied value,
	// a value > 0 when the Orderable is greater than the supplied value,
	// and 0 when the two values are equal
	Order(any) int
}
```
Using the `Orderable` interface, we can modify the `Tree` to support it:
```go
type Tree struct {
  val         Orderable
  left, right *Tree
}

func (t *Tree) Insert(val Orderable) *Tree {
  if t == nil {
    return &Tree{val: val}
  }

	switch comp := val.Order(t.val); {
	case comp < 0:
		t.left = t.left.Insert(val)
	case comp > 0:
		t.right = t.right.Insert(val)
	}
	return t
}
```
With an `OrderableInt` type, we can insert `int` values:
```go
type OrderableInt int

func (oi OrderableInt) Order(val any) int {
  return int(oi - val.(OrderableInt))
}

func main() {
  var it *Tree
  it = it.Insert(OrderableInt(5))
  it = it.Insert(OrderableInt(3))
  // etc...
}
```
While this code works, it doesn't allow the compiler to validate that the values inserted into our data structure are all the same. If we also had the following type:
```go
type OrderableString string

func (os OrderableString) Order(val any) int {
  return strings.Compare(string(os), val.(string))
}
```
The following code compiles:
```go
var it *Tree
it = it.Insert(OrderableInt(5))
it = it.Insert(OrderableString("nope"))
```
While the compiler accepts this code, the program will panic when attempting to insert an `OrderableString` into a `Tree` that already contains an `OrderableInt`:
```shell
panic: interface conversion: interface {} is main.OrderableInt, not string
```
The full implementation of the tree can be found [here](./examples/nonGenericTree/main.go)

## Introducing Generics in Go
We'll take a first look at generics by implementing a stack:
```go
type Stack[T any] struct {
	vals []T
}

func (s *Stack[T]) Push(val T) {
	s.vals = append(s.vals, val)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.vals) == 0 {
		var zero T
		return zero, false
	}
	top := s.vals[len(s.vals)-1]
	s.vals = s.vals[:len(s.vals)-1]
	return top, true
}
```


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
