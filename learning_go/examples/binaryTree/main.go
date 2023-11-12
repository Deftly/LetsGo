package main

import "fmt"

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
	fmt.Println(it.Contains(3))  // true
	fmt.Println(it.Contains(5))  // true
	fmt.Println(it.Contains(10)) // true
	fmt.Println(it.Contains(12)) // false
}
