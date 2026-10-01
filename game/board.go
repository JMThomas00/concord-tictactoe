package game

import (
	"math"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/JMThomas00/Concord/sdk/table"
	"github.com/JMThomas00/Concord/sdk/wire"
	"github.com/JMThomas00/concord-tictactoe/engine"
)

// Board is what a viewer sees and plays on. With room it draws a large
// framed grid: in Concord the pieces are pictures in the viewer's theme
// colors (Images), and everywhere else large box-drawn X and O glyphs. In a
// small pane it falls back to a compact grid. ? shows the rules and keys.
type Board struct {
	seat   *table.Seat
	cursor int
	err    string
	help   bool
	width  int
	height int
}

func (b *Board) Init() tea.Cmd { return nil }

func (b *Board) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		b.width, b.height = msg.Width, msg.Height
	case table.ChangedMsg:
		b.err = ""
	case tea.KeyMsg:
		k := msg.String()
		if b.help {
			if k == "?" || k == "esc" || k == "q" || k == "enter" {
				b.help = false
			}
			return b, nil
		}
		switch k {
		case "?":
			b.help = true
		case "up", "k":
			b.cursor = (b.cursor + 6) % 9
		case "down", "j":
			b.cursor = (b.cursor + 3) % 9
		case "left", "h":
			b.cursor = b.cursor/3*3 + (b.cursor+2)%3
		case "right", "l":
			b.cursor = b.cursor/3*3 + (b.cursor+1)%3
		case "enter", " ":
			b.move(strconv.Itoa(b.cursor + 1))
		case "q", "esc":
			return b, tea.Quit // hand the keyboard back to Concord
		default:
			if len(k) == 1 && k[0] >= '1' && k[0] <= '9' {
				b.cursor = int(k[0] - '1')
				b.move(k)
			}
		}
	}
	return b, nil
}

func (b *Board) move(cell string) {
	if !b.seat.MyTurn() {
		if b.seat.Outcome().Over {
			b.err = "the game is over — Tab for a rematch"
		} else {
			b.err = "not your turn"
		}
		return
	}
	if err := b.seat.Play(cell); err != nil {
		b.err = err.Error()
	}
}

func (b *Board) game() *engine.Game { return b.seat.Game().(*engine.Game) }

// ── Colors ─────────────────────────────────────────────────────────────────

// color is a theme color by name, with Dracula's as the fallback.
func (b *Board) color(name string) lipgloss.TerminalColor {
	fallback := map[string]string{
		"red": "#FF5555", "cyan": "#8BE9FD", "comment": "#6272A4", "selection": "#44475A",
		"current_line": "#6272A4", "yellow": "#F1FA8C", "green": "#50FA7B", "foreground": "#F8F8F2", "purple": "#BD93F9",
	}[name]
	return b.seat.Color(name, lipgloss.Color(fallback))
}

func (b *Board) style(name string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(b.color(name))
}

func (b *Board) pieceStyle(piece string) lipgloss.Style {
	if piece == "X" {
		return b.style("red").Bold(true)
	}
	return b.style("cyan").Bold(true)
}

// ── Layout ─────────────────────────────────────────────────────────────────

// statusRows are the lines under the grid.
const statusRows = 3

// cellSize is a large cell's size in cells (inside its frame lines), or
// 0, 0 when the pane is too small and the compact grid is used.
func (b *Board) cellSize() (w, h int) {
	h = min(6, (b.height-4-statusRows)/3)
	w = h*2 + 1
	if h < 3 || 3*w+4 > b.width {
		return 0, 0
	}
	return w, h
}

// ── View ───────────────────────────────────────────────────────────────────

func (b *Board) View() string {
	if b.help {
		return b.helpView()
	}
	var grid []string
	if w, h := b.cellSize(); w > 0 {
		grid = b.largeGrid(w, h)
	} else {
		grid = b.compactGrid()
	}
	return strings.Join(append(grid, b.status()...), "\n")
}

