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
The type parameter information(`[T any]`)is placed within brackets and has two parts. The first is the parameter name, it is customary to use capital letters for them. The second part is the *type constraint*, which uses a Go interface to specify which types are valid. If any type is usable, this is specified using the universe block identifier `any`. In the stack declaration we declare `vals` to be of the type `[]T`.

Looking at the method declarations, we refer to the type in the receiver section with `Stack[T]`.

Handling zero values is a little different with generics. In `Pop`, we can't just return `nil` because that's not a valid value for a value type like `int`. The easiest way to get a zero value for a generic is to declare a variable with `var` and return it since `var` always initializes its variable to the zero value if no value is assigned.

Using a generic type is similar to a non-generic one:
```go
func main() {
	var intStack Stack[int]
	intStack.Push(10)
	intStack.Push(20)
	intStack.Push(30)
	v, ok := intStack.Pop()
	fmt.Println(v, ok)
}
```
The only difference is that when we declare our variable, we include the type that we want to use with our stack. If you try to push a string onto this stack, the compiler will catch it and produce a compiler error:
```shell
cannot use "nope" (untyped string constant) as int value in argument to intStack.Push
```
You can view the full implementation of this stack [here](./examples/stack/main.go)

Let's add another method to the stack to tell us if the stack contains a value:
```go
func (s Stack[T]) Contains(val T) bool {
	for _, v := range s.vals {
		if v == val { // invalid operation: v == val (type parameter T is not comparable with ==)
			return true
		}
	}
	return false
}
```
This doesn't work because `any`, just like the `interface{}` doesn't say anything. You can only store and retrieve values of `any` type. To use `==`, you need a different type, and since nearly all Go types can be compared with `==` and `!=`, a built-in interface called `comparable` is defined in the universe block.
```go
type Stack[T comparable] struct {
  vals []T
}
```
Now our `Contains` method works as expected.

## Generic Functions Abstract Algorithms
We saw earlier that not having generics made it hard to write map, reduce, and filter implementations that work for all types. Generics make it easy, here are implementations from the type parameters proposal:
```go
// Map turns a []T1 to a []T2 using a mapping function.
// This function has two type parameters, T1 and T2.
// This works with slices of any type.
func Map[T1, T2 any](s []T1, f func(T1) T2) []T2 {
	r := make([]T2, len(s))
	for i, v := range s {
		r[i] = f(v)
	}
	return r
}

// Reduce reduces a []T1 to a single value using a reduction function.
func Reduce[T1, T2 any](s []T1, initializer T2, f func(T2, T1) T2) T2 {
	r := initializer
	for _, v := range s {
		r = f(r, v)
	}
	return r
}

// Filter filters values from a slice using a filter function.
// It returns a new slice with only the elements of s
// for which f returned true.
func Filter[T any](s []T, f func(T) bool) []T {
	var r []T
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}
```
Here's how we might use these functions:
```go
func main() {
	words := []string{"One", "Potato", "Two", "Potato"}
	filtered := Filter(words, func(s string) bool {
		return s != "Potato"
	})
	fmt.Println(filtered) // [One Two]
	lengths := Map(filtered, func(s string) int {
		return len(s)
	})
	fmt.Println(lengths) // [3 3]
	sum := Reduce(lengths, 0, func(acc int, val int) int {
		return acc + val
	})
	fmt.Println(sum) // 6
}
```
The complete example code can be found [here](./examples/mapFilterReduce/main.go)

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
