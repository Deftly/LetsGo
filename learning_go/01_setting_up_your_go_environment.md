# Setting Up Your Go Environment

<!--toc:start-->
- [Setting Up Your Go Environment](#setting-up-your-go-environment)
  - [Installing the Go Tools](#installing-the-go-tools)
    - [Go Tooling](#go-tooling)
  - [Your First Go Program](#your-first-go-program)
    - [Making a Go Module](#making-a-go-module)
    - [go build](#go-build)
    - [go fmt](#go-fmt)
  - [go vet](#go-vet)
  - [Makefiles](#makefiles)
  - [The Go Compatibility Promise](#the-go-compatibility-promise)
  - [Staying Up to Date](#staying-up-to-date)
  - [Wrapping Up](#wrapping-up)
<!--toc:end-->

## Installing the Go Tools
To build Go code you need to download and install the Go development tools, the latest version can be found on the [Go website](https://go.dev/dl/). The installers for Mac and Windows automatically installs Go in the right location, remove old installs, and puts the Go binary in the default executable path.

The various Linux and BSD installers are gzipped tar files and expand to a directory named *go*. Copy this directory to */usr/local* and add */usr/local/go/bin* to you r `$PATH` to access the `go` command:
```bash
$ sudo rm -rf /usr/local/go
$ sudo tar -C /usr/local -xzf go<version>.linux-amd64.tar.gz
$ echo 'export PATH=$PATH:/usr/local/go/bin' >> $HOME/.bash_profile
$ source $HOME/.bash_profile
```
> **_NOTE:_** Go programs compile to a single native binary and do not require any additional software to run them, making it very easy to distribute programs written in Go.

You can validate that your environment is set up correctly with the following command:
```bash
$ go version
$ go version go1.20.10 linux/amd64
```
### Go Tooling
All of the Go development tools are accessed via the `go` command. In addition to `go version`, there's a compiler(`go build`), code formatter(`go fmt`), dependency manager(`go mod`), test runner(`go test`), a tool to scan for common mistakes(`go vet`) and more.

## Your First Go Program
Now let's go over the basics of writing a Go program.

### Making a Go Module
The first thing we need to do is create a directory to hold our program:
```bash
$ mkdir ch1
$ cd ch1
```
Inside the directory run `go mod init` to mark this directory as a Go *module*:
```bash
$ go mod init hello_world
go: creating new go.mod: module hello_world
```
We'll learn more about what a module is in [section 10](./10_modules_packages_and_imports.md), for now just remember that a Go project is called a module. A module isn't just source code, it is also an exact specification of the dependencies of the code within the module. Every module has a `go.mod` file in its root directory, this create when we run `go mod init`. This is what our `go.mod` file looks like:
```
module hello_world

go 1.20
```
The `go.mod` file declares the name of the module, the minimum supported version of Go for the module, and any modules that your module depends on.

You shouldn't edit the `go.mod` file directly, instead use the `go get` and `go mod tidy` commands.

### go build
Now let's write our first program:
```go
package main

import "fmt"

func main() {
	fmt.Println("Hello, world!")
}
```
The first line is the package declaration. Within a Go module, code is organized into one or more packages. The `main` package in a go Module contains the code that starts a Go program.

The `import` statement lists the packages referenced in this file. We're using a function in the `fmt` package from the standard library so we list it here.

All Go programs start form the `main` function in the `main` package. In the body of the function we are calling the `Println` function in the `fmt` package with the argument `"Hello, World!"`.
```bash
$ go build -o hello
$ ./hello
Hello, world!
```
The `go build` command creates an executable, by default the name of the binary matches the name in the module declaration, but you use the `-o` flag to change the name. 

### go fmt
Most languages allow a great deal of flexibility in how code is formatted, Go does not. Enforcing a standard format makes it easier to write tools that manipulate source code and simplifies the compiler.

The Go development tools include `go fmt`, which automatically adjusts your code to match the standard format. Run it with:
```bash
$ go fmt ./...
```
Using `./...` tells a Go tool to apply the command to all the files in the current directly and all subdirectories.

## go vet
There is a class of bugs where the code is syntactically valid but quite likely incorrect. The `go vet` command can detect some of these kinds of errors. Let's add one such error to our program:
```go
fmt.Printf("Hello, %s!\n")
```
`fmt.Printf` is a function with a template for its first parameter, and values for the placeholders in the template as the remaining parameters. In our example `%s` is the placeholder we specified no values for it. This code will compile and run but it isn't correct and we can detect it with `go vet`:
```bash
$ go vet ./...
# hello_world
./hello.go:6:2: fmt.Printf format %s reads arg #1, but call has 0 args
```
Once `go vet` finds the bug we can easily fix it:
```go
fmt.Printf("Hello, %s!", "world")
```
While `go vet` catches several common programming errors, there are things that it cannot detect. Luckily there are third-party Go code quality tools that can help, we'll cover those in a [later section](./11_go_tooling.md#code-quality-scanners)

## Makefiles
Modern software development relies on repeatable, automatable builds that can be run by anyone, anywhere, at any time. The way to do this is to use some kind of script to specify your build steps, one possible solution is `make`. It lets developers specify a set of operations that are necessary to build a program and the order in which the steps must be performed. Let's create a file called `Makefile` in our directory with the following:
```
.DEFAULT_GOAL := build

.PHONY:fmt vet build
fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build
```
Each possible operation is called a *target*. The `.DEFAULT_GOAL` defines which target to run if no target is specified. Next we have the target definitions. The word before the colon is the name of the target, any words after the target(like `vet` in the line `build: vet`) are the other targets that must run before the specified target runs. The tasks that are performed by the target are on the indented lines after the target.
```bash
$ make
go fmt ./...
go vet ./...
go build
```
Now entering a single command formats the code correctly, checks it for errors, and compiles it. You can also vet the code with `make vet` or just formatting with `make fmt`.

## The Go Compatibility Promise
There are periodic updates to the Go development tools, since Go 1.2 there has been a new release every 6 months with patch releases with bug and security fixes released as needed. The [Go Compatibility Promise](https://go.dev/doc/go1compat) is a detailed description of how the Go team plans to avoid breaking Go code. It says that there won't be any backward-breaking changes to the language or standard library for any Go version that starts with 1, unless the change is required for a bug or security fix.

This doesn't apply to the `go` commands. There have been backward-incompatible changes to the flags and functionality of the `go` commands, and it can happen again in the future.

## Staying Up to Date
Go programs compile to a standalone native binary, so you don't need to worry that updating your development environment could cause your currently deployed programs to fail. You can have programs compiled with different version of Go running simultaneously on the same computer or VM.

To update Go on Linux/BSD you need to download the latest version and follow these steps:
```bash
$ sudo mv /usr/local/go /usr/local/old-go
$ sudo tar -C /usr/local -xzf go<version>.tar.gz
$ sudo rm -rf /usr/local/old-go
```
## Wrapping Up
In this section we covered how to install the Go development environment and some of the tools for building Go programs. In the [next section](./02_predeclared_types_and_declarations.md) we'll look at the built-in types in Go and how to declare variables.
