package engine

import (
	"sort"

	"github.com/RubikNube/GoInGo/pkg/game"
)

// AlphaBetaEngine implements Engine using alpha-beta pruning with killer move heuristic, transposition table, and history heuristic.
type AlphaBetaEngine struct {
	killerMoves        map[int]*game.Point // depth -> killer move
	transpositionTable map[uint64]int      // board hash -> score
	historyHeuristic   map[game.Point]int  // move -> score for ordering
}

func NewAlphaBetaEngine() *AlphaBetaEngine {
	return &AlphaBetaEngine{
		killerMoves:        make(map[int]*game.Point),
		transpositionTable: make(map[uint64]int),
		historyHeuristic:   make(map[game.Point]int),
	}
}

// Move in AlphaBetaEngine uses alpha-beta pruning to select the best move or pass if no beneficial move exists.
func (e *AlphaBetaEngine) Move(board game.Board, player game.FieldState, ko *game.Point) *game.Point {
	bestScore := -1 << 30
	var bestMove *game.Point
	depth := 4 // Shallow for performance; increase for stronger player
	moveFound := false

	// Ensure killerMoves map is initialized
	if e.killerMoves == nil {
		e.killerMoves = make(map[int]*game.Point)
	}
	if e.transpositionTable == nil {
		e.transpositionTable = make(map[uint64]int)
	}
	if e.historyHeuristic == nil {
		e.historyHeuristic = make(map[game.Point]int)
	}

	for i := int8(0); i < 9; i++ {
		for j := int8(0); j < 9; j++ {
			if board[i][j] != game.Empty {
				continue
			}
			pt := game.Point{Row: i, Col: j}
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
			for _, n := range game.Neighbors(pt) {
				if nextBoard[n.Row][n.Col] == opp {
					group, libs := game.Group(nextBoard, n)
					if len(libs) == 0 {
						for stonePt := range group {
							nextBoard[stonePt.Row][stonePt.Col] = game.Empty
						}
					}
				}
			}
			_, libs := game.Group(nextBoard, pt)
			if len(libs) == 0 {
				continue
			}
			score := -e.alphaBeta(nextBoard, opp, player, ko, depth-1, -1<<30, 1<<30)
			moveFound = true
			if score > bestScore {
				bestScore = score
				move := pt
				bestMove = &move
			}
		}
	}
	// Pass if no move found or if passing is as good or better than any move
	passScore := -e.alphaBeta(board, opponent(player), player, ko, depth-1, -1<<30, 1<<30)
	if !moveFound || passScore >= bestScore {
		return nil // pass
	}
	return bestMove
}

// opponent returns the opposite FieldState (Black <-> White).
func opponent(player game.FieldState) game.FieldState {
	if player == game.Black {
		return game.White
	}
	return game.Black
}

// alphaBeta is a minimax search with alpha-beta pruning, killer move heuristic, transposition table, and history heuristic.
func (e *AlphaBetaEngine) alphaBeta(board game.Board, player, opp game.FieldState, ko *game.Point, depth, alpha, beta int) int {
	if depth == 0 {
		return evaluate(board, player, opp)
	}
	foundMove := false

	// Transposition table lookup
	boardHash := boardHash(board, player)
	if val, ok := e.transpositionTable[boardHash]; ok {
		return val
	}

	// Null Move Pruning: try skipping a move (pass) if depth is sufficient
	if depth >= 2 {
		passScore := -e.alphaBeta(board, opp, player, ko, depth-2, -beta, -beta+1)
		if passScore >= beta {
			e.transpositionTable[boardHash] = passScore
			return passScore
		}
	}

	// Try killer move first if available
	if killer, ok := e.killerMoves[depth]; ok && killer != nil && board[killer.Row][killer.Col] == game.Empty {
		pt := *killer
		if ko == nil || pt.Row != ko.Row || pt.Col != ko.Col {
			var nextBoard game.Board
			copy(nextBoard[:], board[:])
			nextBoard[pt.Row][pt.Col] = player
			for _, n := range game.Neighbors(pt) {
				if nextBoard[n.Row][n.Col] == opp {
					group, libs := game.Group(nextBoard, n)
					if len(libs) == 0 {
						for stonePt := range group {
							nextBoard[stonePt.Row][stonePt.Col] = game.Empty
						}
					}
				}
			}
			_, libs := game.Group(nextBoard, pt)
			if len(libs) != 0 {
				foundMove = true
				score := -e.alphaBeta(nextBoard, opp, player, ko, depth-1, -beta, -alpha)
				// History heuristic update
				e.historyHeuristic[pt] += 1 << uint(depth)
				if score > alpha {
					alpha = score
					// Update killer move if this move caused a beta cutoff
					if alpha >= beta {
						e.killerMoves[depth] = &pt
						e.transpositionTable[boardHash] = alpha
						return alpha
					}
				}
			}
		}
	}

	for _, pt := range e.orderedMoves(board, player, depth) {
		if board[pt.Row][pt.Col] != game.Empty {
			continue
		}
		if ko != nil && pt.Row == ko.Row && pt.Col == ko.Col {
			continue
		}
		// Skip killer move (already tried)
		if killer, ok := e.killerMoves[depth]; ok && killer != nil && pt.Row == killer.Row && pt.Col == killer.Col {
			continue
		}
		var nextBoard game.Board
		copy(nextBoard[:], board[:])
		nextBoard[pt.Row][pt.Col] = player
		for _, n := range game.Neighbors(pt) {
			if nextBoard[n.Row][n.Col] == opp {
				group, libs := game.Group(nextBoard, n)
				if len(libs) == 0 {
					for stonePt := range group {
						nextBoard[stonePt.Row][stonePt.Col] = game.Empty
					}
				}
			}
		}
		_, libs := game.Group(nextBoard, pt)
		if len(libs) == 0 {
			continue
		}
		foundMove = true
		score := -e.alphaBeta(nextBoard, opp, player, ko, depth-1, -beta, -alpha)
		// History heuristic update
		e.historyHeuristic[pt] += 1 << uint(depth)
		if score > alpha {
			alpha = score
			// Update killer move if this move caused a beta cutoff
			if alpha >= beta {
				move := pt
				e.killerMoves[depth] = &move
				e.transpositionTable[boardHash] = alpha
				return alpha
			}
		}
	}
	// Consider passing if no move found or passing is better
	passScore := -e.alphaBeta(board, opp, player, ko, depth-1, -beta, -alpha)
	if !foundMove || passScore > alpha {
		alpha = passScore
	}
	e.transpositionTable[boardHash] = alpha
	return alpha
}

