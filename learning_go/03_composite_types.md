# Composite Types

<!--toc:start-->
- [Composite Types](#composite-types)
  - [Arrays-Too Rigid to Use Directly](#arrays-too-rigid-to-use-directly)
  - [Slices](#slices)
    - [len](#len)
    - [append](#append)
    - [Capacity](#capacity)
      - [The Go Runtime](#the-go-runtime)
    - [make](#make)
    - [Emptying A Slice](#emptying-a-slice)
    - [Declaring Your Slice](#declaring-your-slice)
    - [Slicing Slices](#slicing-slices)
      - [Slices share storage sometimes](#slices-share-storage-sometimes)
    - [copy](#copy)
    - [Converting Arrays to Slices](#converting-arrays-to-slices)
    - [Converting Slices to Arrays](#converting-slices-to-arrays)
  - [Strings and Runes and Bytes](#strings-and-runes-and-bytes)
  - [Maps](#maps)
    - [Reading and Writing a Map](#reading-and-writing-a-map)
    - [The comma ok idiom](#the-comma-ok-idiom)
    - [Deleting from Maps](#deleting-from-maps)
    - [Emptying a Map](#emptying-a-map)
    - [Comparing Maps](#comparing-maps)
    - [Using Maps as Sets](#using-maps-as-sets)
  - [Structs](#structs)
    - [Anonymous Structs](#anonymous-structs)
    - [Comparing and Converting Structs](#comparing-and-converting-structs)
  - [Exercises](#exercises)
  - [Wrapping Up](#wrapping-up)
<!--toc:end-->

## Arrays-Too Rigid to Use Directly
Go has arrays like most languages but they are rarely used directly. All the elements in the array must be of the specified type and there are a fed different declaration styles:
```go
var a [3]int // creates an array of three ints, all values are initialized to the zero value
var b = [3]int{10, 20, 30} // Initial values can be specified with an array literal 
var c = [12]int{1, 5: 4, 6, 10: 100, 15} // [1, 0, 0, 0, 0, 4, 6, 0, 0, 0, 100, 15]
var d = [...]int{10, 20, 30} // can use ... when using an array literal
fmt.Println(b == d) // prints true; you can use == and != to compare arrays
var e [2][3]int // can simulate multidimensional arrays, does not have true matrix support
```
Like most languages, arrays in Go are read and written using bracket syntax:
```go
x[0] = 10
fmt.Println(x[2])
```
You cannot read or write past the end of an array or use a negative index. If you do this with a constant or literal index it is a compile time error. Out-of-bounds read or write with a variable index compiles but fails with a runtime *panic*(more on panics later).

The built-in function `len` takes in an array an returns its length:
```go
fmt.Println(len(x))
```
The reason arrays are rarely used directly is because the *size* of the array is considered to be part of the *type* of the array. This means an array declared to be `[3]int` is a different type than `[4]int`. It also means that you cannot use a variable to specify the size of an array because types must be resolved at compile time.

You also can't use type conversion to convert arrays of different sizes to identical types. Because of this you can't write a function that works with arrays of any size and you can't assign arrays of different sizes to the same variable.

For these reasons only use arrays when you know the exact length you need ahead of time. The main reason why arrays exist in Go is to provide the backing store for *slices*.

## Slices
Most of the time when you want a data structure to hold a sequence of values you will use a slice. For slices, the length is not part of the type for a slice, removing the limitation of arrays. A single function can process a slice of any size and we can grow slices as needed.

Working with slices looks very similar to working with arrays with a few key differences. First, notice that we don't specify the size of the slice when we declare it:
```go
var a = []int{10, 20, 30}
var b = []int{1, 5: 4, 6, 10: 100, 15} // sparse slice
var c [][]int // simulate multidimensional slice
a[0] = 20 // reads and writes are the same as arrays
fmt.Println(a[2])
```
We see some differences when declaring a slice without using a literal like so: `var x []int`. This creates a slice of `int`s. Since no value is assigned `x` is assigned the zero value for a slice which is `nil`. In Go `nil` is an identifier that represents the lack of a value for some types. Like untyped numeric constants `nil` has no type so it can be assigned or compared against values of different types.

A slice is the first type we've covered that isn't *comparable*. It is a compile time error to use `==` or `!=` to compare slices. The only thing you can compare a slice with is `nil`.
```go
var x []int
fmt.Println(x == nil) // true
```
Since Go 1.21, the `slices` package in the standard library includes two functions to compare slices. The `slices.Equal` function takes int two slices and returns true if the slices are the same length, and all of the elements are equal. It requires that the elements of the slices are comparable. The other function, `slices.EqualFunc` lets you pass in a function to determine equality and does not require the slice elements to be comparable.
```go
x := []int{1, 2, 3, 4, 5}
y := []int{1, 2, 3, 4, 5}
z := []int{1, 2, 3, 4, 5, 6}
s := []string{"a", "b", "c"}
fmt.Println(slices.Equal(x, y)) // true
fmt.Println(slices.Equal(x, z)) // false
fmt.Println(slices.Equal(x, s)) // does not compile
```
### len
Go provides several built-in functions to work with its built-in types, `len` being one of them. We've already seen how `len` works with arrays and it works for slices too. If you pass a `nil` slice to `len` it returns 0.

> **_NOTE:_** Functions like `len` are built in to Go because they can do things that can't be done by functions that you can write. `len`'s parameter can be any type of array or any type of slice. It also works for strings, maps and channels. Trying to pass a variable of any other type to `len` is a compile-time error.

### append
The built-in `append` function is used to grow slices:
```go
var x []int
x = append(x, 10)
var y []int{1, 2, 3, 4}
y = append(y, 5)
y = append(y, 6, 7, 8) // can append multiple values at once
```
The append function takes at least two parameters, a slice of any type and a value of that type. It returns a slice of the same type.

One slice is appended onto another by using the `...` operator to expand the source slice into individual values:
```go
var x []int
y := []int{20, 30, 40}
x = append(x, y...)
```
It is a compile-time error if you forget to assign the value returned from `append`. This is because Go is a *call by value* language. Every time you pass a parameter to a function Go makes a copy of the value that's passed in. Passing a slice to the `append` function actually passes a copy of the slice and returns the copy. You then assign the returned slice back to the variable in the calling function.

### Capacity
Each element in a slice is assigned to consecutive memory locations which makes it quick to read and write values. Every slice has a capacity which is the number of consecutive memory locations reserved, this can be larger than the length. Each time you append to a slice the length is increased by 1, and when length equals capacity there is no more room to put values. If you call `append` when length equals capacity the `append` function uses the Go runtime to allocate a new slice with a larger capacity. The values in the original slice are copied to the new slice, new values are added to the end, and the new slice is returned.

> #### The Go Runtime
> The Go runtime provides services like memory allocation and garbage collection, concurrency support, networking, and implementations of built-in types and functions.
>
> The Go runtime is compiled into every Go binary. This makes it easier to distribute Go programs and avoid worries about about compatibility issues between the runtime and the program.

Just as the built-in `len` function returns the current length of a slice, the built-in `cap` function returns the current capacity of a slice. This will mostly be used to check if a slice is large enough to hold new data, or if a call to `make` is needed to create a new slice.

### make
The built-in `make` function allows us to specify the type, length, and optionally, the capacity when creating a slice;
```go
x := make([]int, 5) // creates an int slice with a length of 5 and capacity of 5, all elements initialized to 0 
y := make([]int, 5, 10) // creates an int slice with length of 5 and capacity of 10
z := make([]int, 0, 10) // creates a non-nil slice with a length 0 but capacity of 10
z = append(z, 5, 6, 7, 8) // Since length is 0 can't directly index but we can append
```
### Emptying A Slice
Go 1.21 added a `clear` function that takes in a slice and sets all of the slice's elements to their zero value. The length of the slice remains unchanged:
```go
s := []string{"first", "second", "third"}
fmt.Println(s, len(s)) // [first second third] 3
clear(s)
fmt.Println(s, len(s)) // [   ] 3
```
### Declaring Your Slice
The primary goal when declaring a slice is to minimize the number of times the slice needs to grow. If it's possible that the slice won't need to grow at all(because your function might return nothing), use a `var` declaration with no assigned value to create a `nil` slice:
```go
var data []int
```
If you have some starting values or if the values aren't going to change then a slice literal is a good choice.
```go
data := []int{2, 4, 6, 8}
```
If you have a good idea of how large your slice needs to be, but don't know what those values will be use `make`. The next consideration is what to set the length and capacity to:
- If you are using a slice as a buffer then specify a nonzero length
- If you are sure you know the exact size you want you can specify the length and index into the slice to set values. This is often done when transforming values in one slice and storing them in a second.
- Otherwise use `make` with a zero length and a specified capacity. This allows you to use `append` to add items to the slice. This means you won't have any extra zero values at the end if the number of items is less than the initial capacity and your code won't panic if the number of items is larger than the initial capacity.

### Slicing Slices
A *slice expression* creates a slice from a slice. Here are a few examples:
```go
x := []int{1, 2, 3, 4}
y := x[:2]  // [1, 2]
z := x[1:]  // [2, 3, 4]
d := x[1:3] // [2, 3]
e := x[:]   // [1, 2, 3, 4]
```
#### Slices share storage sometimes
When you take a slice from a slice, you are *not* making a copy of the data. Instead, you have two variables that are sharing memory. This means that changes to an element in a slice affect all slices that share that element:
```go
x := []int{1, 2, 3, 4}
y := x[:2]
z := x[1:]
x[1] = 20
y[0] = 10
z[1] = 30
fmt.Println("x:", x) // [10, 20, 30, 4]
fmt.Println("y:", y) // [10, 20]
fmt.Println("z:", z) // [20, 30, 4]
```
Slicing slices can get extra confusing when using `append`:
```go
x := []int{1, 2, 3, 4}
y := x[:2]
fmt.Println(cap(x), cap(y)) // 4 4
y = append(y, 30)
fmt.Println(x) // [1, 2, 30, 4]
fmt.Println(y) // [1, 2, 30]
```
Whenever you take a slice from another slice, the sub-slice's capacity is set to the capacity of the original slice, minus the offset of the sub-slice within the original slice. This means any unused capacity in the original slice is also shared with any sub-slices.

To avoid complicated slice situations, you should either never use `append` with a sub-slice or make sure that `append` doesn't overwrite by using a *full slice expression*. The full slice expression includes a third part which indicates the last position in the parent slice's capacity that's available for the sub-slice:
```go
y := x[:2:2]
z := x[2:4:4]
```
### copy
If you need to create a slice that's independent from the original, use the built-in `copy` function:
```go
x := []int{1, 2, 3, 4}
y := make([]int, 4)
z := make([]int, 2)
numY := copy(y, x)
numZ := copy(z, x[2:])
fmt.Println(y, numY) // [1, 2, 3, 4] 4
fmt.Println(z, numZ) // [3, 4] 2
numX = copy(x[:3], x[1:])
fmt.Println(x, numX) // [2, 3, 4, 4] 3
```
The `copy` function takes two parameters, the first is the destination slice and the second is the source slice. It copies as many elements as it can from source to destination, limited by whichever slice is smaller, and returns the number of elements copied. The *capacity* of the slices doesn't matter, it's the length that is important.

You can also sue `copy` with arrays by taking a slice of the array. The array can be either the source or the destination of the copy:
```go
x := []int{1, 2, 3, 4}
d := [4]int{5, 6, 7, 8}
y := make([]int, 2)
copy(y, d[:])
fmt.Println(y) // [5 6]
copy(d[:], x)
fmt.Println(d) // [1 2 3 4]
```
### Converting Arrays to Slices
If you have an array, you can take a slice from it using a slice expression. This is a useful way to bridge an array to a function that only takes slices. Do be aware that taking a slice from an array has the same memory sharing properties as taking a slice from a slice:
```go
x := [4]int{5, 6, 7, 8}
y := x[:2]
z := x[2:]
x[0] = 10
fmt.Println("x:", x) // x: [10 6 7 8]
fmt.Println("y:", y) // y: [10 6]
fmt.Println("z:", z) // z: [7 8]
```
### Converting Slices to Arrays
When you convert a slice to an array, the data in the slice is copied to new memory. That means that changes tot he slice won't affect the array, and vice-versa.
```go
xSlice := []int{1, 2, 3, 4}
xArray := [4]int(xSlice)
smallArray := [2]int(xSlice)
xSlice[0] = 10
fmt.Println(xSlice) // [10 2 3 4]
fmt.Println(xArray) // [1 2 3 4]
fmt.Println(smallArray) // [1 2]
```
The size of the array must be specified at compile time. It is an error to use `[...]` in a slice to array type conversion.

Also note that while the array can be smaller than the slice, it cannot be bigger. If it is your code will panic at runtime if you specify an array size that's bigger than the length(not the capacity) of the slice:
```go
panicArray := [5]int(xSlice)
fmt.Println(panicArray) // panic: runtime error: cannot convert slice with length 4 to array or pointer to array with length 5
```
## Strings and Runes and Bytes
You might think that a string in Go is made out of runes, but that's not the case. Under the covers, Go uses a sequence of bytes to represent a string. These bytes don't have to be in any particular character encoding but several Go library functions(like `for-range`) assume that a string is composed of a sequence of UTF-8 encoded code points.

Just like you can extract a single value from an array or a slice, you can extract a single value from an array or a slice expression you can extract a single value from a string using an *index expression*:
```go
var s string = "Hello there"
var b byte = s[6] // t(116)
var s2 string = s[4:7] // "o t"
var s3 string = s[:5]  // "Hello"
var s4 string = s[6:]  // "there"
```

There is a problem with using slice notation to make substrings and index notation to extract individual values. A string is composed of a sequence of bytes, while a code point in UTF-8 can be anywhere from one to four bytes long. The previous examples were composed of code points that are one byte long in UTF-8 so everything worked. But when dealing with languages other than English and symbols you will run into code points that are multiple bytes long in UTF-8:
```go
var s string := "Hello ☀️"
var s2 string = s[4:7] // We don't get "o ☀️" since we only copied the first byte of the emojis code point
var s3 string = s[:5]  // "Hello"
var s4 string = s[6:]  // ☀️
```
Go allows you to pass a string to the built-in `len` function to find the length of a string, this returns the length in bytes, not in code points:
```go
var s string := "Hello ☀️"
fmt.Println(len(s)) // Prints 10 not 7 because the emoji takes 4 bytes
```
Because of the complicated relationship between runes, strings, and bytes, Go has some interesting type conversions between these types:
```go
var a rune = 'x'
var s string = string(a) // "x"
var b byte = 'y'
var s2 string = string(b) // "y"
var x int = 65
var y = string(x) // y is "A" not "65"
```
A string can be converted back and forth to a slice of bytes or a slice of runes like so:
```go
var s string := "Hello ☀️"
var bs []byte = []byte(s) // [72 101 108 108 111 44 32 240 159 140 158]
var rs []rune = []rune(s) // [72 101 108 108 111 44 32 127774]
```
Most data in Go is read and written as a sequence of bytes, so the most common string type conversions are back and forth with a slice of bytes.

Instead of using slice and index expressions with strings, you should extract substrings and code points from strings using the functions in the `strings` and `unicode/utf8` packages in the standard library.

## Maps
Go provides the `map` built-in data type for situations where you want to associate one value to another. There are a few different ways to declare maps:
```go
var nilMap map[string]int
```
`nilMap` is declared to be a map with `string` keys and `int` values. The zero value for a map is `nil`. A `nil` map has a length of 0. Attempting to read a `nil` map always returns the zero value for the map's value type. However, attempting to write to a `nil` map causes a panic.

We can use a `:=` declaration to create a map variable by assigning a map literal:
```go
totalWins := map[string]int{} // Different than a nil map, has a length of 0 but you can read and write to it
teams := map[string][]string { // Non-empty map literal
  "Orcas": []string{"Fred", "Ralph", "Bijou"},
  "Lions": []string{"Sarah", "Peter", "Billie"},
  "Kittens": []string{"Waldo", "Raul", "Ze"},
}
```
If you know how many key-value pairs you intend to put in the map but not the exact values you can use `make` to create a map with a default size:
```go
ages := make(map[int][]string, 10)
```
Maps created with `make` still have a length of 0 and they can grow past the initially specified size.

Maps are like slices in a number of ways:
- Maps automatically grow as you add key-value pairs to them.
- If you know how many key-value pairs you plant to insert into a map, you can use `make` to create a map with a specific initial size.
- Passing a map to the `len` function tells you the number of key value pairs in the map.
- The zero value for a map is `nil`
- Maps are not comparable. You can check if they are equal to `nil`, but you cannot check for equality using `==`.

The key for a map can be any comparable type. This means you cannot use a slice or a map as the key for a map.

### Reading and Writing a Map
```go
totalWins := map[string]int{}
totalWins["Orcas"] = 1
totalWins["Lions"] = 2
fmt.Println(totalWins["Orcas"]) // 1
fmt.Println(totalWins["Kittens"]) // 0
totalWins["Kittens"]++
fmt.Println(totalWins["Kittens"]) // 1
totalWins["Lions"] = 3
fmt.Println(totalWins["Lions"]) // 3
```
When we try to read the value assigned to a map key that was never set, the map returns the zero value for the map's value type. In this case, the value type is an `int`, so we get back a 0 which makes sense in this context.

### The comma ok idiom
We've seen that a map returns the zero value if you look up the value for a key that's not in the map. This can be useful but sometimes you need to find out if the key is actually in the map or not. Go provides the *comma ok idiom* to tell the difference between a key that's associated with a zero value and a key that's not in the map:
```go
m := map[string]int{
  "hello": 5,
  "world": 0,
}
v, ok := m["hello"]
fmt.Println(v, ok) // 5 true

v, ok = m["world"]
fmt.Println(v, ok) // 0 true

v, ok = m["goodbye"] // 0 false
fmt.Println(v, ok)
```
The result of map reads are assigned to two variables. The first gets the value associated with the key, the second value is a `bool` usually named `ok`. If `ok` is `true`, the key is present in the map and it is `false` if the key is not present.

### Deleting from Maps
Key-value pairs are removed from a map via the built-in `delete` function:
```go
m := map[string]int{
  "hello": 5,
  "world": 10,
}
delete(m, "hello")
```
The `delete` function takes a map and a key and remove the key-value pair with the specified key. If the key isn't present, nothing happens. The `delete` function doesn't return a value.

### Emptying a Map
The clear function we saw [earlier](#emptying-a-slice) works on maps too. However, unlike a cleared slice a cleared map has its length set to zero:
```go
m := map[string]int{
  "hello": 5,
  "world": 10,
}
fmt.Println(m) // map[hello:5 world:10] 2
clear(m)
fmt.Println(m) // map[] 0
```
### Comparing Maps
Go 1.21 added a package to the standard library called `maps` that contains helper functions for working with maps. We'll cover this package in more detail in a [later section](./08_generics.md#adding-generics-to-the-standard-library). Two of the function sin the package are useful for comparing maps, `maps.Equal` and `maps.EqualFunc`. They are analogous to the `slices.Equal` and `slices.EqualFunc` functions: 
```go
m := map[string]int{
  "hello": 5,
  "world": 10,
}
n := map[string]int{
  "world": 10,
  "hello": 5,
}
fmt.Println(maps.Equal(m, n)) // true
```
### Using Maps as Sets
Go doesn't include a set but you can use a map to simulate some of the features of a set. Use the key of the map for the type that you want to put into the set and use a `bool` for the value:
```go
intSet := map[int]bool{}
vals := []int{5, 10, 2, 5, 8, 7, 3, 9, 1, 2, 10}
for _, v := range vals {
  intSet[v] = true
}
fmt.Println(len(vals), len(intSet)) // 11 8
fmt.Println(intSet[5])   // true
fmt.Println(intSet[500]) // false
if intSet[2] {
  fmt.Println("2 is in the set")
}
```
Some people prefer to use `struct{}` for the value when a map is being used to implement a set. The advantage is that an empty struct uses zero bytes, while a boolean uses one byte:
```go
intSet := map[int]struct{}{}
vals := []int{5, 10, 2, 5, 8, 7, 3, 9, 1, 2, 10} 
for _, v := range vals {
  intSet[v] = struct{}{}
}
if _, ok := intSet[5]; ok {
  fmt.Println("5 is in the set")
}
```
Unless you have very large sets, it is unlikely that the difference in memory usage is significant enough to warrant this approach.

## Structs
Maps are convenient for storing some kinds of data but they have limitations. They don't define an API since there's no way to constrain a map to only allow certain keys. Also, all of the values in a map must be of the same type. When you have related data that you want to group together you should define a *struct*.

Here is what a struct in Go looks like:
```go
type person struct {
  name string
  age int
  pet string
}
```
A struct type is defined with the keyword `type`, the name of the struct, the keyword `struct`, a pair of braces, and within the braces you list the fields in the struct. A struct can be defined inside or outside of a function. A struct defined within a function can only be used within that function.

Once a struct type is declared, we can define variables of that type:
```go
var fred person // A zero value struct has every field set to the field's zero value
bob := person{} // Unlike maps there is no difference between an empty struct literal and not assigning a value
```
There are two different styles for nonempty struct literals:
```go
julia := person { // In this format a value for every field must be specified
  "Julia",        // and values are assigned to the fields in the order they
  40,             // where declared in the struct definition
  "cat",
}
beth := person { // map style declaration
  age: 30,
  name: "Beth",
}
```
A field in a struct is accessed with dotted notation:
```go
bob.name = "Bob"
fmt.Println(beth.name) // "Beth"
```
### Anonymous Structs
You can also declare that a variable implements a struct type without first giving the struct type a name. This is an *anonymous struct*:
```go
var person struct {
  name string
  age int
  pet string
}
person.name = "bob"
person.age = 50
person.pet = "dog"

pet := struct {
  name string
  kind string
}{
  name: "Fido",
  kind: "dog",
}
```
In this example, the type of the variables `person` and `pet` are anonymous structs. You assign and read fields in an anonymous struct the same way you would with a normal struct. And you can initialize an instance of a anonymous struct with a struct literal as well.

There are two common situations where anonymous structs are useful. The first is when translating external data into a struct or a struct into external data(like JSON or protocol buffers). This is called *unmarshaling* and *marshaling* data. Writing tests is another place where anonymous structs pop up.

### Comparing and Converting Structs
Whether or not a struct is comparable depends on its fields. Structs that are entirely composed of comparable types are comparable; those with slices, maps, or channel fields are not. You can, of course, write your own function that you use to compare structs.

Just like how Go doesn't allow comparisons between variables of different types, Go doesn't allow comparisons between variables that represent structs of different types. You can perform a type conversion from one struct type to another *if the fields of both structs have the same names, order, and types*:
```go
type firstPerson struct {
  name string
  age int
}
type secondPerson struct { // We can use a type conversion to convert an instance of firstPerson to secondPerson
  name string
  age int
}
```
An anonymous struct can be compared to any other struct without a type conversion if both struct fields have the same names, order, and types. You can even assign between named and anonymous struct types if they are comparable:
```go
type firstPerson struct {
  name string
  age int
}
f := firstPerson{
  name: "Bob",
  age: 50,
}

var g struct {
  name string
  age int
}
g = f // can use = and == between identical named and anonymous structs
fmt.Println(f == g) // true
```
## Exercises
1. Write a program that defines a variable named `greetings` of type slice of strings with the following values: "Hello", "Hola", "नमस्कार", "こんにちは", and "Привіт". Create a sub-slice containing the first two values, a second subslice with the second, third, and fourth values, and a third subslice with the fourth and fifth values. Print out all four slices.
2. Write a program that defines a string variable called `message` with the value `"Hi 👩 and 👨"` and prints the 4th rune in it as a character, not a number.
3. Write a program that defines a struct called `Employee` with three fields: `firstName`, `lastName`, and `id`. The first two fields are of type `string` and the last field (`id`) is of type `int`. Create three instances of this struct using whatever values you’d like. Initialize the first one using the struct literal style without names, the second using the struct literal style with names, and the third with a var declaration. Use dot notation to populate the fields in the third struct. Print out all three structs.

> [Solutions](./exercises/03/)

## Wrapping Up
This section covered how to use the built-in generic container types, slices, and maps as well as how to build our own composite types via structs. The [next section](./04_blocks_shadows_and_control_structures.md) will cover Go's control structures as well as how Go organizes code into blocks.
