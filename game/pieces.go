package game

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Piece pictures come in every red (X) and cyan (O) that Concord's themes
// use, so a board matches the viewer's theme. A custom theme gets the
// nearest color drawn.

// ansiStandIns are rough colors for the terminal-default theme, whose
// palette is ANSI color numbers rather than colors.
var ansiStandIns = map[string]string{"1": "#CD3131", "9": "#CD3131", "6": "#11A8CD", "14": "#11A8CD"}

// pieceAsset is the picture for piece ("X" or "O") in the theme color c.
func pieceAsset(piece string, c lipgloss.TerminalColor) string {
	choices, prefix := xColors, "x"
	if piece == "O" {
		choices, prefix = oColors, "o"
	}
	want := ""
	if lc, ok := c.(lipgloss.Color); ok {
		want = string(lc)
	}
	if s, ok := ansiStandIns[want]; ok {
		want = s
	}
	best := choices[0]
	if r, g, b, ok := hexRGB(want); ok {
		bestD := -1
		for _, ch := range choices {
			cr, cg, cb, _ := hexRGB(ch)
			if d := sq(r-cr) + sq(g-cg) + sq(b-cb); bestD < 0 || d < bestD {
				best, bestD = ch, d
			}
		}
	}
	return "assets/" + prefix + "-" + strings.ToLower(strings.TrimPrefix(best, "#")) + ".png"
}

func hexRGB(s string) (int, int, int, bool) {
	var r, g, b int
	if len(s) != 7 || s[0] != '#' {
		return 0, 0, 0, false
	}
	if _, err := fmt.Sscanf(strings.ToLower(s), "#%02x%02x%02x", &r, &g, &b); err != nil {
		return 0, 0, 0, false
	}
	return r, g, b, true
}

func sq(n int) int { return n * n }