// cellBackground is how a cell is highlighted: the cursor, or part of the
// winning line.
func (b *Board) cellBackground(i int) (lipgloss.TerminalColor, bool) {
	g := b.game()
	for _, c := range g.WinningLine() {
		if c == i {
			return b.color("current_line"), true
		}
	}
	if i == b.cursor && b.seat.MyTurn() {
		return b.color("selection"), true
	}
	return nil, false
}

func (b *Board) largeGrid(w, h int) []string {
	g := b.game()
	frame := b.style("comment")
	seg := strings.Repeat("─", w)
	rule := func(l, m, r string) string { return frame.Render(l + seg + m + seg + m + seg + r) }
	rows := []string{rule("╭", "┬", "╮")}
	for r := 0; r < 3; r++ {
		for line := 0; line < h; line++ {
			var b2 strings.Builder
			b2.WriteString(frame.Render("│"))
			for c := 0; c < 3; c++ {
				i := r*3 + c
				b2.WriteString(b.cellLine(g, i, w, h, line))
				b2.WriteString(frame.Render("│"))
			}
			rows = append(rows, b2.String())
		}
		if r < 2 {
			rows = append(rows, rule("├", "┼", "┤"))
		}
	}
	return append(rows, rule("╰", "┴", "╯"))
}

// cellLine is one line of a large cell: a margin, the piece's glyph (which
// Concord covers with its picture), a margin, all on the cell's highlight
// if it has one. The last move is marked in the left margin.
func (b *Board) cellLine(g *engine.Game, i, w, h, line int) string {
	base := lipgloss.NewStyle()
	if bg, ok := b.cellBackground(i); ok {
		base = base.Background(bg)
	}
	left := base.Render(" ")
	if g.LastMove() == i && line == h/2 {
		left = base.Foreground(b.color("yellow")).Render("•")
	}
	glyph := base.Render(strings.Repeat(" ", w-2))
	if p := g.Cell(i); p != "" {
		glyph = base.Inherit(b.pieceStyle(p)).Render(glyphLine(p, w-2, h, line))
	}
	return left + glyph + base.Render(" ")
}

// glyphLine draws line `line` of a large X or O glyph inner cells wide and
// h tall. The shape is worked out on a grid of half-cells (two per
// character, top and bottom, which are about square) and drawn with block
// characters, so it looks the same in every font.
func glyphLine(piece string, inner, h, line int) string {
	w, hh := float64(inner), float64(2*h)
	on := func(x, y int) bool {
		px, py := float64(x)+0.5, float64(y)+0.5
		switch piece {
		case "X":
			// Near either diagonal of the box.
			n := math.Hypot(w, hh)
			return math.Abs(hh*px-w*py)/n < 0.8 || math.Abs(hh*(w-px)-w*py)/n < 0.8
		case "O":
			// A ring about 1.5 half-cells thick.
			r := math.Min(w, hh)/2 - 0.1
			d := math.Hypot(px-w/2, (py-hh/2)*w/hh)
			return d <= r && d >= r-1.5
		}
		return false
	}
	var sb strings.Builder
	for x := 0; x < inner; x++ {
		top, bottom := on(x, 2*line), on(x, 2*line+1)
		switch {
		case top && bottom:
			sb.WriteRune('█')
		case top:
			sb.WriteRune('▀')
		case bottom:
			sb.WriteRune('▄')
		default:
			sb.WriteRune(' ')
		}
	}
	return sb.String()
}

func (b *Board) compactGrid() []string {
	g := b.game()
	frame := b.style("comment")
	var rows []string
	for r := 0; r < 3; r++ {
		var cells []string
		for c := 0; c < 3; c++ {
			i := r*3 + c
			mark := " "
			if p := g.Cell(i); p != "" {
				mark = b.pieceStyle(p).Render(p)
			}
			cell := " " + mark + " "
			if bg, ok := b.cellBackground(i); ok {
				st := lipgloss.NewStyle().Background(bg)
				if p := g.Cell(i); p != "" {
					cell = st.Render(" ") + st.Inherit(b.pieceStyle(p)).Render(p) + st.Render(" ")
				} else {
					cell = st.Render("   ")
				}
			}
			cells = append(cells, cell)
		}
		rows = append(rows, strings.Join(cells, frame.Render("│")))
		if r < 2 {
			rows = append(rows, frame.Render("───┼───┼───"))
		}
	}
	return rows
}

