package game

// BoardSize defines the size of the Go board.
const BoardSize = 9

// Point represents a coordinate on the board.
type Point struct {
	Row int
	Col int
}

// Board is a BoardSize x BoardSize Go board.
type Board [BoardSize][BoardSize]FieldState

// Neighbors returns the adjacent points of a given point.
func Neighbors(p Point) []Point {
	var n []Point
	dirs := []struct{ dr, dc int }{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for _, d := range dirs {
		r, c := p.Row+d.dr, p.Col+d.dc
		if r >= 0 && r < BoardSize && c >= 0 && c < BoardSize {
			n = append(n, Point{r, c})
		}
	}
	return n
}

// Group returns all stones connected to the given point and their liberties.
func Group(b Board, start Point) (stones map[Point]struct{}, liberties map[Point]struct{}) {
	color := b[start.Row][start.Col]
	stones = make(map[Point]struct{})
	liberties = make(map[Point]struct{})
	var stack []Point
	stack = append(stack, start)
	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if _, seen := stones[p]; seen {
			continue
		}
		stones[p] = struct{}{}
		for _, n := range Neighbors(p) {
			switch b[n.Row][n.Col] {
			case Empty:
				liberties[n] = struct{}{}
			case color:
				if _, seen := stones[n]; !seen {
					stack = append(stack, n)
				}
			}
		}
	}
	return
}

// IsLegalMove checks if placing a stone at p for color is legal (no suicide, no ko).
func IsLegalMove(b Board, p Point, color FieldState, prev Board) bool {
	if b[p.Row][p.Col] != Empty {
		return false
	}
	// Copy board and play move
	var next Board = b
	next[p.Row][p.Col] = color
	// Remove opponent groups with no liberties
	opp := Black
	if color == Black {
		opp = White
	}
	for _, n := range Neighbors(p) {
		if next[n.Row][n.Col] == opp {
			group, libs := Group(next, n)
			if len(libs) == 0 {
				for stone := range group {
					next[stone.Row][stone.Col] = Empty
				}
			}
		}
	}
	// Check if own group has liberties
	_, libs := Group(next, p)
	if len(libs) == 0 {
		return false // suicide
	}
	// Ko: board must not repeat previous position
	same := true
	for i := 0; i < BoardSize; i++ {
		for j := 0; j < BoardSize; j++ {
			if next[i][j] != prev[i][j] {
				same = false
				break
			}
		}
		if !same {
			break
		}
	}
	if same {
		return false // Ko: position repeats previous
	}
	return true
}
