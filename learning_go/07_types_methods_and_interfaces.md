# Types, Methods, and Interfaces

<!--toc:start-->
- [Types, Methods, and Interfaces](#types-methods-and-interfaces)
  - [Types in Go](#types-in-go)
  - [Methods](#methods)
    - [Pointer Receivers and Value Receivers](#pointer-receivers-and-value-receivers)
    - [Code Your Methods for nil Instances](#code-your-methods-for-nil-instances)
    - [Methods Are Functions Too](#methods-are-functions-too)
    - [Functions Versus Methods](#functions-versus-methods)
    - [Type Declarations Aren't Inheritance](#type-declarations-arent-inheritance)
    - [Types Are Executable Documentation](#types-are-executable-documentation)
  - [iota Is for Enumerations-Sometimes](#iota-is-for-enumerations-sometimes)
  - [Use Embedding for Composition](#use-embedding-for-composition)
  - [Embedding Is Not Inheritance](#embedding-is-not-inheritance)
  - [A Quick Lesson on Interfaces](#a-quick-lesson-on-interfaces)
  - [Interfaces Are Type-Safe Duck Typing](#interfaces-are-type-safe-duck-typing)
  - [Embedding and Interfaces](#embedding-and-interfaces)
  - [Accept Interfaces, Return Structs](#accept-interfaces-return-structs)
  - [Interfaces and nil](#interfaces-and-nil)
  - [Interfaces are comparable](#interfaces-are-comparable)
  - [The Empty Interface Says Nothing](#the-empty-interface-says-nothing)
  - [Type Assertions and Type Switches](#type-assertions-and-type-switches)
  - [Use Type Assertions and Type Switches Sparingly](#use-type-assertions-and-type-switches-sparingly)
  - [Function Types Are a Bridge to Interfaces](#function-types-are-a-bridge-to-interfaces)
  - [Implicit Interfaces Make Dependency Injection Easier](#implicit-interfaces-make-dependency-injection-easier)
  - [Go Isn't Particularly Object-Oriented(and That's Great)](#go-isnt-particularly-object-orientedand-thats-great)
  - [Exercises](#exercises)
  - [Wrapping Up](#wrapping-up)
<!--toc:end-->

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
To understand the relationship between interfaces and `nil` requires understanding a little bit about how interfaces are implemented. In The Go runtime, interfaces are implemented as a struct with two pointer fields, one for the value and one for the type of the value. As long as the type field is non-nil, the interface is non-nil.

For an interface to be considered `nil` *both* the type and the value must be `nil`:
```go
var pointerCounter *Counter
fmt.Println(pointerCounter == nil) // true
var incrementer Incrementer
fmt.Println(incrementer == nil) // true
incrementer = pointerCounter
fmt.Println(incrementer == nil) // false
```
What `nil` indicates for a variable with an interface type is whether or not you can invoke methods on it. We saw earlier that you can invoke methods on `nil` concrete instances, so it makes sense that you can invoke methods on an interface that was assigned a `nil` concrete instance. If an interface is `nil`, invoking any methods on it triggers a panic.

Because an interface instance with a non-nil type is not equal to `nil`, it isn't straightforward to tell whether or not the value associated with the interface is `nil`. To find out we need to use reflection(we'll cover this [later](./16_here_be_dragons_reflect_unsafe_and_cgo.md#use-reflection-to-check-if-an-interfaces-value-is-nil)).

## Interfaces are comparable
Just as an interface is equal to `nil` only if its type and value fields are both `nil`, two instances of an interface type are only equal if their types are equal and their values are equal. But what happens if the type isn't comparable? Let's look at an example:
```go
type Doubler interface {
	Double()
}

type DoubleInt int

func (d *DoubleInt) Double() {
	*d = *d * 2
}

type DoubleIntSlice []int

func (d DoubleIntSlice) Double() {
	for i := range d {
		d[i] = d[i] * 2
	}
}
```
The `*DoubleInt` type is comparable(all pointer types are) and the `DoubleIntSlice` is not comparable(slices aren't comparable).
```go
func DoublerCompare(d1, d2 Doubler) {
	fmt.Println(d1 == d2)
}

func main() {
	var di DoubleInt = 10
	var di2 DoubleInt = 10
	dis := DoubleIntSlice{1, 2, 3}
	dis2 := DoubleIntSlice{1, 2, 3}
	DoublerCompare(&di, &di2) // false - the types match but the pointers point to different addresses
	DoublerCompare(&di, dis)  // false - the types do not match
	DoublerCompare(dis, dis2) // compiles but panics at runtime: runtime error: comparing uncomparable type main.DoubleIntSlice
}
```
Also be aware that the key of a map must be comparable, so a map can be defined to have an interface as a key:
```go
m := map[Doubler]int{}
```
Adding a key-value pair to this map where the key isn't comparable will also trigger a panic.

Given this behavior be careful using `==` and `!=` with interfaces or using an interface as a map key as this can easily generate a panic that will crash your program. To be extra-safe, you can use the `Comparable` method on `reflect.Value` to inspect an interface before using it with `==` or `!=`(We'll learn more about reflection [later on](./16_here_be_dragons_reflect_unsafe_and_cgo.md#reflection-lets-us-work-with-types-at-runtime)).

## The Empty Interface Says Nothing
Sometimes you need a way to say that a variable could store a value of any type. Go uses an *empty interface*, `interface{}`, to represent this:
```go
var i interface{}
i = 20
i = "hello"
i = struct {
  FirstName string
  LastName  string
}{"Fred", "Fredson"}
````
> **_NOTE:_** `interface{}` isn't special case syntax. An empty interface states that the variable can store any value whose type implements zero or more methods, which happens to match every type in Go.

To improve readability, Go added `any` as a type alias for `interface{}`.

Because an empty interface doesn't tell you anything about the value it represents, there isn't a lot you can do with it. A common use of `any` is a placeholder for data of uncertain schema that's read from an external source, like a JSON file:
```go
data := map[string]any{}
contents, err := os.ReadFile("testdata/sample.json")
if err != nil {
  return err
}
json.Unmarshal(contents, &data) // the contents are now in the data map
```
> **_NOTE:_** User-created data containers that were written before generics were added to Go use an empty interface to store a value(we'll cover generics in [next section](./08_generics.md)). Now that generics are part of Go, use them for any newly created data containers.

If you see a function that takes in an empty interface, it's likely that it is using reflection to either populate or read the value.

These situation should be relatively rare and you should avoid using `any`. Go is designed as a strongly typed language and attempts to work around this are unidiomatic. 

## Type Assertions and Type Switches
Go has two ways to see if a variable of an interface type has a specific concrete type or if the concrete type implements another interface. We'll start with *type assertions*.

A type assertion names the concrete type that implemented the interface, or names another interface that is also implemented by the concrete type whose value is stored in the interface:
```go
type MyInt int

func main() {
	var i any
	var mine MyInt = 20
	i = mine
	i2 := i.(MyInt)     // i2 is of type MyInt
	fmt.Println(i2 + 1) // 21
}
```
If the type assertion is wrong your code will panic like in this example:
```go
i2 := i.(string) // panic: interface conversion: interface {} is main.MyInt, not string
```
Go is very careful about concrete types. Even if two types share an underlying type, a type assertion must match the type of the value stored in the interface. The following will also panic:
```go
i2 := i.(int) // panic: interface conversion: interface {} is main.MyInt, not int
```
Crashing is not a desired behavior in most cases and we can avoid this with the comma ok idiom:
```go
i2, ok := i.(int) // if ok is set to false the other variable(i2) is set to its zero value
if !ok {
  return fmt.Errorf("unexpected type for %v", i)
}
fmt.Println(i2 + 1)
```
Even if you are absolutely certain that your type assertion is valid, use the comma ok idiom. You don't know how other people(or you in 6 months) will reuse your code. Sooner or later, your unvalidated type assertions will fail at runtime.

When an interface could be one of multiple possible types, use a *type switch* instead:
```go
func doThings(i any) {
	switch j := i.(type) {
	case nil:
		// i is nil, type of j is any
	case int:
		// j is of type int
	case MyInt:
		// j is of type MyInt
	case io.Reader:
		// j is of type io.Reader
	case string:
		// j is of type string
	case bool, rune:
		// i is either a bool or rune, so j is of type any
	default:
		// no idea what i is, so j is of type any
	}
}
```
This is similar to a normal `switch` statement but instead of specifying a boolean operation you specify a variable of an interface type and follow it with `.(type)`. Usually, you assign the variable being checked to another variable that's only valid within the `switch`.

> **_NOTE:_** Since the purpose of a type `switch` is to derive a new variable from an existing one, it is idiomatic to assign the variable being switched on to a variable of the same name: `i := i.(type)`. This is one of the few places where shadowing is a good idea.

## Use Type Assertions and Type Switches Sparingly
For the most part you should treat a parameter or return value as the type that was supplied and not what it could be. Otherwise, your function's API isn't accurately declaring what types it needs to perform its task. 

That being said, there are some cases where type assertions and type switches are useful. A common use of a type assertion is to see if the concrete type behind the interface also implements another interface. This allows you to specify optional interfaces.

The standard library uses this technique to allow more efficient copies when the `io.Copy` function is called:
```go
// copyBuffer is the actual implementation of Copy and CopyBuffer.
// if buf is nil, one is allocated.
func copyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error) {
    // If the reader has a WriteTo method, use it to do the copy.
    // Avoids an allocation and a copy.
    if wt, ok := src.(WriterTo); ok {
        return wt.WriteTo(dst)
    }
    // Similarly, if the writer has a ReadFrom method, use it to do the copy.
    if rt, ok := dst.(ReaderFrom); ok {
        return rt.ReadFrom(src)
    }
    // function continues...
}
```
One drawback to the optional interface technique arises when the decorator pattern is used to wrap other implementations of the same interface to layer behavior. If an optional interface is implemented by one of the wrapped implementations, you cannot detect it with a type assertion or type switch. 

As an example, the standard library includes a `bufio` package that provides a buffered reader. You can buffer any other `io.Reader` implementation by passing it to the `bufio.NewReader` function and using the returned `*bufio.Reader`. If the passed-in `io.Reader` also implemented `io.ReaderFrom`, wrapping it in a buffered reader prevents the optimization we saw in the previous example.

Type `switch` statements provide the ability to differentiate between multiple implementations of an interface that require different processing. This is most useful when there are only certain possible valid types that can be supplied for an interface:
```go
func walkTree(t *treeNode) (int, error) {
	switch val := t.val.(type) {
	case nil:
		return 0, errors.New("invalid expression")
	case number:
		return int(val), nil
	case operator:
		left, err := walkTree(t.lchild)
		if err != nil {
			return 0, err
		}
		right, err := walkTree(t.rchild)
		if err != nil {
			return 0, err
		}
		return val.process(left, right), nil
	default:
		return 0, errors.New("unknown node type")
	}
}
```
You can see the full implementation of this example [here](./examples/operationsTree/main.go)

## Function Types Are a Bridge to Interfaces
Go allows methods on *any* user-defined type including user-defined functions. This isn't just a corner case but actually quite useful as it allows function to implement interfaces. The most common use of this is for HTTP handlers. An HTTP handler processes an HTTP server request and is defined by an interface:
```go
type Handler interface {
  ServeHTTP(http.ResponseWriter, *http.Request)
}
```
By using a type conversion to `http.HandlerFunc`, any function that has the signature `func(http.ResponseWriter, *http.Request)` can be used as an `http.Handler`:
```go
type HandlerFunc func(http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  f(w, r)
}
```
This lets you implement HTTP handlers using functions, methods, or closures using the exact same code path as the one used for other types that meet the `http.Handler` interface.

The question then becomes, when should your function or method specify an input parameter of a function type and when should you use an interface?

If your single function is likely to depend on many other functions or other state that's not specified in its input parameters, use an interface parameter and define a function type to bridge a function to the interface. However if it's a simple function then a parameter of function type is a good choice.

## Implicit Interfaces Make Dependency Injection Easier
Over time developers will be asked to update programs to fix bugs, add features, and run in new environments. This means that you should structure your programs in ways that make them easier to modify. One way to do this is through *decoupling* code, so that changes to different parts of a program have no effect on each other.

One technique to ease decoupling is called *dependency injection*, which is the concept that your code should explicitly specify the functionality it needs to perform its task.

One of the benefits of Go's implicit interfaces is that they make dependency injection easy. Other languages of use large complicated frameworks to inject their dependencies, in Go this can be implemented without any additional libraries.

To demonstrate let's build a very simple web application(we'll cover Go's built in HTTP server support [later](./13_the_standard_library.md#the-server)). We'll start with a utility logger function:
```go
func LogOutput(message string) {
	fmt.Println(message)
}
```
Next we'll create a data store, a retrieval function and a factory function to create an instance of our data store:
```go
type SimpleDataStore struct {
	userData map[string]string
}

func (s SimpleDataStore) UserNameForID(userID string) (string, bool) {
	name, ok := s.userData[userID]
	return name, ok
}

func NewSimpleDataStore() SimpleDataStore {
	return SimpleDataStore{
		userData: map[string]string{
			"1": "Fred",
			"2": "Mary",
			"3": "Pat",
		},
	}
}
```
Next we'll want to write some business logic that looks up a user and says hello or goodbye. This business logic needs a data store and the ability to log when it is invoked. However, we don't want to force it to depend on `LogOutput` or `SimpleDataStore` since we might want to use a different logger or data store later on. To handle this we'll create interfaces to describe what our business logic needs:
```go
type DataStore interface {
	UserNameForID(userID string) (string, bool)
}

type Logger interface {
	Log(message string)
}
```
To make the `LogOutput` function meet this interface, we define a function type with a `Log` method on it:
```go
type LoggerAdapter func(message string)

func (lg LoggerAdapter) Log(message string) {
	lg(message)
}
```
Now our `LoggerAdapter` and `SimpleDataStore` meet the interfaces needed by our business logic, though neither type has any idea that it does. With the dependencies defined we can implement the business logic:
```go
type SimpleLogic struct {
	l  Logger
	ds DataStore
}

func (sl SimpleLogic) SayHello(userID string) (string, error) {
	sl.l.Log("in SayHello for " + userID)
	name, ok := sl.ds.UserNameForID(userID)
	if !ok {
		return "", errors.New("unknown user")
	}
	return "Hello, " + name, nil
}

func (sl SimpleLogic) SayGoodbye(userID string) (string, error) {
	sl.l.Log("in SayGoodbye for " + userID)
	name, ok := sl.ds.UserNameForID(userID)
	if !ok {
		return "", errors.New("unknown user")
	}
	return "Goodbye, " + name, nil
}
```
Notice that there is nothing in `SimpleLogic` that mentions the concrete types, so there's no dependency on them. There's no problem if we later swap in new implementations from an entirely different provider because the provider has nothing to do with out interface.

We'll also want a factory function to create an instance of `SimpleLogic`, passing in interfaces and returning a struct:
```go
func NewSimpleLogic(l Logger, ds DataStore) SimpleLogic {
	return SimpleLogic{
		l:  l,
		ds: ds,
	}
}
```
For our example we'll have a single endpoint,`/hello`, which says hello to the person whose user ID is supplied(Don't use query parameters in real applications for authentication information). Our controller needs business logic that says hello, so we define an interface for that:
```go
type Logic interface {
	SayHello(userID string) (string, error)
}
```
This method is available on our `SimpleLogic` struct, but once again, the concrete type is not aware of the interface. The other method on `SimpleLogic`, `SayGoodbye`, is not the interface because our controller doesn't care about it. The interface is owned by the client code, so its method set is customized to the needs of the client code:
```go
type Controller struct {
	l     Logger
	logic Logic
}

func (c Controller) SayHello(w http.ResponseWriter, r *http.Request) {
	c.l.Log("In SayHello")
	userID := r.URL.Query().Get("user_id")
	message, err := c.logic.SayHello(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(message))
}
```
Let's also a factory function for the `Controller`:
```go
func NewController(l Logger, logic Logic) Controller {
	return Controller{
		l:     l,
		logic: logic,
	}
}
```
Finally we'll wire up all of our components in `main` and start the server:
```go
func main() {
	l := LoggerAdapter(LogOutput)
	ds := NewSimpleDataStore()
	logic := NewSimpleLogic(l, ds)
	c := NewController(l, logic)
	http.HandleFunc("/hello", c.SayHello)
	http.ListenAndServe(":8080", nil)
}
```
You can find the complete code for this example [here](./examples/simpleWebApp/main.go)

The `main` function is the only part of the code that knows what all the concrete types actually are. If we want to swap implementations, this is the only place that we need to change. Externalizing the dependencies via dependency injection limits the changes that are needed to evolve our code and makes testing easier. We'll talk more about testing in a [later section](./15_writing_tests.md).

## Go Isn't Particularly Object-Oriented(and That's Great)
After looking at idiomatic use of types in Go we can see it's hard to categorize Go as a particular style of language. It isn't a strictly procedural language. It lacks method overriding, inheritance, and objects meaning it's not an object oriented language. Go has function types and closures, but it isn't a function language either.

Go is practical, it borrows concepts from many places in order to create a language that is simple, readable, and maintainable.

## Exercises

## Wrapping Up
This section covered types, methods, interfaces, and their best practices. The [next section](./08_generics.md) will cover generics.