// status is the lines under the grid: what's happening, and the help hint.
func (b *Board) status() []string {
	g := b.game()
	var line string
	switch o := b.seat.Outcome(); {
	case b.err != "":
		line = b.style("red").Render(b.err)
	case o.Over && o.Winner >= 0 && o.Winner == b.seat.Index:
		line = b.style("green").Bold(true).Render("You win!")
	case o.Over && o.Winner >= 0 && b.seat.Index >= 0:
		line = b.style("comment").Render("You lose this one.")
	case o.Over && o.Winner < 0:
		line = b.style("yellow").Render("A draw — nobody can win.")
	case o.Over:
		line = ""
	case b.seat.MyTurn():
		me := b.seat.SeatName(b.seat.Index)
		line = b.pieceStyle(me).Render(me) + " " + b.style("foreground").Render("Your move") +
			b.style("comment").Render(" · arrows + Enter, or 1-9")
	case b.seat.Index >= 0:
		turn := g.Turn()
		if turn >= 0 {
			name := b.seat.SeatName(turn)
			line = b.style("comment").Render("Waiting for ") + b.pieceStyle(name).Render(name) + b.style("comment").Render("…")
		}
	default:
		line = b.style("comment").Render("Watching")
	}
	return []string{"", line, b.style("comment").Render("? rules and keys")}
}

func (b *Board) helpView() string {
	h := b.style("purple").Bold(true)
	k := b.style("cyan")
	t := b.style("foreground")
	d := b.style("comment")
	lines := []string{
		h.Render("Tic-tac-toe"),
		"",
		t.Render("Take turns placing your mark; X goes first."),
		t.Render("Three in a row — across, down or diagonal — wins."),
		t.Render("A full board with no three in a row is a draw."),
		"",
		h.Render("Keys"),
	}
	or := d.Render(" or ")
	for _, row := range [][2]string{
		{k.Render("← → ↑ ↓") + or + k.Render("h j k l"), "move the cursor"},
		{k.Render("Enter") + or + k.Render("Space"), "place your mark"},
		{k.Render("1") + d.Render("–") + k.Render("9"), "place it directly (1 2 3 / 4 5 6 / 7 8 9)"},
		{k.Render("Tab"), "sit, stand, rematch, play the computer"},
		{k.Render("?"), "this help"},
		{k.Render("q") + or + k.Render("Esc"), "back to Concord"},
	} {
		pad := max(1, 22-lipgloss.Width(row[0]))
		lines = append(lines, row[0]+strings.Repeat(" ", pad)+t.Render(row[1]))
	}
	lines = append(lines, "", d.Render("Press ? or Esc to close."))
	return strings.Join(lines, "\n")
}

// Images puts piece pictures in the viewer's theme colors over the large
// grid's cells (Concord draws them however each terminal can; the glyphs
// underneath show where it can't).
func (b *Board) Images() []wire.PaneImage {
	w, h := b.cellSize()
	if w == 0 || b.help {
		return nil
	}
	g := b.game()
	var out []wire.PaneImage
	for i := 0; i < 9; i++ {
		p := g.Cell(i)
		if p == "" {
			continue
		}
		c := b.color("red")
		if p == "O" {
			c = b.color("cyan")
		}
		r, col := i/3, i%3
		out = append(out, wire.PaneImage{
			Asset: pieceAsset(p, c),
			Col:   1 + col*(w+1) + 1, // inside the frame and the cell's margin
			Row:   1 + r*(h+1),
			Cols:  w - 2,
			Rows:  h,
		})
	}
	return out
}
