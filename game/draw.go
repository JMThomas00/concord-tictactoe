package game

import (
	"strings"

	"github.com/JMThomas00/Concord/sdk/arcade"
)

// Drawing on the arcade canvas: a 3 x 3 board in one of the board styles,
// with pieces from one set, at whatever size fits.

// look is how to draw one board.
type look struct {
	style  string // a board style id (without boardPrefix)
	set    *pieceSet
	cursor int   // the cell to mark for placing, or -1
	win    []int // the winning line, highlighted when lit
	lit    bool
	last   int  // the last move's cell, marked, or -1
	ghost  bool // a locked silhouette
	bare   bool // no frame around the grid (set by drawBoard when space is short)
}

// cellsFor picks a cell size for a board in w x h cells: pictures need 11 x
// 5 (a 9 x 8 pixel piece with a margin); smaller boards use characters.
func cellsFor(w, h int) (cw, ch int, big bool) {
	ch = min(5, (h-2)/3)
	cw = min(11, (w-2)/3)
	if ch >= 4 && cw >= 10 {
		return cw, ch, true
	}
	ch = max(1, min(3, ch))
	cw = max(3, min(2*ch+1, cw))
	return cw, ch, false
}

// boardSize is the size of a board with cw x ch cells.
func boardSize(cw, ch int) (int, int) { return 3*cw + 2, 3*ch + 2 }

// frameOf is how far a board style draws outside its grid: columns each
// side, rows above and below.
func frameOf(style string) (mx, my int) {
	switch style {
	case "chalk", "wood", "vineyard":
		return 1, 1
	case "neon", "picnic":
		return 2, 1
	case "notebook":
		return 4, 0
	}
	return 0, 0
}

// fit picks the cell size for a board in w x h cells, keeping its frame
// inside too. When pictures fit only without the frame, the frame is left
// off (bare): the pieces matter more.
func fit(style string, w, h int) (cw, ch int, big, bare bool) {
	mx, my := frameOf(style)
	if cw, ch, big = cellsFor(w-2*mx, h-2*my); big {
		return cw, ch, true, false
	}
	if cw2, ch2, big2 := cellsFor(w, h); big2 {
		return cw2, ch2, true, true
	}
	bw, bh := boardSize(cw, ch)
	return cw, ch, false, bw+2*mx > w || bh+2*my > h
}

// drawBoard draws a board centred in w x h cells at (x, y). cell gives each
// square's piece ("", "X" or "O").
func drawBoard(c *arcade.Canvas, x, y, w, h int, cell func(i int) string, l look) {
	cw, ch, big, bare := fit(l.style, w, h)
	l.bare = bare
	bw, bh := boardSize(cw, ch)
	x += max(0, (w-bw)/2)
	y += max(0, (h-bh)/2)
	drawGrid(c, x, y, cw, ch, l)
	for i := 0; i < 9; i++ {
		cx, cy := x+(i%3)*(cw+1), y+(i/3)*(ch+1)
		for _, k := range l.win {
			if k == i && l.lit && !l.ghost {
				c.Shade(cx, cy, cw, ch, "yellowB")
			}
		}
		p := cell(i)
		if p == "" {
			if i == l.cursor {
				corner := func(dx, dy int, s string) { c.Text(cx+dx, cy+dy, s, "yellow", "", true) }
				corner(0, 0, "┌")
				corner(cw-1, 0, "┐")
				corner(0, ch-1, "└")
				corner(cw-1, ch-1, "┘")
				if ch == 1 {
					c.Text(cx, cy, "[", "yellow", "", true)
					c.Text(cx+cw-1, cy, "]", "yellow", "", true)
				}
			}
			continue
		}
		if big {
			rows, pal := l.set.sprite(p, l.ghost)
			px := cx + (cw-9)/2
			py := 2*cy + ch - 4 // centred: the picture is 4 rows tall
			c.Sprite(px, py, rows, pal, 1)
		} else {
			s, role := l.set.mini(p, l.ghost)
			c.Text(cx+cw/2, cy+ch/2, s, role, "", true)
		}
		if i == l.last && !l.ghost && big { // small cells have no room for it
			c.Text(cx, cy+ch/2, "•", "yellow", "", true)
		}
	}
}

