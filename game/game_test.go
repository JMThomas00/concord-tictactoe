package game

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

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

// Two players sit down and play to a win: pictures on the large board,
// a sound per move and the win chime at the end.
func TestAGameInAChannel(t *testing.T) {
	srv, ch := start(t)
	x := srv.Enter(ch, "alice", 80, 30)
	o := srv.Enter(ch, "bob", 80, 30)
	for _, v := range []*plugintest.Viewer{x, o} {
		srv.FrameContaining(v, "Tab: sit down")
		srv.Key(v, "tab")
		srv.Key(v, "enter")
	}
	srv.FrameContaining(x, "Your move")

	srv.Key(x, "5")
	frame := srv.FrameContaining(o, "Your move")
	if !strings.Contains(frame, "╭") || !strings.Contains(frame, "█") {
		t.Fatalf("no large board with the X glyph:\n%s", frame)
	}
	images := srv.Images(o)
	if len(images) != 1 || images[0].Asset != "assets/x-ff5555.png" || images[0].Cols <= 0 {
		t.Fatalf("images %+v", images)
	}
	if s, _ := srv.NextSound(); s.Asset != "sounds/move.wav" || s.ChannelID != ch {
		t.Fatalf("sound %+v", s)
	}

	for i, mv := range []string{"1", "4", "2", "6"} { // X: 5, 4, 6 — the middle row
		v := o
		if i%2 == 1 {
			v = x
		}
		srv.Key(v, mv)
		srv.NextSound()
	}
	srv.FrameContaining(x, "You win!")
	srv.FrameContaining(o, "You lose this one.")
	// The win chime is the last sound.
	srv.DrainEvents()
	if got := len(srv.Images(x)); got != 5 {
		t.Fatalf("%d pictures on the final board, want 5", got)
	}
}

func TestHelpAndTheCompactBoard(t *testing.T) {
	srv, ch := start(t)
	small := srv.Enter(ch, "carol", 40, 12)
	frame := srv.FrameContaining(small, "Tab: sit down")
	if strings.Contains(frame, "╭") {
		t.Fatalf("a small pane should get the compact grid:\n%s", frame)
	}
	if len(srv.Images(small)) != 0 {
		t.Fatal("pictures on the compact grid")
	}
	srv.Key(small, "?")
	srv.FrameContaining(small, "Three in a row")
	srv.Key(small, "esc")
	srv.FrameContaining(small, "? rules and keys")
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
