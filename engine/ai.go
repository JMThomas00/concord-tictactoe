package engine

import (
	"math/rand"
	"strconv"
)

// BestMove picks the computer's move. Level 1 plays at random, 2 wins or
// blocks when it can, 3 never loses.
func BestMove(g *Game, level int) string {
	var free []int
	for i, c := range g.cells {
		if c == 0 {
			free = append(free, i)
		}
	}
	if len(free) == 0 {
		return ""
	}
	pick := func(i int) string { return strconv.Itoa(i + 1) }
	switch {
	case level <= 1:
		return pick(free[rand.Intn(len(free))])
	case level == 2:
		me := g.moves%2 + 1
		for _, who := range []int{me, 3 - me} { // win first, then block
			for _, i := range free {
				g.cells[i] = who
				won := g.winner() >= 0
				g.cells[i] = 0
				if won {
					return pick(i)
				}
			}
		}
		return pick(free[rand.Intn(len(free))])
	}
	// Perfect play, choosing at random among equally good moves so games
	// against it don't all look the same.
	var best []int
	bestScore := -2
	for _, i := range free {
		g.cells[i] = g.moves%2 + 1
		g.moves++
		score := -minimax(g)
		g.moves--
		g.cells[i] = 0
		switch {
		case score > bestScore:
			best, bestScore = []int{i}, score
		case score == bestScore:
			best = append(best, i)
		}
	}
	return pick(best[rand.Intn(len(best))])
}

// minimax scores the position for the side to move: 1 win, 0 draw, -1 loss.
func minimax(g *Game) int {
	if g.winner() >= 0 {
		return -1 // the previous move won
	}
	if g.moves == 9 {
		return 0
	}
	best := -2
	for i := range g.cells {
		if g.cells[i] != 0 {
			continue
		}
		g.cells[i] = g.moves%2 + 1
		g.moves++
		if s := -minimax(g); s > best {
			best = s
		}
		g.moves--
		g.cells[i] = 0
	}
	return best
}
