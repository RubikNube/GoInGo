// Package game provides a simple Go engine interface and a random-move engine.
package game

import (
	"math/rand"
	"time"
)

// Engine is an interface for Go engines.
type Engine interface {
	// Move returns the next move as a Point, or nil if passing.
	Move(board Board, player FieldState, ko *Point) *Point
}

// RandomEngine implements Engine by picking a random legal move.
type RandomEngine struct{}

func (e *RandomEngine) Move(board Board, player FieldState, ko *Point) *Point {
	empty := []Point{}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if board[i][j] == Empty {
				empty = append(empty, Point{Row: i, Col: j})
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
		var nextBoard Board
		copy(nextBoard[:], board[:])
		nextBoard[pt.Row][pt.Col] = player
		opp := Black
		if player == Black {
			opp = White
		}
		captured := []Point{}
		for _, n := range Neighbors(pt) {
			if nextBoard[n.Row][n.Col] == opp {
				group, libs := Group(nextBoard, n)
				if len(libs) == 0 {
					for stonePt := range group {
						nextBoard[stonePt.Row][stonePt.Col] = Empty
						captured = append(captured, stonePt)
					}
				}
			}
		}
		_, libs := Group(nextBoard, pt)
		if len(libs) == 0 {
			continue
		}
		// Legal move found
		return &pt
	}
	// No legal move, pass
	return nil
}
