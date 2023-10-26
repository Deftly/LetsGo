# Predeclared Types And Declarations

<!--toc:start-->
- [Predeclared Types And Declarations](#predeclared-types-and-declarations)
  - [The Predeclared Types](#the-predeclared-types)
    - [The Zero Value](#the-zero-value)
    - [Literals](#literals)
    - [Booleans](#booleans)
    - [Numeric Types](#numeric-types)
      - [Integer types](#integer-types)
      - [The special integer types](#the-special-integer-types)
      - [Integer operators](#integer-operators)
      - [Floating point types](#floating-point-types)
      - [Complex types(you're probably not going to use these)](#complex-typesyoure-probably-not-going-to-use-these)
    - [A Taste of Strings and Runes](#a-taste-of-strings-and-runes)
    - [Explicit Type Conversion](#explicit-type-conversion)
    - [Literals are Untyped](#literals-are-untyped)
  - [var Versus :=](#var-versus)
  - [Using const](#using-const)
  - [Typed and Untyped Constants](#typed-and-untyped-constants)
  - [Unused Variables](#unused-variables)
  - [Naming Variables and Constants](#naming-variables-and-constants)
  - [Exercises](#exercises)
  - [Wrapping Up](#wrapping-up)
<!--toc:end-->

## The Predeclared Types
Go has many built-in types, which are also called *predeclared types*(booleans, integers, floats, and strings). We'll cover each of these types and see how they work best in Go, but first we'll cover some concepts that apply to all types.

### The Zero Value
Go assigns a default *zero value* to any variable that is declared but not assigned a value. This removes a source of bugs found in C and C++ programs. As we cover each type we'll cover the zero value for the type.

### Literals
A literal represents a fixed value in source code. Literals are used to specify values that are to be used directly by the program or as initializers for variables. 

*Integer literals* are a sequence of numbers. They are base ten by default, but prefixes are used to indicate other bases: `0b` for binary, `0o` for octal, or `0x` for hexadecimal.

*Floating point literals* have a decimal point to indicate the fractional portion of the value. They can also have an exponent specified with the letter e and a positive or negative number(such as `6.03e23`)

*Rune literals* represent a character and is surrounded by single quotes. Rune literals can be written as single Unicode characters(`'a'`), 8-bit octal numbers(`'\141'`), 8-bit hexadecimal numbers(`'\x61'`), 16-bit hexadecimal numbers(`'\u0061'`), or 32-bit Unicode numbers(`'\U00000061'`). There are also several backslash escaped rune literals, the most useful ones being newline(`'\n'`), tab(`'\t'`), single quote(`'\''`), and backslash(`'\\'`).

There are two different ways to indicate *string literals*. You will mostly use double quotes to create an *interpreted string literal*(e.g. `"Greetings and Salutations"`). These contain zero or more rune literals. They are called "interpreted" because they interpret rune literals(both numeric and backslash escaped) into a single character.

If you need to include backslashes, double quotes, or newlines in your string without using escapes you can use a *raw string literal*. These are delimited with backquotes(```) and can contain any character except a backquote. Here is a raw string literal and an equivalent interpreted string literal:
```go
// Interpreted String Literal
"Greetings and\n\"Salutations\""
// Raw String Literal
`Greetings and
"Salutations"`
```
Literals are considered *untyped*. We'll talk about this more [later in this section](#literals-are-untyped). There are situations in Go where the type isn't explicitly declared. In those situations Go uses the *default type* for a literal.

### Booleans
The `bool` type represents boolean variables, they can have one of two values: `true` or `false`. The zero value for a `bool` is `false`: 
```go
var flag bool // no value assigned, set to false
var isAwesome = true
```
### Numeric Types
Go has 12 different numeric types(and a few special names/aliases) that are grouped into 3 categories(integer types, floating point types, and complex types).

#### Integer types
Go provides both signed and unsigned integers in a variety of sizes:
| Type   | Value range    |
|--------------- | --------------- |
| int8 | -128 to 127 |
| int16 | -32768 to 32767 |
| int32 | -2147483648 to 2147483647 |
| int64 | -9223372036854775808 to -9223372036854775807 |
| uint8 | 0 to 255 |
| uint16 | 0 to 65536 |
| uint32 | 0 to 4294967295 |
| uint64 | 0 to 18446744073709551615 |

The zero value for all integer types is 0.

#### The special integer types
A `byte` is an alias for `uint8`. It is legal to assign, compare, or perform mathematical operations between a `byte` and a `uint8`. You will rarely see `uint8` used, just use `byte` instead.

The second name is `int`. On a 32-bit CPU, `int` is a 32-bit signed integer like `int32`. On most 64-bit CPUs, `int` is a 64-bit signed integer like `int64`. Because `int` isn't consistent from platform to platform it is a compile time error to assign, compare, or perform mathematical operations between an `int` and an `int32` or `int64` without a type conversion. Integer literals default to being of type `int`.

The third special name `uint` which follows the same rules as `int`, only it is unsigned.

The two other special names for integer types are `rune` and `uintptr`.

> **_NOTE:_** Unless you need to be explicit about the size of an integer for performance or integration purposes, use the `int` type.

#### Integer operators
Go integers support the usual arithmetic operators: `+`, `-`, `*`, `/`, and `%` for modulus. The result of integer division is an integer, to get a floating point result you need to use type conversion to make your integers into floating point numbers.

Go also has bit-manipulation operators for integers. You can bit shift left and right with `<<` and `>>`, or do bit masks with `&`(logical AND), `|`(logical OR), `^`(logical XOR), and `&^`(logical AND NOT).

#### Floating point types
There are two floating point types in Go, `float32` and `float64`. They have a zero value of 0 and function similar to floating point numbers in other languages with a large range and limited precision. Floating point literals have a default type of `float64` and you should always use `float64` unless you need to be compatible with an existing format calls for `float32`.

Because floats aren't exact they can only be used in situations where inexact values are acceptable or the rules of floating point are well understood, limiting their use to things like graphics and scientific operations. Never use them to represent money or any other value that needs an exact decimal representation.

You can use all the standard mathematical and comparison operators with floats, except `%`. Floating point division has a couple interesting properties. Dividing a nonzero floating point variable by 0 returns `+Inf` or `-Inf`(positive or negative infinity) depending on the sign of the number. Dividing a floating point variable set to 0 by 0 returns `NaN`(Not a Number).

While Go allows you to use `==` and `!=` to compare floats you shouldn't due to the inexact nature of floats. Instead define a maximum allowed variance and see if the difference between two floats is less than that.

#### Complex types(you're probably not going to use these)
While Go does have complex numbers as a built-in type it is not a popular language for numerical computation. This is likely because features like matrix support are not part of the language and libraries have to use inefficient replacements. There has been discussion about removing complex numbers from a future version of Go but it's easier to just ignore the feature.

### A Taste of Strings and Runes
Like most languages Go includes strings as a built-in type and the zero value for a string is an empty string. You can put any Unicode character into a string. Like integers and floats they can be compared for equality using `==`, difference with `!=`, or ordering with `>`, `>=`, `<`, or `<=`.

Strings in Go are immutable, you can reassign the value of a string variable, but you cannot change the value of the string that is assigned to it.

Go also has a type that represents a single code point. The `rune` type is an alias for `int32` just like `byte` is an alias for `uint8`. As you could probably guess, a rune literals default type is a `rune`, and a string literal's default type is a `string`.

> **_NOTE:_** If you are referring to a character use the `rune` type not the `int32` type. They might be the same to the compiler but you want to use the type that clarifies the intent of your code.

### Explicit Type Conversion
As a language that values clarity of intent and readability, Go doesn't allow automatic type promotion between variables. You must use a *type conversion* when variable types do not match. Even different sized integers and floats must be converted to the same type to interact.
```go
var x int = 10
var y float64 = 30.2
var z float64 = float64(x) + y
var d int = x + int(y)
```
Since all type conversions in Go are explicit, you cannot treat another Go type as a boolean. Many languages allow a nonzero number or nonempty string to be interpreted as a boolean `true`. In Go *no other type can be converted to a `bool`, implicitly or explicitly*. To convert from another data type to a boolean, you must use one of the comparison operators.

### Literals are Untyped
While you can't add two integer variables together if they are declared to be of different types of integers, Go lets you use an integer literal in floating point expressions or even assign an integer literal to a floating point variable:
```go
var x float64 = 10
var y float64 = 200.3 * 5
```
This is because literals in Go are untyped, meaning they can be used with any variable whose type is compatible with the literal.

## var Versus :=

## Using const

## Typed and Untyped Constants

## Unused Variables

## Naming Variables and Constants

## Exercises

## Wrapping Up
