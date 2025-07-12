package engine

import "github.com/RubikNube/GoInGo/pkg/game"

func EmptyBoard() game.Board {
	var b game.Board
	for i := range b {
		for j := range b[i] {
			b[i][j] = game.Empty
		}
	}
	return b
}

func MidGameBoard() game.Board {
	b := EmptyBoard()
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
