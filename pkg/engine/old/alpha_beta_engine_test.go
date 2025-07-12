package old

import (
	"testing"

	"github.com/RubikNube/GoInGo/pkg/game"
)

func TestAlphaBetaEngine_MoveReturnsLegalMove(t *testing.T) {
	board := game.Board{}
	engine := &AlphaBetaEngine{}
	player := game.Black
	var ko *game.Point

	move := engine.Move(board, player, ko)
	if move == nil {
		t.Error("Expected a move, got nil")
	}
	if move != nil && (move.Row < 0 || move.Row > 8 || move.Col < 0 || move.Col > 8) {
		t.Errorf("Move out of bounds: %+v", move)
	}
}

func TestAlphaBetaEngine_MoveReturnsNilWhenNoMoves(t *testing.T) {
	board := game.Board{}
	// Fill the board
	for i := range board {
		for j := range board[i] {
			board[i][j] = game.Black
		}
	}
	engine := &AlphaBetaEngine{}
	player := game.White
	var ko *game.Point

	move := engine.Move(board, player, ko)
	if move != nil {
		t.Errorf("Expected nil (pass), got %+v", move)
	}
}

func TestAlphaBetaEngine_PassIsOptimal(t *testing.T) {
	board := game.Board{}
	// Set up a board where any move would be suicide
	for i := range board {
		for j := range board[i] {
			board[i][j] = game.Black
		}
	}
	board[4][4] = game.Empty // Only one empty spot, but surrounded by Black
	engine := &AlphaBetaEngine{}
	player := game.White
	var ko *game.Point

	move := engine.Move(board, player, ko)
	if move != nil {
		t.Errorf("Expected nil (pass) due to suicide, got %+v", move)
	}
}
