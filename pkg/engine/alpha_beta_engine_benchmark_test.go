package engine

import (
	"testing"

	"github.com/RubikNube/GoInGo/pkg/game"
)

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
