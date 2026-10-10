package game

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/JMThomas00/Concord/sdk/arcade"
	"github.com/JMThomas00/Concord/sdk/plugin"
	"github.com/JMThomas00/Concord/sdk/plugintest"
	"github.com/JMThomas00/Concord/sdk/table"
	"github.com/JMThomas00/Concord/sdk/wire"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

// Every piece the board can ask for exists in client/assets.
func TestPieceAssetsExistForEveryThemeColor(t *testing.T) {
	for _, set := range []struct {
		piece  string
		colors []string
	}{{"X", xColors}, {"O", oColors}} {
		for _, c := range set.colors {
			asset := pieceAsset(set.piece, lipgloss.Color(c))
			if !strings.Contains(asset, strings.ToLower(c[1:])) {
				t.Fatalf("%s in %s picked %s", set.piece, c, asset)
			}
			if _, err := os.Stat(filepath.Join("..", "client", asset)); err != nil {
				t.Fatalf("missing %s", asset)
			}
		}
	}
	// A custom theme gets the nearest drawn color; terminal-default's ANSI
	// numbers get stand-ins.
	if got := pieceAsset("X", lipgloss.Color("#FE5656")); got != "assets/x-ff5555.png" {
		t.Fatalf("nearest red: %s", got)
	}
	if got := pieceAsset("O", lipgloss.Color("6")); got != "assets/o-11a8cd.png" {
		t.Fatalf("ANSI cyan: %s", got)
	}
}

func start(t *testing.T) (*plugintest.Server, uuid.UUID) {
	t.Helper()
	srv := plugintest.NewServer(t)
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	go plugin.Run(ctx, srv.Config(), table.New(Rules).Handler())
	srv.WaitReady()
	ch := uuid.New()
	srv.Channel(wire.Channel{ID: ch, Name: "tictactoe"})
	return srv, ch
}

// toTable takes a viewer through the front door to the seats table: the
// title, then TAKE A SEAT (which sits them down if there's a seat).
func toTable(srv *plugintest.Server, v *plugintest.Viewer) {
	srv.FrameContaining(v, "PRESS ENTER")
	srv.Key(v, "enter")
	srv.FrameContaining(v, "TAKE A SEAT")
	srv.Key(v, "down")
	srv.Key(v, "enter")
}

// Two players come in through the arcade's front door, sit down and play
// to a win: pieces drawn as pixels in each viewer's set, a sound per move,
// a Gold Star for the first win, and a rematch once both ask for one.
func TestAGameInAChannel(t *testing.T) {
	srv, ch := start(t)
	x := srv.Enter(ch, "alice", 80, 30)
	o := srv.Enter(ch, "bob", 80, 30)
	if c := srv.Claimed(x); len(c) != 0 {
		t.Fatalf("the title screen claims %v; Esc should leave the pane", c)
	}
	toTable(srv, x)
	toTable(srv, o)
	srv.FrameContaining(x, "YOUR MOVE")
	if c := srv.Claimed(x); len(c) != 1 || c[0] != wire.PaneKeyEsc {
		t.Fatalf("the table claims %v; Esc should go back to the menu", c)
	}

	srv.Key(x, "5")
	frame := srv.FrameContaining(o, "YOUR MOVE")
	if !strings.ContainsAny(frame, "█▀▄") {
		t.Fatalf("no pixel piece on the board:\n%s", frame)
	}
	if len(srv.Images(o)) != 0 {
		t.Fatal("the arcade draws pieces itself, with no pictures")
	}
	for i := 0; ; i++ { // menu sounds come first
		if s, _ := srv.NextSound(); s.Asset == "sounds/move.wav" {
			break
		}
		if i > 10 {
			t.Fatal("no move sound")
		}
	}

	for i, mv := range []string{"1", "4", "2", "6"} { // X: 5, 4, 6 -- the middle row
		v := o
		if i%2 == 1 {
			v = x
		}
		srv.Key(v, mv)
	}
	srv.FrameContaining(x, "YOU WIN! ENTER: RESULTS")
	srv.FrameContaining(o, "YOU LOSE THIS ONE.")

	srv.Key(x, "enter")
	frame = srv.FrameContaining(x, "REMATCH?")
	if !strings.Contains(frame, "+1 GOLD STAR") || !strings.Contains(frame, "WON IN 5 MOVES") {
		t.Fatalf("the winner's results:\n%s", frame)
	}
	srv.Key(o, "enter")
	frame = srv.FrameContaining(o, "REMATCH?")
	if strings.Contains(frame, "GOLD STAR") {
		t.Fatalf("a Gold Star for losing the first game:\n%s", frame)
	}
	srv.Key(x, "enter")
	srv.Key(o, "enter")
	srv.FrameContaining(o, "YOUR MOVE") // seats swapped: bob is X now
}

// A draw is a cat's game.
func TestCatsGame(t *testing.T) {
	srv, ch := start(t)
	x := srv.Enter(ch, "alice", 80, 30)
	o := srv.Enter(ch, "bob", 80, 30)
	toTable(srv, x)
	toTable(srv, o)
	srv.FrameContaining(x, "YOUR MOVE")
	for i, mv := range []string{"1", "2", "3", "5", "4", "6", "8", "7", "9"} {
		v := x
		if i%2 == 1 {
			v = o
		}
		srv.Key(v, mv)
	}
	srv.FrameContaining(x, "A DRAW. ENTER: RESULTS")
	srv.Key(x, "enter")
	srv.FrameContaining(x, "Meow. Nobody wins.")
}

