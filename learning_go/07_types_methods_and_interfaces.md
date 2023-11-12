# Types, Methods, and Interfaces

## Types in Go
```go
type Person struct {
	FirstName string
	LastName  string
	Age       int
}
```
This should be read as declaring a user defined type with the name `Person` to have the *underlying type* of the struct literal that follows. You can also use any primitive type or compound type literal to define a concrete type:
```go
type Score int
type Converter func(string)Score
type TeamScores map[string]Score
```
Types can be declared at any block level, from the package block down, though they can only be accessed from within its scope. The exceptions are types exported from other packages.

## Methods
The methods for a type are defined at the package block level:
```go
func (p Person) String() string {
	return fmt.Sprintf("%s %s, age %d", p.FirstName, p.LastName, p.Age)
}
```
Method declarations are like function declarations with one addition: the *receiver* specification. By convention, the receiver name is a short abbreviation of the type's name, usually the first letter.

Another key difference between functions and methods is that methods can only be defined at the package block level, while functions can be defined inside any block.

We'll cover packages in [section 10](./10_modules_packages_and_imports.md), for now be aware that methods must be declared in the same package as their associated type, you can't add methods to types you don't control.

Here is how you invoke a method:
```go
p := Person{
  FirstName: "Sang-hyeok",
  LastName:  "Lee",
  Age:       27,
}
output := p.String()
```
### Pointer Receivers and Value Receivers
We saw in the [previous section](./06_pointers.md) that we use parameters of pointer type to indicate that a parameter might be modified by a function. The same rules apply for method receivers. The can be *pointer receivers* or *value receivers*, here are some rules for picking which kind of receiver to use:
- If your method modifies the receiver, you *must* use a pointer receiver.
- If your method needs to handle `nil` instances(see [Code Your Methods for nil Instances](#code-your-methods-for-nil-instances)), then it *must* use a pointer receiver.
- If your method doesn't modify the receiver, you *can* use a value receiver.

When a type has *any* pointer receiver methods, a common practice is to be consistent and use pointer receiver methods for *all* methods, even those that don't modify the receiver.

Here's a simple example:
```go
type Counter struct {
	lastUpdate time.Time
	total      int
}

func (c *Counter) Increment() {
	c.total++
	c.lastUpdate = time.Now()
}

func (c Counter) String() string {
	return fmt.Sprintf("total: %d, last updated: %v", c.total, c.lastUpdate)
}

func main() {
	var c Counter
	fmt.Println(c.String()) // total: 0, last updated: 0001-01-01 00:00:00 +0000 UTC
	c.Increment()
	fmt.Println(c.String()) // total: 1, last updated: 2023-11-11 18:17:27.104062603 -0800 PST m=+0.000073011
}
```
You might have noticed that we were able to call the pointer receiver method even though `c` is a value type. This is because Go automatically takes the address of the local variable when calling the method. In this case, `c.Increment()` is converted to `(&c).Increment()`. Similarly, if you call a value receiver on a pointer variable, Go automatically dereferences the pointer when calling the method.

Go considers both pointer and value receiver methods to be in the *method set* for a pointer instance. For a value instance, only the value receiver methods are in the method set. This concept will be important when we talk about interfaces in just a little bit.

> **_NOTE:_**: Go's automatic conversion from pointer types to value types and vice-versa is syntactic sugar and is independent of the method set concept. Read the [addendum on method sets](./method_set_addendum.md) for more details on why the method set of pointer instances have both pointer and value receivers, but the method set of value instances only have value receiver methods.

Do not write getter and setter method for Go structs unless you need them to meet an interface([We'll cover them shortly](#a-quick-lesson-on-interfaces)) or you need to update multiple fields as a single operation/the update isn't a straightforward assignment. Reserve methods for business logic.

### Code Your Methods for nil Instances


### Methods Are Functions Too

### Functions Versus Methods

### Type Declarations Aren't Inheritance

### Types Are Executable Documentation

## iota Is for Enumerations-Sometimes

## Use Embedding for Composition

## Embedding Is Not Inheritance

## A Quick Lesson on Interfaces

## Interfaces Are Type-Safe Duck Typing

## Embedding and Interfaces

## Accept Interfaces, Return Structs

## Interfaces and nil

## Interfaces are comparable

## The Empty Interface Says Nothing

## Type Assertions and Type Switches

## Use Type Assertions and Type Switches Sparingly

## Function Types Are a Bridge to Interfaces

## Implicit Interfaces Make Dependency Injection Easier

## Wire

## Go Isn't Particularly Object-Oriented(and That's Great)

## Exercises

## Wrapping Up
