package game

import (
	"math/rand/v2"
	"strconv"

	"github.com/JMThomas00/Concord/sdk/arcade"
	"github.com/JMThomas00/Concord/sdk/table"
	"github.com/JMThomas00/concord-tictactoe/engine"
)

// The arcade front door's personality: pencil and paper, the schoolyard
// classic. Gold Stars unlock SETS of pieces and boards; attract mode is
// the computer playing itself in a different set on a different board
// each game; a draw is a cat's game, with the cat.
var arcadeLook = &table.Arcade{
	Title:   "TIC-TAC-TOE",
	Tagline: "THREE IN A ROW",
	HowTo: []string{
		"Take turns placing your piece. X goes first.",
		"Three in a row -- across, down or corner to corner -- wins.",
		"A full board with no three in a row is a draw: a cat's game.",
		"",
		"Win games to earn Gold Stars: one for each new achievement,",
		"every three wins in a row and every ten games. Spend them in",
		"SETS on new pieces and boards. Everyone sees their own.",
	},
	Keys: []arcade.Key{
		{Key: "←↑↓→", Does: "aim"},
		{Key: "1-9", Does: "place"},
		{Key: "Enter", Does: "place"},
	},
	Sounds:      true,
	Reward:      "GOLD STAR",
	Collection:  "SETS",
	Unlockables: unlockables(),
	Kinds:       []table.Kind{{ID: kindPieces, Label: "PIECES"}, {ID: kindBoard, Label: "BOARD"}},
	Preview:     preview,
	Attract:     attract,
	Result:      result,
	ResultArt:   resultArt,
}

func result(g table.Game, o table.Outcome, seats []string) string {
	if o.Winner < 0 {
		return "CAT'S GAME!"
	}
	return seats[o.Winner] + " WINS!"
}

// cat is the cat of a cat's game.
var cat = []string{".k.....k.", ".kk...kk.", ".kkkkkkk.", "kkYkkkYkk", "kkkkPkkkk", ".kk.k.kk.", "..kkkkk..", "...k.k..."}

func resultArt(c *arcade.Canvas, g table.Game, o table.Outcome, x, y, w, h, frame int) bool {
	if o.Winner >= 0 {
		return false
	}
	pal := map[rune]string{'k': "orange", 'Y': "green", 'P': "pink"}
	if frame/4%6 == 5 { // a blink
		pal['Y'] = "orange"
	}
	c.Sprite(x+(w-18)/2, 2*(y+1), cat, pal, 2)
	c.CenterIn(x, w, y+10, "Meow. Nobody wins.", "dim", "", false)
	return true
}

// ── Attract mode ───────────────────────────────────────────────────────────

// The computer plays itself: a move every three ticks, a pause on the
// finished board, then the next game in the next set on the next board.
const (
	attractStep  = 3
	attractPause = 5 // steps to hold a finished board
)

// demoGame is game n's moves, the same every time it's asked for (it's
// drawn frame by frame, so it must not change between frames).
var demoCache struct {
	n     int
	moves []int
	ok    bool
}

func demoGame(n int) []int {
	if demoCache.ok && demoCache.n == n {
		return demoCache.moves
	}
	rnd := rand.New(rand.NewPCG(uint64(n)+1, 0x7ac7ac))
	g := engine.New()
	var moves []int
	for !g.Outcome().Over {
		m := demoMove(g, rnd)
		_ = g.Play(strconv.Itoa(m + 1))
		moves = append(moves, m)
	}
	demoCache.n, demoCache.moves, demoCache.ok = n, moves, true
	return moves
}

// demoMove plays like a decent child: wins when it can, usually blocks,
// otherwise anywhere.
func demoMove(g *engine.Game, rnd *rand.Rand) int {
	mine, theirs := "X", "O"
	if g.Turn() == 1 {
		mine, theirs = "O", "X"
	}
	var free []int
	for i := 0; i < 9; i++ {
		if g.Cell(i) == "" {
			free = append(free, i)
		}
	}
	for _, i := range free {
		if lineFor(g, i, mine) {
			return i
		}
	}
	for _, i := range free {
		if lineFor(g, i, theirs) && rnd.IntN(10) < 8 {
			return i
		}
	}
	return free[rnd.IntN(len(free))]
}

// lineFor reports whether piece at cell i would complete a line.
func lineFor(g *engine.Game, i int, piece string) bool {
	for _, l := range engine.Lines {
		n, has := 0, false
		for _, c := range l {
			if c == i {
				has = true
			} else if g.Cell(c) == piece {
				n++
			}
		}
		if has && n == 2 {
			return true
		}
	}
	return false
}

func attract(c *arcade.Canvas, x, y, w, h, frame int) {
	step := frame / attractStep
	// Find which game this step is in.
	n := 0
	for {
		length := len(demoGame(n)) + attractPause
		if step < length {
			break
		}
		step -= length
		n++
		if n > 1000 { // a still screen after a very long time is fine
			break
		}
	}
	moves := demoGame(n)
	var cells [9]string
	g := engine.New()
	for k := 0; k < min(step+1, len(moves)); k++ {
		_ = g.Play(strconv.Itoa(moves[k] + 1))
	}
	for i := range cells {
		cells[i] = g.Cell(i)
	}
	set := &pieceSets[n%len(pieceSets)]
	board := boardStyles[n%len(boardStyles)]
	bw := w - 22
	drawBoard(c, x, y, bw, h, func(i int) string { return cells[i] }, look{
		style: board.ID, set: set, cursor: -1, last: -1,
		win: g.WinningLine(), lit: true,
	})
	c.Text(x+bw+1, y+3, "NOW SHOWING", "pink", "", true)
	c.Text(x+bw+1, y+5, set.Name, set.Tier.Role(), "", true)
	c.Text(x+bw+1, y+6, "on "+board.Name, "dim", "", false)
	if set.Tier != arcade.Starter {
		c.Text(x+bw+1, y+8, set.Tier.Stars()+" "+set.Tier.Name(), set.Tier.Role(), "", false)
	}
}
