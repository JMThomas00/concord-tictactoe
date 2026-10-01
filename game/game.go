// Package game hosts tic-tac-toe on Concord's table kit: the rules hookup,
// the board people see, its pictures and its sounds.
package game

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/JMThomas00/Concord/sdk/table"
	"github.com/JMThomas00/concord-tictactoe/engine"
)

// Rules is everything the table kit needs to host tic-tac-toe.
var Rules = table.Rules{
	Name:      "Tic-tac-toe",
	SeatNames: []string{"X", "O"},
	New:       func(map[string]string) table.Game { return engine.New() },
	NewBoard:  func(s *table.Seat) tea.Model { return &Board{seat: s, cursor: 4} },
	AI:        func(g table.Game, level int) string { return engine.BestMove(g.(*engine.Game), level) },
	Sound:     sound,
}

// sound is what everyone watching hears after a move.
func sound(g table.Game, move string) string {
	switch o := g.Outcome(); {
	case o.Over && o.Winner >= 0:
		return "sounds/win.wav"
	case o.Over:
		return "sounds/draw.wav"
	}
	return "sounds/move.wav"
}
