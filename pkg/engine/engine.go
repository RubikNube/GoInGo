// Package engine provides a simple Go engine interface and a random-move engine.
package engine

import (
	"github.com/RubikNube/GoInGo/pkg/game"
)

// Engine is an interface for Go engines.
type Engine interface {
	// Move returns the next move as a Point, or nil if passing.
	Move(board game.Board, player game.FieldState, ko *game.Point) *game.Point
}