// orderedMoves returns a list of all empty points, ordered by killer move, history heuristic, proximity, and capture potential.
func (e *AlphaBetaEngine) orderedMoves(board game.Board, player game.FieldState, depth int) []game.Point {
	type moveScore struct {
		pt    game.Point
		score int
	}
	var moves []moveScore
	killer, hasKiller := e.killerMoves[depth]
	for i := int8(0); i < 9; i++ {
		for j := int8(0); j < 9; j++ {
			if board[i][j] != game.Empty {
				continue
			}
			pt := game.Point{Row: i, Col: j}
			score := 0
			// Killer move gets highest priority
			if hasKiller && killer != nil && pt.Row == killer.Row && pt.Col == killer.Col {
				score += 10000
			}
			// History heuristic
			score += e.historyHeuristic[pt] * 10
			// Proximity: +1 for each neighbor that is not empty
			for _, n := range game.Neighbors(pt) {
				if board[n.Row][n.Col] != game.Empty {
					score += 2
				}
			}
			// Capture potential: +5 for each neighbor group with 1 liberty
			opp := game.Black
			if player == game.Black {
				opp = game.White
			}
			for _, n := range game.Neighbors(pt) {
				if board[n.Row][n.Col] == opp {
					_, libs := game.Group(board, n)
					if len(libs) == 1 {
						score += 5
					}
				}
			}
			moves = append(moves, moveScore{pt, score})
		}
	}
	// Sort moves by descending score
	sort.Slice(moves, func(i, j int) bool {
		return moves[i].score > moves[j].score
	})
	result := make([]game.Point, len(moves))
	for i, m := range moves {
		result[i] = m.pt
	}
	return result
}

// evaluate is a sophisticated evaluation function considering liberties, groups, and captures.
func evaluate(board game.Board, player, opp game.FieldState) int {
	playerStones, oppStones := 0, 0
	playerLibs, oppLibs := 0, 0
	playerGroups, oppGroups := 0, 0
	playerCapturable, oppCapturable := 0, 0

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
				playerGroups++
				if len(libs) == 1 {
					playerCapturable += len(group)
				}
			} else if board[i][j] == opp {
				oppStones += len(group)
				oppLibs += len(libs)
				oppGroups++
				if len(libs) == 1 {
					oppCapturable += len(group)
				}
			}
		}
	}
	// Weighted sum: stones, liberties, groups, capturability
	return (playerStones-oppStones)*10 +
		(playerLibs-oppLibs)*2 +
		(oppCapturable-playerCapturable)*8 +
		(playerGroups - oppGroups)
}

// boardHash returns a simple hash for the board and player.
// You may want to replace this with Zobrist hashing for better collision resistance.
func boardHash(board game.Board, player game.FieldState) uint64 {
	var h uint64
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			h = h*3 + uint64(board[i][j])
		}
	}
	h = h*3 + uint64(player)
	return h
}