// SETS shows the player's pieces and board, and what's still locked.
func TestSets(t *testing.T) {
	srv, ch := start(t)
	v := srv.Enter(ch, "alice", 80, 24)
	srv.FrameContaining(v, "PRESS ENTER")
	srv.Key(v, "enter")
	frame := srv.FrameContaining(v, "SETS")
	if !strings.Contains(frame, "YOUR PIECES") {
		t.Fatalf("the menu doesn't show the player's pieces:\n%s", frame)
	}
	for _, k := range []string{"down", "down", "enter"} { // 1P, TAKE A SEAT, SETS
		srv.Key(v, k)
	}
	frame = srv.FrameContaining(v, "✓ IN USE")
	for _, want := range []string{"PIECES", "BOARD", "◂ CLASSIC ▸", "GOLD STARS"} {
		if !strings.Contains(frame, want) {
			t.Fatalf("SETS is missing %q:\n%s", want, frame)
		}
	}
	srv.Key(v, "right")
	srv.Key(v, "right") // grapes & leaves: locked
	frame = srv.FrameContaining(v, "? ? ?")
	if !strings.Contains(frame, "LOCKED: EARN GOLD STARS") {
		t.Fatalf("a locked set:\n%s", frame)
	}
	srv.Key(v, "enter")
	srv.FrameContaining(v, "EARN GOLD STARS TO UNLOCK MORE.")
	srv.Key(v, "esc")
	srv.FrameContaining(v, "HALL OF FAME")
}

// A pane too small for the arcade gets the plain front door, then the
// compact grid.
func TestHelpAndTheCompactBoard(t *testing.T) {
	srv, ch := start(t)
	small := srv.Enter(ch, "carol", 40, 12)
	srv.FrameContaining(small, "Enter to play")
	srv.Key(small, "enter")
	frame := srv.FrameContaining(small, "M: sit down")
	if strings.Contains(frame, "╭") {
		t.Fatalf("a small pane should get the compact grid:\n%s", frame)
	}
	if len(srv.Images(small)) != 0 {
		t.Fatal("pictures on the compact grid")
	}
	srv.Key(small, "?")
	srv.FrameContaining(small, "Three in a row")
	if c := srv.Claimed(small); len(c) != 1 || c[0] != wire.PaneKeyEsc {
		t.Fatalf("help open claims %v", c)
	}
	srv.Key(small, "esc")
	srv.FrameContaining(small, "? rules and keys")
	if len(srv.Claimed(small)) != 0 {
		t.Fatal("Esc still claimed after closing help")
	}
}

// Every set and board draws at every size it's shown at, and attract mode
// draws the same picture for the same frame.
func TestDrawing(t *testing.T) {
	sizes := [][2]int{{42, 18}, {40, 12}, {21, 6}, {20, 6}, {38, 17}, {54, 14}}
	for _, s := range pieceSets {
		for _, p := range []string{s.xMini, s.oMini} {
			if arcade.TextWidth(p) != 1 {
				t.Errorf("%s: mini piece %q is %d cells wide", s.ID, p, arcade.TextWidth(p))
			}
		}
		for _, rows := range [][]string{s.x, s.o} {
			if len(rows) != 8 || arcade.SpriteWidth(rows, 1) != 9 {
				t.Errorf("%s: a piece isn't 9 x 8", s.ID)
			}
		}
		for _, b := range boardStyles {
			for _, size := range sizes {
				c := arcade.New(80, 24, arcade.NewPalette(nil))
				drawBoard(c, 2, 2, size[0], size[1], func(i int) string { return sampleCells[i] },
					look{style: b.ID, set: &s, cursor: 3, win: []int{0, 4, 8}, lit: true, last: 8})
				if out := c.String(); len(strings.Split(out, "\n")) != 24 {
					t.Fatalf("%s on %s at %v: %d lines", s.ID, b.ID, size, len(strings.Split(out, "\n")))
				}
			}
		}
	}
	draw := func(frame int) string {
		c := arcade.New(80, 24, arcade.NewPalette(nil))
		attract(c, 2, 6, 76, 14, frame)
		return c.String()
	}
	for f := 0; f < 600; f += 7 {
		if draw(f) != draw(f) {
			t.Fatalf("attract mode changes within frame %d", f)
		}
	}
	if draw(0) == draw(300) {
		t.Fatal("attract mode doesn't move")
	}
}

func TestSounds(t *testing.T) {
	g := Rules.New(nil)
	for _, m := range []string{"1", "2", "3", "5", "4", "6", "8", "7"} {
		_ = g.Play(m)
		if s := sound(g, m); s != "sounds/move.wav" {
			t.Fatalf("%s: %s", m, s)
		}
	}
	_ = g.Play("9")
	if s := sound(g, "9"); s != "sounds/draw.wav" {
		t.Fatalf("draw: %s", s)
	}
}
