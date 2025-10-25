package engine

import (
	"math/rand"
	"sort"
	"sync"

	"github.com/RubikNube/GoInGo/pkg/game"
)

const (
	minScore = -1 << 30
	maxScore = 1 << 30
)

type ttFlag uint8

const (
	ttExact ttFlag = iota
	ttLowerBound
	ttUpperBound
)

type ttEntry struct {
	score    int
	depth    int
	flag     ttFlag
	bestMove game.Point
	hasMove  bool
}

// AlphaBetaEngine implements Engine using alpha-beta pruning with killer move heuristic, transposition table, and history heuristic.
type AlphaBetaEngine struct {
	killerMoves        map[int]game.Point // depth -> killer move
	transpositionTable map[uint64]ttEntry // board hash -> entry
	historyHeuristic   map[game.Point]int // move -> score for ordering
}

func NewAlphaBetaEngine() *AlphaBetaEngine {
	return &AlphaBetaEngine{
		killerMoves:        make(map[int]game.Point),
		transpositionTable: make(map[uint64]ttEntry),
		historyHeuristic:   make(map[game.Point]int),
	}
}

// Move in AlphaBetaEngine uses alpha-beta pruning to select the best move or pass if no beneficial move exists.
func (e *AlphaBetaEngine) Move(board game.Board, player game.FieldState, ko *game.Point) *game.Point {
	bestScore := minScore
	var bestMove *game.Point
	depth := 4 // Shallow for performance; increase for stronger player
	moveFound := false

	// Ensure killerMoves map is initialized
	if e.killerMoves == nil {
		e.killerMoves = make(map[int]game.Point)
	}
	if e.transpositionTable == nil {
		e.transpositionTable = make(map[uint64]ttEntry)
	}
	if e.historyHeuristic == nil {
		e.historyHeuristic = make(map[game.Point]int)
	}

	opp := opponent(player)
	seen := make(map[game.Point]struct{})
	tryMove := func(pt game.Point) {
		if _, exists := seen[pt]; exists {
			return
		}
		seen[pt] = struct{}{}
		if board[pt.Row][pt.Col] != game.Empty {
			return
		}
		if ko != nil && pt.Row == ko.Row && pt.Col == ko.Col {
			return
		}
		nextBoard := board
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
		if _, libs := game.Group(nextBoard, pt); len(libs) == 0 {
			return
		}
		moveFound = true
		score := -e.alphaBeta(nextBoard, opp, player, ko, depth-1, minScore, maxScore)
		if score > bestScore {
			bestScore = score
			move := pt
			bestMove = &move
		}
	}

	if entry, ok := e.transpositionTable[boardHash(board, player)]; ok && entry.hasMove {
		tryMove(entry.bestMove)
	}
	for _, pt := range e.orderedMoves(board, player, depth) {
		tryMove(pt)
	}

	// Pass if no move found or if passing is as good or better than any move
	passScore := -e.alphaBeta(board, opp, player, ko, depth-1, minScore, maxScore)
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

	boardHashValue := boardHash(board, player)
	entry, entryFound := e.transpositionTable[boardHashValue]
	if entryFound && entry.depth >= depth {
		switch entry.flag {
		case ttExact:
			return entry.score
		case ttLowerBound:
			if entry.score > alpha {
				alpha = entry.score
			}
		case ttUpperBound:
			if entry.score < beta {
				beta = entry.score
			}
		}
		if alpha >= beta {
			return entry.score
		}
	}
	alphaStart := alpha
	betaStart := beta

	// Null Move Pruning: try skipping a move (pass) if depth is sufficient
	if depth >= 2 {
		passScore := -e.alphaBeta(board, opp, player, ko, depth-2, -beta, -beta+1)
		if passScore >= beta {
			e.transpositionTable[boardHashValue] = ttEntry{
				score:   passScore,
				depth:   depth,
				flag:    ttLowerBound,
				hasMove: false,
			}
			return passScore
		}
	}

	bestValue := minScore
	bestMove := game.Point{}
	bestMoveValid := false
	seen := make(map[game.Point]struct{})

	tryMove := func(pt game.Point) bool {
		if _, exists := seen[pt]; exists {
			return false
		}
		seen[pt] = struct{}{}
		if board[pt.Row][pt.Col] != game.Empty {
			return false
		}
		if ko != nil && pt.Row == ko.Row && pt.Col == ko.Col {
			return false
		}
		nextBoard := board
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
		if _, libs := game.Group(nextBoard, pt); len(libs) == 0 {
			return false
		}
		foundMove = true
		score := -e.alphaBeta(nextBoard, opp, player, ko, depth-1, -beta, -alpha)
		// History heuristic update to encourage good moves.
		e.historyHeuristic[pt] += 1 << uint(depth)

		if score > bestValue {
			bestValue = score
			bestMove = pt
			bestMoveValid = true
		}
		if score > alpha {
			alpha = score
			if alpha >= beta {
				e.killerMoves[depth] = pt
				flag := ttExact
				switch {
				case bestValue <= alphaStart:
					flag = ttUpperBound
				case bestValue >= betaStart:
					flag = ttLowerBound
				}
				e.transpositionTable[boardHashValue] = ttEntry{
					score:    bestValue,
					depth:    depth,
					flag:     flag,
					bestMove: bestMove,
					hasMove:  bestMoveValid,
				}
				return true
			}
		}
		return false
	}

	if entryFound && entry.hasMove {
		if tryMove(entry.bestMove) {
			return bestValue
		}
	}

	if killer, ok := e.killerMoves[depth]; ok {
		if tryMove(killer) {
			return bestValue
		}
	}

	for _, pt := range e.orderedMoves(board, player, depth) {
		if tryMove(pt) {
			return bestValue
		}
	}

	passScore := -e.alphaBeta(board, opp, player, ko, depth-1, -beta, -alpha)
	if !foundMove || passScore > alpha {
		alpha = passScore
	}
	if !foundMove || passScore > bestValue {
		bestValue = passScore
		bestMoveValid = false
	}
	flag := ttExact
	switch {
	case bestValue <= alphaStart:
		flag = ttUpperBound
	case bestValue >= betaStart:
		flag = ttLowerBound
	}
	e.transpositionTable[boardHashValue] = ttEntry{
		score:    bestValue,
		depth:    depth,
		flag:     flag,
		bestMove: bestMove,
		hasMove:  bestMoveValid,
	}
	return bestValue
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
			if hasKiller && pt.Row == killer.Row && pt.Col == killer.Col {
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

// boardHash returns a Zobrist hash for the board and player to reduce collisions.
var (
	zobristOnce   sync.Once
	zobristTable  [game.BoardSize][game.BoardSize][3]uint64
	zobristPlayer [3]uint64
)

func initZobrist() {
	r := rand.New(rand.NewSource(0x102030405060708))
	for i := 0; i < game.BoardSize; i++ {
		for j := 0; j < game.BoardSize; j++ {
			for k := 0; k < 3; k++ {
				zobristTable[i][j][k] = r.Uint64()
			}
		}
	}
	for k := 0; k < 3; k++ {
		zobristPlayer[k] = r.Uint64()
	}
}

func boardHash(board game.Board, player game.FieldState) uint64 {
	zobristOnce.Do(initZobrist)
	var h uint64
	for i := 0; i < game.BoardSize; i++ {
		for j := 0; j < game.BoardSize; j++ {
			stone := board[i][j]
			if stone != game.Empty {
				h ^= zobristTable[i][j][int(stone)]
			}
		}
	}
	h ^= zobristPlayer[int(player)]
	return h
}
