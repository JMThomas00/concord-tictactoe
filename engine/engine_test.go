package engine

import "testing"

func play(t *testing.T, g *Game, moves ...string) {
	t.Helper()
	for _, m := range moves {
		if err := g.Play(m); err != nil {
			t.Fatalf("%s: %v", m, err)
		}
	}
}

func TestRules(t *testing.T) {
	g := New()
	if g.Turn() != 0 || g.LastMove() != -1 {
		t.Fatal("X moves first on an empty board")
	}
	play(t, g, "5", "1")
	if err := g.Play("5"); err == nil {
		t.Fatal("played on a taken cell")
	}
	for _, bad := range []string{"0", "10", "x"} {
		if err := g.Play(bad); err == nil {
			t.Fatalf("accepted %q", bad)
		}
	}
	play(t, g, "3", "2", "7") // X: 5, 3, 7 — a diagonal
	o := g.Outcome()
	if !o.Over || o.Winner != 0 || g.Turn() != -1 {
		t.Fatalf("outcome %+v", o)
	}
	if l := g.WinningLine(); len(l) != 3 || l[0] != 2 || l[1] != 4 || l[2] != 6 {
		t.Fatalf("winning line %v", l)
	}
	if g.LastMove() != 6 || g.Cell(6) != "X" || g.Cell(0) != "O" {
		t.Fatal("cells or last move wrong")
	}
	if err := g.Play("9"); err == nil {
		t.Fatal("played after the game ended")
	}
}

func TestDraw(t *testing.T) {
	g := New()
	play(t, g, "1", "2", "3", "5", "4", "6", "8", "7", "9")
	if o := g.Outcome(); !o.Over || o.Winner != -1 || g.WinningLine() != nil {
		t.Fatalf("outcome %+v", o)
	}
}

// The perfect computer never loses, whoever moves first, against every
// line of play.
func TestPerfectComputerNeverLoses(t *testing.T) {
	var explore func(g *Game, computer int)
	explore = func(g *Game, computer int) {
		if o := g.Outcome(); o.Over {
			if o.Winner >= 0 && o.Winner != computer {
				t.Fatalf("the computer (%d) lost: %v", computer, g.cells)
			}
			return
		}
		if g.Turn() == computer {
			c := *g
			if err := c.Play(BestMove(&c, 3)); err != nil {
				t.Fatal(err)
			}
			explore(&c, computer)
			return
		}
		for i := 0; i < 9; i++ {
			if g.cells[i] == 0 {
				c := *g
				_ = c.Play(string(rune('1' + i)))
				explore(&c, computer)
			}
		}
	}
	explore(New(), 0)
	explore(New(), 1)
}

func TestMediumComputerWinsAndBlocks(t *testing.T) {
	g := New()
	play(t, g, "1", "4", "2") // X threatens 3; O to move
	if m := BestMove(g, 2); m != "3" {
		t.Fatalf("didn't block: %s", m)
	}
	g = New()
	play(t, g, "1", "4", "9", "5") // O threatens 6, X to move but X can't win: block
	if m := BestMove(g, 2); m != "6" {
		t.Fatalf("didn't block: %s", m)
	}
}
