package main

import "fmt"

func main() {
	message := "Hi 👩 and 👨"
	runes := []rune(message)
	fmt.Println(runes)
	fmt.Println(string(runes[3]))
}
