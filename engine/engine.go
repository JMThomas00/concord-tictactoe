// Package engine is tic-tac-toe's rules: the board, whose turn it is, who
// won and how, and a computer opponent. It knows nothing about screens.
package engine

import (
	"fmt"
	"strconv"

	"github.com/JMThomas00/Concord/sdk/table"
)

// Game is a board of nine cells, numbered 1-9 left to right, top to bottom
// (like a phone's keypad). A move is the cell's number. X moves first.
type Game struct {
	cells [9]int // 0 empty, 1 X, 2 O
	moves int
	last  int // the last move's cell (0-8), or -1
}

// New starts an empty board.
func New() *Game { return &Game{last: -1} }

// Lines are the eight ways to get three in a row.
var Lines = [8][3]int{{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, {0, 3, 6}, {1, 4, 7}, {2, 5, 8}, {0, 4, 8}, {2, 4, 6}}

// WinningLine is the three cells that won, or nil.
func (g *Game) WinningLine() []int {
	for _, l := range Lines {
		if c := g.cells[l[0]]; c != 0 && c == g.cells[l[1]] && c == g.cells[l[2]] {
			return l[:]
		}
	}
	return nil
}

func (g *Game) winner() int {
	if l := g.WinningLine(); l != nil {
		return g.cells[l[0]] - 1
	}
	return -1
}

// Turn is 0 (X) or 1 (O), or -1 once the game is over.
func (g *Game) Turn() int {
	if g.Outcome().Over {
		return -1
	}
	return g.moves % 2
}

// Play takes a cell number, "1" to "9".
func (g *Game) Play(move string) error {
	n, err := strconv.Atoi(move)
	if err != nil || n < 1 || n > 9 {
		return fmt.Errorf("pick a cell from 1 to 9")
	}
	if g.Outcome().Over {
		return fmt.Errorf("the game is over")
	}
	if g.cells[n-1] != 0 {
		return fmt.Errorf("that cell is taken")
	}
	g.cells[n-1] = g.moves%2 + 1
	g.moves++
	g.last = n - 1
	return nil
}

// Outcome reports whether the game is over and who won.
func (g *Game) Outcome() table.Outcome {
	if w := g.winner(); w >= 0 {
		return table.Outcome{Over: true, Winner: w, Reason: "three in a row"}
	}
	if g.moves == 9 {
		return table.Outcome{Over: true, Winner: -1, Reason: "board full"}
	}
	return table.Outcome{}
}

// Cell reports what's in cell i (0-8): "", "X" or "O".
func (g *Game) Cell(i int) string { return [...]string{"", "X", "O"}[g.cells[i]] }

// LastMove is the cell (0-8) of the last move, or -1.
func (g *Game) LastMove() int {
	if g.moves == 0 {
		return -1
	}
	return g.last
}
