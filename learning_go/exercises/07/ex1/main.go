package main

type Team struct {
	Name    string
	Players []string
}

type League struct {
	Teams map[string]Team
	Wins  map[string]int
	Name  string
}
