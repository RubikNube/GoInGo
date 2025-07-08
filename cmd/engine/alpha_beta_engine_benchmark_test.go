package engine

import (
	"testing"

	"github.com/RubikNube/GoInGo/cmd/game"
)

func emptyBoard() game.Board {
	var b game.Board
	for i := range b {
		for j := range b[i] {
			b[i][j] = game.Empty
		}
	}
	return b
}

func midGameBoard() game.Board {
	b := emptyBoard()
	// Place some stones for a simple mid-game scenario
	b[2][2] = game.Black
	b[2][3] = game.White
	b[3][2] = game.White
	b[3][3] = game.Black
	b[4][4] = game.Black
	b[4][5] = game.White
	b[5][4] = game.White
	b[5][5] = game.Black
	return b
}

func BenchmarkAlphaBetaEngine_EmptyBoard(b *testing.B) {
	engine := NewAlphaBetaEngine()
	board := emptyBoard()
	for i := 0; i < b.N; i++ {
		engine.Move(board, game.Black, nil)
	}
}

func BenchmarkAlphaBetaEngine_MidGame(b *testing.B) {
	engine := NewAlphaBetaEngine()
	board := midGameBoard()
	for i := 0; i < b.N; i++ {
		engine.Move(board, game.White, nil)
	}
}
