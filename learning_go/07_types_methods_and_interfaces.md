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

<!--TODO: Finish method set addendum-->
> **_NOTE:_**: Go's automatic conversion from pointer types to value types and vice-versa is syntactic sugar and is independent of the method set concept. Read the [addendum on method sets](./method_set_addendum.md) for more details on why the method set of pointer instances have both pointer and value receivers, but the method set of value instances only have value receiver methods.

Do not write getter and setter method for Go structs unless you need them to meet an interface([We'll cover them shortly](#a-quick-lesson-on-interfaces)) or you need to update multiple fields as a single operation/the update isn't a straightforward assignment. Reserve methods for business logic.

### Code Your Methods for nil Instances
In most languages if you try to invoke a method on a `nil` instance you get some sort of error. In Go if the method has a value receiver you'll get a panic since there is no value being pointed to. If the method has a pointer receiver, it can work if the method is written to handle the possibility of a `nil` instance.

Here's an implementation of a binary tree that takes advantage of `nil` values for the receiver:
```go
type IntTree struct {
	left, right *IntTree
	val         int
}

func (it *IntTree) Insert(val int) *IntTree {
	if it == nil {
		return &IntTree{val: val}
	}
	if val < it.val {
		it.left = it.left.Insert(val)
	} else if val > it.val {
		it.right = it.right.Insert(val)
	}
	return it
}

func (it *IntTree) Contains(val int) bool {
	switch {
	case it == nil:
		return false
	case val < it.val:
		return it.left.Contains(val)
	case val > it.val:
		return it.right.Contains(val)
	default:
		return true
	}
}

func main() {
	it := &IntTree{}
	it = it.Insert(5)
	it.Insert(3)
	it.Insert(10)
	it.Insert(2)
	fmt.Println(it.Contains(2))  // true
	fmt.Println(it.Contains(12)) // false
}
```
Pointer receivers work like pointer function parameters, it's a copy of the pointer that's passed into the method. Just like `nil` parameters passed to functions, if you change the copy of the pointer, you haven't changed the original. This means that you can't write a pointer receiver method that handles `nil` and makes the original pointer non-nil.

If your method has a pointer receiver and won't work for a `nil` receiver you have to decide how to handle it. You could treat it as a fatal flaw, in this case just let the code panic. If a nil receiver is something that is recoverable, check for `nil` and return an error(errors will be covered in [section 9](./09_errors.md))

### Methods Are Functions Too
Methods in Go are so much like functions that you can use a method in place of a function any time there's a variable or parameter of a function type.
```go
type Adder struct {
	start int
}

func (a Adder) AddTo(val int) int {
	return a.start + val
}

func main() {
	myAdder := Adder{start: 10}
	fmt.Println(myAdder.AddTo(5)) // 15

	// You can assign the method to a variable or pass it to a parameter
	// of type func(int) int. This is called a method value
	f1 := myAdder.AddTo
	fmt.Println(f1(10)) // 20

	// This is called a method expression. When using a method expression the
	// first parameter is the receiver for the method
	f2 := Adder.AddTo
	fmt.Println(f2(myAdder, 15)) // 25
}
```
We'll see how to leverage method values and method expressions when we look at dependency injection [later in this section](#implicit-interfaces-make-dependency-injection-easier).

### Functions Versus Methods
Any time your logic depends on values that are configured on startup or changed while your program is running, those values should be stored in a struct and that logic should be implemented as a method. If your logic only depends on the input parameters, then it should be a function.

### Type Declarations Aren't Inheritance
In addition to declaring types based on built-in Go types and struct literals, you can declare a user-defined type based on another user-defined type:
```go
type HighScore Score
type Employee Person
```
Declaring a type based on another type looks like inheritance, but that isn't the case in Go. You can't assign an instance of type `HighScore` to a variable of type `Score` or vice versa without a type conversion. Also, any methods defined on `Score` aren't defined on `HighScore`.
```go
// assigning untyped constants is valid
var i int = 300
var s Score = 100
var hs HighScore = 200
hs = s            // compilation error
s = i             // compilation error
s = Score(i)      // ok
hs = HighScore(s) // ok
```
User defined types whose underlying types are built-in types can be assigned literals and constants compatible with the underlying type, as well as be used with operators for those types.

### Types Are Executable Documentation
Types are documentation. They make code clearer by providing a name for a concept and describing the kind of data that is expected. It's clearer for someone reading your code when a method has a parameter of type `Percentage` than of type `int`, and it's harder for it to be invoked with an invalid value.

## iota Is for Enumerations-Sometimes
Go doesn't have an enumeration type, instead it has `iota`, which lets you assign an increasing value to a set of constants.

When using `iota` it's best practice to first define a type based on `int` that will represent all of the valid values. Then we define a `const` block to define a set of values for the type:
```go
type MailCategory int

const (
	Uncategorized MailCategory = iota // 0
	Personal                          // 1
	Spam                              // 2
	Social                            // 3
	Advertisement                     // 4
)
```
The first constant in the `cosnt` block has the type specified and its value is set to `iota`. The subsequent lines have neither the type or value assigned to it. The Go compiler will repeat the type assignment to all of the subsequent constants in the block, the value of `iota` increments for each constant in the block starting with `0`.

If you were to insert a new identifier in the middle of your list of literals, all of the subsequent identifiers will be renumbered. This can break your application if those constant represented values in another system or database. Given this limitation only use `iota`-based enumerations to differentiate between a set of values without caring what the value is behind the scenes. If the actual value matters you should specify it explicitly.

## Use Embedding for Composition
Go encourages code reuse via built-in support for composition and promotion:
```go
type Employee struct {
	Name string
	ID   string
}

func (e Employee) Description() string {
	return fmt.Sprintf("%s (%s)", e.Name, e.ID)
}

type Manager struct {
	Employee // This is an embedded field
	Reports []Employee
}

func (m Manager) FindNewEmployees() []Employee {
	// do business logic
}

func main() {
	m := Manager{
		Employee: Employee{
			Name: "Bob",
			ID:   "12345",
		},
		Reports: []Employee{},
	}
	fmt.Println(m.ID)            // 12345
	fmt.Println(m.Description()) // Bob (12345)
}
```
The `Manager` struct contains a field of type `Employee`, but no name is assigned to that field. This makes `Employee` and embedded field. Any fields or methods declared on an embedded field are *promoted* to the containing struct and can be invoked directly on it.

> **_NOTE:_** You can embed any type within a struct, not just another struct. This promotes the methods on the embedded type to the containing struct.

If the containing struct has fields or methods with the same name as an embedded field, you need to use the embedded field's type to refer to the obscured fields or methods:
```go
type Inner struct {
	X int
}

type Outer struct {
	Inner
	X int
}

func main() {
	o := Outer{
		Inner: Inner{
			X: 10,
		},
		X: 20,
	}
	fmt.Println(o.X)       // 20
	fmt.Println(o.Inner.X) // 10
}
```
## Embedding Is Not Inheritance
Many developers try to understand embedding by treating it as inheritance, don't do this. You cannot assign a variable of type `Manager` to a variable of type `Employee`. To access the `Employee` field in `Manager` you must do so explicitly:
```go
var eFail Employee = m        // compilation error: cannot use m (type Manger) as type Employee in assignment
var eOk Employee = m.Employee // ok
```
<!--TODO: Finish dynamic dispatch addendum---->
There is no *dynamic dispatch*(see [dynamic dispatch addendum](./dynamic_dispatch_addendum.md)) for concrete types in Go. The methods on an embedded field have no idea that they are embedded. If you have a method on an embedded field that calls another method on the embedded field, and the containing struct has a method of the same name, the method on the embedded field is invoked, not the method on the containing struct:
```go
type Inner struct {
	A int
}

func (i Inner) IntPrinter(val int) string {
	return fmt.Sprintf("Inner: %d", val)
}

func (i Inner) Double() string {
	return i.IntPrinter(i.A * 2)
}

type Outer struct {
	S string
	Inner
}

func (o Outer) IntPrinter(val int) string {
	return fmt.Sprintf("Outer: %d", val)
}

func main() {
	o := Outer{
		Inner: Inner{
			A: 10,
		},
		S: "Hello",
	}
	fmt.Println(o.Double()) // 20
}
```
While embedding one concrete type inside another doesn't allow you to treat the outer type as the inner type, the methods on an embedded field do count toward the *method set* of the containing struct. This is important when it comes to implementing an interface.

## A Quick Lesson on Interfaces
A their core, interfaces are simple. Like other user-defined types, you use the `type` keyword to declare them:
```go
// The Stringer interface in the fmt package
type Stringer interface {
  String() string
}
```
An interface literal lists the methods that must be implemented by a concrete type to meet the interface. These defined methods are the method set of the interface. We mentioned earlier that the method set of a pointer instance contains the methods defined with both pointer and value receivers, while the method set of a value instance only contains methods with value receivers. Here's an example:
```go
type Counter struct {
	lastUpdated time.Time
	total       int
}

func (c *Counter) Increment() {
	c.total++
	c.lastUpdated = time.Now()
}

func (c Counter) String() string {
	return fmt.Sprintf("total: %d, last updated: %v", c.total, c.lastUpdated)
}

type Incrementer interface {
	Increment()
}

func main() {
	var myStringer fmt.Stringer
	var myIncrementer Incrementer
	pointerCounter := &Counter{}
	valueCounter := Counter{}

	myStringer = pointerCounter    // ok
	myStringer = valueCounter      // ok
	myIncrementer = pointerCounter // ok
	myIncrementer = valueCounter   // compile-time error
	// cannot use valueCounter (variable of type Counter) as Incrementer value in
  // assignment: Counter does not implement Incrementer (method Increment has pointer receiver)

	fmt.Println(myStringer, myIncrementer)
}
```
Like other types, interfaces can be declared in any block. Also, the convention is to name them with "er" endings, here are some examples: `io.Reader`, `io.Closer`, `io.ReadCloser`, `json.Marshaler`, `http.Handler`.

## Interfaces Are Type-Safe Duck Typing
What makes Go's interfaces special is that they are implemented *implicitly*, concrete types do not need to declare that they implement an interface. If the method set for a concrete type contains all of the methods in the method set for an interface, the concrete type implements the interface.
<!--TODO: Finish type-safety addendum-->
This implicit behavior enables both [type-safety] and decoupling, bridging the functionality in both static and dynamic languages.

Dynamically typed languages like Python, Ruby, and JavaScript don't have interfaces. Instead they use "duck typing", which is comes from the expression "If it walks like a duck and quacks like a duck, it's a duck.". The concept is that you can pass an instance of a type as a parameter as long as the function can find a method to invoke that it expects.

Interfaces can be shared, we've seen examples in the standard library that are used for input and output. If you write your code to work with `io.Reader` and `io.Writer`, it will function correctly whether it is writing to a file on local disk or a value in memory.

Standard interfaces also encourage the *decoration pattern*. It's common in Go to write factory functions that take in an instance of an interface and return another type that implements the same interface. Say we have the following function definition:
```go
func process(r io.Reader) error
```
We can process data from a file with the following:
```go
r, err := os.Open(fileName)
if err != nil {
  fmt.Println(err)
}
defer r.Close()
return process(r)
```
The `os.File` instance returned by `os.Open` meets the `io.Reader` interface and can be used in any code that read in data. If the file is gzip-compressed, we can wrap the `io.Reader` in another `io.Reader`:
```go
r, err := os.Open(fileName)
if err != nil {
  fmt.Println(err)
}
defer r.Close()
gz, err := gzip.NewReader(r)
if err != nil {
  fmt.Println(err)
}
defer gz.Close()
return process(r)
```
Now the code that was reading from an uncompressed file is reading from a compressed file instead.

> **_NOTE:_** If there's an interface in the standard library that describes what your code needs, use it.

It's fine for a type that meets an interface to specify additional methods that aren't part of the interface.

## Embedding and Interfaces
You can embed an interface in an interface. For example, the `io.ReadCloser` interface is built out of an `io.Reader` and an `io.Closer`:
```go
type Reader interface {
	Read(p []byte) (n int, err error)
}

type Closer interface {
	Close() error
}

type ReadCloser interface {
	Reader
	Closer
}
```
## Accept Interfaces, Return Structs
We've already covered why functions should accept interfaces: they make your code more flexible and explicitly declare exactly what functionality is being used.

The primary reason why your functions should return concrete types is they make it easier to update a function's return values in new versions of your code. When you return a concrete type, new methods and fields can be added without breaking existing code that calls the function, because the new fields and methods are ignored. However, if you add a new method to an interface all existing implementations of that interface must be updated or your code breaks.

Instead of writing a single factory function that returns different instances behind an interface based on input parameters, try to write separate factory functions for each concrete type. In some situations(such as a parser that returns different kinds of token), it's unavoidable and you have no choice but to return an interface.

Errors are the exception to this rule. Go functions and methods often declare a return parameter of the `error` interface type. With `error` it's quite likely that different implementations of the interface could be returned, so you need to use an interface to handle all possible options, as interfaces are the only abstract type in Go.

The one potential drawback to this pattern of accepting interfaces and returning structs has to do with heap allocations. We [previously covered](./06_pointers.md#reducing-the-garbage-collectors-workload) how reducing heap allocations improves performance. Returning a struct avoid a heap allocation which is good. However, when invoking a function with parameters of interface types, a heap allocation occurs for each of the interface parameters. Figuring out the trade-off between better abstraction and better performance should be done over the life of your program. Always write your code so that it is readable and maintainable. If you find that your program is too slow *and* you have profiled it *and* you have determined that the performance problems are due to heap allocations caused by an interface parameter, then you should rewrite the function to use a concrete type parameter.

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
