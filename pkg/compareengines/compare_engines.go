// Package compareengines provides functionality to compare two Go engines by playing a game.
package compareengines

import (
	"github.com/RubikNube/GoInGo/pkg/engine"
	"github.com/RubikNube/GoInGo/pkg/game"
)

// CompareEngines lets two engines play against each other and returns the winner.
// Returns: 1 if engineA wins, -1 if engineB wins, 0 for draw.
func CompareEngines(engineA, engineB engine.Engine, board game.Board, firstPlayer game.FieldState, maxMoves int) int {
	player := firstPlayer
	var ko *game.Point
	moveCount := 0
	passCount := 0
	for moveCount < maxMoves && passCount < 2 {
		var move *game.Point
		if player == firstPlayer {
			move = engineA.Move(board, player, ko)
		} else {
			move = engineB.Move(board, player, ko)
		}
		if move == nil {
			passCount++
		} else {
			passCount = 0
			board[move.Row][move.Col] = player
			ko = nil // Optionally update ko if needed
		}
		player = opponent(player)
		moveCount++
	}
	score := evaluate(board, firstPlayer, opponent(firstPlayer))
	if score > 0 {
		return 1
	} else if score < 0 {
		return -1
	}
	return 0
}

// opponent returns the opposite FieldState (Black <-> White).
func opponent(player game.FieldState) game.FieldState {
	if player == game.Black {
		return game.White
	}
	return game.Black
}

// evaluate returns a score for the board from the perspective of 'player'.
// Positive means advantage for 'player', negative for opponent.
func evaluate(board game.Board, player, opp game.FieldState) int {
	playerStones, oppStones := 0, 0
	playerLibs, oppLibs := 0, 0
	visited := make(map[game.Point]bool)
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			pt := game.Point{Row: int8(i), Col: int8(j)}
			if visited[pt] || board[i][j] == game.Empty {
				continue
			}
			group, libs := game.Group(board, pt)
			for stone := range group {
				visited[stone] = true
			}
			if board[i][j] == player {
				playerStones += len(group)
				playerLibs += len(libs)
			} else if board[i][j] == opp {
				oppStones += len(group)
				oppLibs += len(libs)
			}
		}
	}
	return (playerStones-oppStones)*10 + (playerLibs-oppLibs)*2
}
