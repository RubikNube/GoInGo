package old

import (
	"testing"

	"github.com/RubikNube/GoInGo/pkg/engine"
	"github.com/RubikNube/GoInGo/pkg/game"
)

func BenchmarkAlphaBetaEngine_EmptyBoard(b *testing.B) {
	engine2 := NewAlphaBetaEngine()
	board := engine.EmptyBoard()
	for i := 0; i < b.N; i++ {
		engine2.Move(board, game.Black, nil)
	}
}

func BenchmarkAlphaBetaEngine_MidGame(b *testing.B) {
	engine2 := NewAlphaBetaEngine()
	board := engine.MidGameBoard()
	for i := 0; i < b.N; i++ {
		engine2.Move(board, game.White, nil)
	}
}
