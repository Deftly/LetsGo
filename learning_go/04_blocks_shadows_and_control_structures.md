# Blocks, Shadows, And Control Structures

<!--toc:start-->
- [Blocks, Shadows, And Control Structures](#blocks-shadows-and-control-structures)
  - [Blocks](#blocks)
  - [Shadowing Variables](#shadowing-variables)
    - [The Universe Block](#the-universe-block)
  - [if](#if)
  - [for, Four Ways](#for-four-ways)
    - [The Complete for Statement](#the-complete-for-statement)
    - [The Condition-Only for Statement](#the-condition-only-for-statement)
    - [The Infinite for Statement](#the-infinite-for-statement)
    - [break and continue](#break-and-continue)
    - [The for-range Statement](#the-for-range-statement)
      - [The `for-range` value is a copy](#the-for-range-value-is-a-copy)
    - [Labeling Your for Statements](#labeling-your-for-statements)
    - [Choosing The Right for Statement](#choosing-the-right-for-statement)
  - [switch](#switch)
  - [Blank Switches](#blank-switches)
  - [Choosing Between if and switch](#choosing-between-if-and-switch)
  - [goto-Yes, goto](#goto-yes-goto)
  - [Exercises](#exercises)
  - [Wrapping Up](#wrapping-up)
<!--toc:end-->

## Blocks
Each place where a declaration occurs is called a *block*. Variables, constants, types, and functions declared outside of any functions are placed in the *package* block. Import statements define names for other packages that are valid for the file that contains the `import` statement. These names are in the *file* block. Each function defines its own block and within a function every set of braces(`{}`) defines another block. You can access an identifier in any outer block from within any inner block.

## Shadowing Variables
A shadowing variable is a variable that has the same name as a variable in a containing block. As long as the shadowing variable exists(is in scope), you cannot access a shadowed variable.
```go
func main() {
  x := 10
  if x > 5 {
    fmt.Println(x) // 10
    x := 5
    fmt.Println(x) // 5
  }
  fmt.Println(x) // 10
}
```
There are situations where you want to avoid using `:=` because it can make it unclear what variables are being used because it's very easy to shadow a variable when using `:=`. Remember, we can use `:=` to create and assign to multiple variables at once and not all of the variables have to be new for the assignment to be legal:
```go
func main() {
  x := 10
  if x > 5 {
    x, y := 5, 20
    fmt.Println(x, y) // 5 20
  }
  fmt.Println(x) // 10
}
```
Although there was an existing definition of `x` in the outer block, `x` was still shadowed within the `if` statement. This is because `:=` only reuses variables that are declared in the current block.

You also need to be careful not to shadow a package import:
```go
func main() {
  x := 10
  fmt.Println(x)
  fmt := "oops"
  fmt.Println(fmt) // error: fmt.Println undefined (type string has no field or method Println)
}
```
Notice that the problem isn't that we named our variable `fmt`, it is that we tried to access something that the local variable `fmt` didn't have. Once the local variable `fmt` is declared it shadows the package named `fmt` in the file block, making it impossible to use the `fmt` package for the rest of the `main` function.

### The Universe Block
There is actually one more block, the universe block, which contains all other blocks. Built-in types(like `int` and `string`), constants(`true` and `false`), and functions(like `make` or `close`) aren't keywords in Go, instead the are *predeclared identifiers* and they are defined in the universe block.

This means that they can be shadowed in other scopes:
```go
fmt.Println(true) // true
true := 10
fmt.Println(true) // 10
```
**Be very careful to never redefine any of the identifiers in the universe block**. Doing so can lead to some very strange behavior that can be difficult to debug.

## if
The most visible difference between `if` statements in Go and other languages is that you don't put parenthesis around the condition. Go also adds the ability to declare variables that are scoped to the condition and to both the `if` and `else` blocks:
```go
if n := rand.Int(10); n == 0 {
  fmt.Println("That's too low")
} else if n > 5 {
  fmt.Println("That's too big:", n)
} else {
  fmt.Println("That's a good number:", n)
}
```
Having this special scope lets you create variables that are available only where they are needed. Once the series of `if/else` statements ends, `n` is undefined.

## for, Four Ways
In Go, `for` is the only looping keyword in the language and it can be used in four different formats:
- A complete, C-style `for`
- A condition-only `for`
- An infinite `for`
- `for-range`

### The Complete for Statement
```go
for i := 0; i < 10; i++ {
  fmt.Println(i)
}
```
Note that you *must* use `:=` to initialize the variables, `var` is *not* legal here.

### The Condition-Only for Statement
Go allows you to leave off both the initialization and increment components in a `for` statement. That leaves us with a `for` statement that functions like a while loop in other languages. 
```go
i := 1
for i < 100 {
  fmt.Println(i)
  i = i * 2
}
```
### The Infinite for Statement
This version of the `for` statement does away with the condition too, giving us an infinite loop:
```go
func main() {
  for {
    fmt.Println("Hello")
  }
}
```
Press Ctrl-C to exit out of the program.

### break and continue
To get out of an infinite loop without using the keyboard or powering down the system we use `break` statements. When we hit a `break` statement we exit the loop immediately.

There is no Go equivalent of the `do` keyword found in other languages. If you want a loop to iterate a t least once the cleanest way is to use an infinite `for` loop that ends with an `if` statement:
```go
for {
  // things to do in the loop
  if !CONDITION {
    break
  }
}
```
Go also includes the `continue` keyword which skips over the rest of the body of a `for` loop and proceeds directly to the next iteration.

### The for-range Statement
The `for-range` format is for iterating over some of Go's built-in types(strings, arrays, slices, maps, and channels) and resembles iterators found in other languages.

> **_NOTE:_** You can only use a `for-range` loop to iterate over the built-in compound types and user-defined types that are based on them.

Lets' look at an example using `for-range` with a slice:
```go
evenVals := []int{2, 4, 6, 8, 10, 12}
for i, v := range evenVals {
  fmt.Println(i, v)
}
```
This outputs the following:
```shell
0 2
1 4
2 6
3 8
4 10
5 12
```
With a `for-range` we get two loop variables. The first is the position in the data structure being iterated, the second is the value at that position. The idiomatic name for the variables depends on what is being looped over. For arrays, slices, or strings, an `i`(for *index*) is commonly used. When iterating over a map, `k`(for *key*) is used instead.

If you don't need to use one of the loop variables you can use an underscore(`_`) as the variable's name. This tells Go to ignore the value:
```go
evenVals := []int{2, 4, 6, 8, 10, 12}
for _, v := range evenVals {
  fmt.Println(v)
}
```
> **_NOTE:_** Any time you are in a situation where there's a value being returned, but you want to ignore it, use an underscore to hide the value.

Here's another example using `for-range`:
```go
samples := []string{"hello", "apple_π!"}
for _, sample := range samples {
  for i, r := range sample {
    fmt.Println(i, r, string(r))
  }
  fmt.Println()
}
```
The output when iterating over "hello" has no surprises:
```shell
0 104 h
1 101 e
2 108 l
3 108 l
4 111 o
```
In the first column we have the index, in the second we have the numeric value of the letter, and in the third we have the numeric value of the letter converted to a string.

The result for "apple_π!" are more interesting:
```shell
0 97 a
1 112 p
2 112 p
3 108 l
4 101 e
5 95 _
6 960 π
8 33 !
```
Notice that the first column skips number 7, and the value at position 6 is 960 which is far larger than what can fit in a byte.

What we are seeing is special behavior from iterating over a string with a `for-range` loop. It iterates over *runes* not *bytes*. Whenever a `for-range` loop encounters a multi-byte rune in a string, it converts the UTF-8 representation into a single 32-bit number and assigns it to the value. The offset is incremented by the number of bytes in the rune.

> **_NOTE:_** Use a `for-range` loop to access the runes in a string in order. The key is the number of bytes from the beginning of the string and the type of the value is rune.

#### The `for-range` value is a copy
Each time the `for-range` loop iterates over your compound type, it *copies* the value from the compound type to the value variable. *Modifying the value variable will not modify the value in the compound type*:
```go
evenVals := []int{2, 4, 6, 8, 10, 12}
for _, v := range evenVals {
  v *= 2
}
fmt.Println(evenVals) // [2, 4, 6, 8, 10, 12]
```
When we cover goroutines in a later section you'll see that if you launch goroutines in a `for-range` loop, you need to be careful how you pass the index and value to the goroutines.

### Labeling Your for Statements
By default the `break` and `continue` keywords apply to the `for` loop that directly contains them. If you have nested `for` loops and you want to exit or skip over an iteration of an outer loop you can use a label:
```go
func main() {
  samples := []string{"hello", "apple_π!"}
outer:
  for _, sample := range samples {
    for i, r := range sample {
      fmt.Println(i, r, string(r))
      if r == 'l' {
        continue outer
      }
    }
    fmt.Println()
  }
}
```
Notice that the label `outer` is indented to the same level as the surrounding function. Labels are always indented like this to make them easier to notice. Here's the output of the above:
```shell
0 104 h
1 101 e
2 108 l
0 97 a
1 112 p
2 112 p
3 108 l
```
Nested `for` loops with labels are rare. They are most commonly used to implement algorithms similar to the pseudocode below:
```go
outer:
  for _, outerVal := range outerValues {
    for _, innerVal := range outerVal {
      // process innerVal
      if invalidSituation(innerVal) {
        continue outer
      }
    }
    // here we have code that runs only when all of the
    // innerVal values were successfully processed
  }
```
### Choosing The Right for Statement
When iterating over all the contents of an instance of one of the built-in types you usually want to use a `for-range` loop. It avoids a lot of the boilerplate code and is the best way to walk through a string since it properly gives you back runes instead of bytes.

The complete `for` loop is best for when you aren't iterating from the first element to the last in a compound type. Although, do note that this doesn't properly handle multi-byte characters in strings, there you would still use a `for-range`.

The other two `for` statement formats are less common, use the condition only `for` when you are looping based on a calculated value. And when using the infinite `for` you should always have a `break` somewhere within the body since it's rare that you want a to loop forever.

## switch
```go
words := []string{"a", "cow", "smile", "gopher", "octopus", "anthropologist"}
for _, word := range words {
  switch size := len(word); size {
  case 1, 2, 3, 4:
    fmt.Println(word, "is a short word!")
  case 5:
    wordLen := len(word)
    fmt.Println(word, "is exactly the right length:", wordLen)
  case 6, 7, 8, 9:
  default:
    fmt.Println(word, "is a long word!")
  }
}
```
Like an `if` statement you can declare a variable that's scoped to all of the branches of the switch statement. In this case we scoped the variable `size` to all of the cases in the `switch` statement. Here's the output of the code:
```shell
a is a short word!
cow is a short word!
smile is exactly the right length: 5
anthropologist is a long word!
```
By default, cases in `switch` statements in Go don't fall through. This means that if you have multiple values that should trigger the same logic you separate multiple matches with commas in the same case block.

In the above example we are switching on the value of an integer, but we can also switch on any type that can be compared with `==`, which includes all of the built-in types except slices, maps, channels, functions, and structs that contain fields of these types.

## Blank Switches
Another way to use `switch` statements is to not specify the value that you're comparing against, this is called a *blank switch*. A regular `switch` only allows you to check for equality, a blank `switch` allows you to use any boolean comparison for each case:
```go
words := []string{"hi", "salutations", "hello"}
for _, word := range words {
  switch wordLen := len(word); {
  case wordLen < 5:
    fmt.Println(word, "is a short word!")
  case wordLen > 10:
    fmt.Println(word, "is a long word!")
  default:
    fmt.Println(word, "is exactly the right length.")
  }
}
```
This gives us the following output:
```shell
hi is a short word!
salutations is a long word!
hello is exactly the right length.
```
## Choosing Between if and switch
As a matter of functionality there isn't a lot of difference between a series of `if/else` statements and a blank `switch` statement. Favor blank `switch` statements over `if/else` chains when you have multiple related cases. Using a `switch` makes the comparisons more visible and reinforces that they are a related set of concerns.

## goto-Yes, goto
In Go, a `goto` statement specifies a labeled line of code and execution jumps to it. Go forbids jumps that skip over variable declarations and jumps that go into an inner or parallel block:
```go
func main() {
  a := 10
  goto skip // error: goto skip jumps over declaration of b
  b := 20
skip:
  c := 30
  fmt.Println(a, b, c)
  if c > a {
    goto inner // error: goto inner jumps into block
  }
  if a < b {
  inner:
    fmt.Println("a is less than b")
  }
}
```
Most of the time you shouldn't use `goto`. Labeled `break` and `continue` statements allow you to jump out of deeply nested loops or skip iteration. The following example shows one of the few valid use cases:
```go
func main() {
  a := rand.Intn(10)
  for a < 100 {
    if a%5 == 0 {
      goto done
    }
    a = a*2 + 1
  }
  fmt.Println("do something when the loop completes normally")
done:
  fmt.Println("do complicated stuff no matter why we left the loop")
  fmt.Println(a)
}
```
You should try very hard to avoid using `goto` but in the rare instances where it makes code more readable it is an option.

## Exercises
1. Write a `for` loop that puts 100 random numbers between 0 and 100 into an `int` slice.
2. Loop over the slice you created in exercise 1. For each value in the slice, apply the following rules:
  - if the value is divisible by 2, print "Two!"
  - if the value is divisible by 3, print “Three!”
  - if the value is divisible by 2 and 3, print Six!”. Don’t print anything else.
  - otherwise, print “Never mind”
3. Start a new program. In `main`, declare an `int` variable called `total`. Write a `for` loop that uses a variable named `i` to iterate from 0 (inclusive) to 10 (exclusive). The body of the for loop should be:
```go
total := total + 1
fmt.Println(total)
```
After the `for` loop, print the value of `total`. What is printed out? What is the likely bug in this code?

> [Solutions]()

## Wrapping Up
This section covered blocks, shadowing, and control structures and how to use them properly. The [next section](/learning_go/05_functions.md) will cover functions.