// drawGrid draws the board's lines (and its surface) in style.
func drawGrid(c *arcade.Canvas, x, y, cw, ch int, l look) {
	bw, bh := boardSize(cw, ch)
	vx := []int{x + cw, x + 2*cw + 1}
	hy := []int{y + ch, y + 2*ch + 1}
	g := func(role string) string {
		if l.ghost {
			return "ghost"
		}
		return role
	}
	bg := func(role string) string {
		if l.ghost {
			return ""
		}
		return role
	}
	lines := func(role, v, hz, cross, back string) {
		for j := 0; j < bh; j++ {
			for _, X := range vx {
				c.Text(X, y+j, v, role, back, false)
			}
		}
		for _, Y := range hy {
			c.Fill(x, Y, bw, hz, role, back)
			for _, X := range vx {
				c.Text(X, Y, cross, role, back, false)
			}
		}
	}
	e := 1 // how far the frame reaches past the grid
	if l.bare {
		e = 0
	}
	switch l.style {
	case "chalk":
		c.Shade(x-e, y-e, bw+2*e, bh+2*e, bg("greenB"))
		lines(g("dim"), "┃", "━", "╋", bg("greenB"))
	case "notebook":
		for j := 1; j < bh; j += 2 {
			c.Fill(x-2*e, y+j, bw+4*e, "┈", "ghost", "")
		}
		for j := 0; j < bh && !l.bare; j++ {
			c.Text(x-4, y+j, "│", g("red"), "", false)
		}
		lines(g("fg"), "│", "─", "┼", "")
	case "wood":
		c.Shade(x-e, y-e, bw+2*e, bh+2*e, bg("orangeB"))
		lines(g("orangeD"), "█", "█", "█", bg("orangeB"))
	case "neon":
		lines(g("pink"), "║", "═", "╬", "")
		if !l.bare {
			c.Box(x-2, y-1, bw+4, bh+2, g("purple"), "", "")
		}
	case "picnic":
		for i := -2; i < bw+2 && !l.bare; i++ {
			top, bot := "▄", "▀"
			if i%2 != 0 {
				top, bot = "▀", "▄"
			}
			c.Text(x+i, y-1, top, g("red"), "", false)
			c.Text(x+i, y+bh, bot, g("red"), "", false)
		}
		for j := 0; j < bh && !l.bare; j++ {
			a, b := "▌", "▐"
			if j%2 != 0 {
				a, b = "▐", "▌"
			}
			c.Text(x-2, y+j, a, g("red"), "", false)
			c.Text(x+bw+1, y+j, b, g("red"), "", false)
		}
		lines(g("red"), "┆", "┄", "┼", "")
	case "vineyard":
		lines(g("greenD"), "│", "─", "┼", "")
		for _, Y := range hy {
			for i := 2; i < bw; i += 6 {
				if !contains(vx, x+i) {
					c.Text(x+i, Y, "~", g("green"), "", false)
				}
			}
			for _, X := range vx {
				c.Text(X, Y, "●", g("purple"), "", true)
			}
		}
		for _, X := range vx {
			if l.bare {
				break
			}
			c.Text(X, y-1, "❦", g("green"), "", false)
		}
	default:
		lines(g("comment"), "│", "─", "┼", "")
	}
}

func contains(xs []int, v int) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}

// drawPair draws a set's X and O side by side, centred in w x h.
func drawPair(c *arcade.Canvas, s *pieceSet, x, y, w, h int, ghost bool) {
	if w >= 20 && h >= 4 {
		x += (w - 20) / 2
		y += (h - 4) / 2
		for i, p := range []string{"X", "O"} {
			rows, pal := s.sprite(p, ghost)
			c.Sprite(x+i*11, 2*y, rows, pal, 1)
		}
		return
	}
	xs, xr := s.mini("X", ghost)
	os, or := s.mini("O", ghost)
	at := x + (w-3)/2
	c.Text(at, y+h/2, xs, xr, "", true)
	c.Text(at+2, y+h/2, os, or, "", true)
}

// sampleCells is a game in progress, for previews.
var sampleCells = [9]string{"X", "O", "X", "", "O", "", "O", "", "X"}

// preview draws an unlockable for the arcade: a piece set as its pair on
// small cards, or in a sample game on the player's board when there's
// room; a board as a sample game in the player's pieces.
func preview(c *arcade.Canvas, id string, sel map[string]string, x, y, w, h int, ghost bool) {
	set := pieces(sel[kindPieces])
	style := strings.TrimPrefix(sel[kindBoard], boardPrefix)
	if !strings.HasPrefix(id, boardPrefix) {
		set = pieces(id)
		if h < 10 {
			drawPair(c, set, x, y, w, h, ghost)
			return
		}
	}
	drawBoard(c, x, y, w, h, func(i int) string { return sampleCells[i] }, look{
		style: style, set: set, cursor: -1, last: -1, ghost: ghost,
	})
}
