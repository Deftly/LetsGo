package main

import (
	"io"
	"os"
	"sort"
)

type Team struct {
	Name    string
	Players []string
}

type League struct {
	Teams map[string]Team
	Wins  map[string]int
	Name  string
}

func (l *League) MatchResult(team1 string, score1 int, team2 string, score2 int) {
	if _, ok := l.Teams[team1]; !ok {
		return
	}
	if _, ok := l.Teams[team2]; !ok {
		return
	}
	if score1 == score2 {
		return
	}
	if score1 > score2 {
		l.Wins[team1]++
	} else {
		l.Wins[team2]++
	}
}

func (l League) Ranking() []string {
	names := make([]string, 0, len(l.Teams))
	for k := range l.Teams {
		names = append(names, k)
	}
	sort.Slice(names, func(i, j int) bool {
		return l.Wins[names[i]] > l.Wins[names[j]]
	})
	return names
}

type Ranker interface {
	Ranking() []string
}

func RankPrinter(r Ranker, w io.Writer) {
	results := r.Ranking()
	for _, v := range results {
		io.WriteString(w, v)
		w.Write([]byte("\n"))
	}
}

func main() {
	l := League{
		Name: "Big League",
		Teams: map[string]Team{
			"USA": {
				Name:    "USA",
				Players: []string{"Player1", "Player2", "Player3", "Player4", "Player5"},
			},
			"Canada": {
				Name:    "Canada",
				Players: []string{"Player1", "Player2", "Player3", "Player4", "Player5"},
			},
			"Serbia": {
				Name:    "Serbia",
				Players: []string{"Player1", "Player2", "Player3", "Player4", "Player5"},
			},
			"Germany": {
				Name:    "Germany",
				Players: []string{"Player1", "Player2", "Player3", "Player4", "Player5"},
			},
		},
		Wins: map[string]int{},
	}
	l.MatchResult("USA", 50, "Canada", 70)
	l.MatchResult("Serbia", 85, "Germany", 80)
	l.MatchResult("USA", 60, "Serbia", 55)
	l.MatchResult("Canada", 100, "Germany", 110)
	l.MatchResult("USA", 65, "Germany", 70)
	l.MatchResult("Canada", 95, "Serbia", 80)
	RankPrinter(l, os.Stdout)
}
