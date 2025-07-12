package engine

import (
	"testing"

	"github.com/RubikNube/GoInGo/pkg/game"
)

func BenchmarkRandomEngine_EmptyBoard(b *testing.B) {
	engine := NewRandomEngine()
	board := EmptyBoard()
	for i := 0; i < b.N; i++ {
		engine.Move(board, game.Black, nil)
	}
}

func BenchmarkRandomEngine_MidGame(b *testing.B) {
	engine := NewRandomEngine()
	board := MidGameBoard()
	for i := 0; i < b.N; i++ {
		engine.Move(board, game.White, nil)
	}
}
