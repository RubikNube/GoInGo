package engine

import (
	"testing"

	"github.com/RubikNube/GoInGo/pkg/game"
)

func TestRandomEngine_MoveReturnsLegalMove(t *testing.T) {
	board := game.Board{}
	engine := &RandomEngine{}
	player := game.Black
	var ko *game.Point

	move := engine.Move(board, player, ko)
	if move == nil {
		t.Error("Expected a move, got nil")
	}
	if move.Row < 0 || move.Row > 8 || move.Col < 0 || move.Col > 8 {
		t.Errorf("Move out of bounds: %+v", move)
	}
}

func TestRandomEngine_MoveReturnsNilWhenNoMoves(t *testing.T) {
	board := game.Board{}
	// Fill the board
	for i := range board {
		for j := range board[i] {
			board[i][j] = game.Black
		}
	}
	engine := &RandomEngine{}
	player := game.White
	var ko *game.Point

	move := engine.Move(board, player, ko)
	if move != nil {
		t.Errorf("Expected nil (pass), got %+v", move)
	}
}
