package engine

import (
	"math/rand"
	"time"

	"github.com/RubikNube/GoInGo/cmd/game"
)

// RandomEngine implements Engine by picking a random legal move.
type RandomEngine struct{}

func (e *RandomEngine) Move(board game.Board, player game.FieldState, ko *game.Point) *game.Point {
	empty := []game.Point{}
	for i := int8(0); i < 9; i++ {
		for j := int8(0); j < 9; j++ {
			if board[i][j] == game.Empty {
				empty = append(empty, game.Point{Row: int8(i), Col: int8(j)})
			}
		}
	}
	rand.Seed(time.Now().UnixNano())
	perm := rand.Perm(len(empty))
	for _, idx := range perm {
		pt := empty[idx]
		// Ko rule
		if ko != nil && pt.Row == ko.Row && pt.Col == ko.Col {
			continue
		}
		var nextBoard game.Board
		copy(nextBoard[:], board[:])
		nextBoard[pt.Row][pt.Col] = player
		opp := game.Black
		if player == game.Black {
			opp = game.White
		}
		captured := []game.Point{}
		for _, n := range game.Neighbors(pt) {
			if nextBoard[n.Row][n.Col] == opp {
				group, libs := game.Group(nextBoard, n)
				if len(libs) == 0 {
					for stonePt := range group {
						nextBoard[stonePt.Row][stonePt.Col] = game.Empty
						captured = append(captured, stonePt)
					}
				}
			}
		}
		_, libs := game.Group(nextBoard, pt)
		if len(libs) == 0 {
			continue
		}
		// Legal move found
		return &pt
	}
	// No legal move, pass
	return nil
}
