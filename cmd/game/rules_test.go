package game

import (
	"fmt"
	"testing"
)

func TestNeighborsTableDriven(t *testing.T) {
	tests := []struct {
		point    Point
		expected int
	}{
		{Point{Row: 0, Col: 0}, 2}, // top-left corner
		{Point{Row: 0, Col: 1}, 3}, // top edge
		{Point{Row: 0, Col: 8}, 2}, // top-right corner (assuming 19x19 board)
		{Point{Row: 1, Col: 0}, 3}, // left edge
		{Point{Row: 1, Col: 1}, 4}, // center
		{Point{Row: 8, Col: 0}, 2}, // bottom-left corner
		{Point{Row: 8, Col: 8}, 2}, // bottom-right corner
		{Point{Row: 8, Col: 1}, 3}, // bottom edge
		{Point{Row: 1, Col: 8}, 3}, // right edge
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%+v", tc.point), func(t *testing.T) {
			checkNeighbors(t, tc.point, tc.expected)
		})
	}
}

func checkNeighbors(t *testing.T, p Point, expected int) {
	n := Neighbors(p)
	if len(n) != expected {
		t.Errorf("Expected %d neighbors for %+v, got %d", expected, p, len(n))
	}
}

func TestIsLegalMove(t *testing.T) {
	var b Board
	prev := b
	move := Point{Row: 4, Col: 4}
	if !IsLegalMove(b, move, Black, prev) {
		t.Errorf("Expected move to be legal")
	}
}

func TestIsLegalMoveOccupied(t *testing.T) {
	var b Board
	b[4][4] = Black
	prev := b
	move := Point{Row: 4, Col: 4}
	if IsLegalMove(b, move, White, prev) {
		t.Errorf("Expected move to be illegal (occupied)")
	}
}

func TestIsLegalMoveSuicide(t *testing.T) {
	var b Board
	// Surround (3,3) with Black stones
	b[2][3], b[3][2], b[3][4], b[4][3] = Black, Black, Black, Black
	prev := b
	move := Point{Row: 3, Col: 3}
	if IsLegalMove(b, move, White, prev) {
		t.Errorf("Expected move to be illegal (suicide)")
	}
}

func TestIsLegalMoveKo(t *testing.T) {
	var b Board
	// Setup a simple Ko situation:
	// Black at (1,0), (0,1)
	// White at (0,0), (1,1)
	b[1][0] = Black
	b[0][1] = Black
	b[0][0] = White
	b[1][1] = White

	// Previous board state before White captured at (0,0)
	var prev Board
	prev[1][0] = Black
	prev[0][1] = Black
	prev[1][1] = White

	// Now Black tries to recapture at (0,0) immediately (should be illegal due to Ko)
	move := Point{Row: 0, Col: 0}
	if IsLegalMove(b, move, Black, prev) {
		t.Errorf("Expected move to be illegal due to Ko rule")
	}
}

func TestGroupAndLibertiesSingleStone(t *testing.T) {
	var b Board
	b[0][0] = Black
	stones, libs := Group(b, Point{0, 0})
	if len(stones) != 1 || len(libs) != 2 {
		t.Errorf("Expected 1 stone and 2 liberties, got %d stones and %d liberties", len(stones), len(libs))
	}
}

func TestGroupAndLibertiesConnectedStones(t *testing.T) {
	var b Board
	b[1][1] = Black
	b[1][2] = Black
	stones, libs := Group(b, Point{1, 1})
	if len(stones) != 2 {
		t.Errorf("Expected 2 stones, got %d", len(stones))
	}
	// (1,0), (0,1), (2,1), (1,3), (0,2), (2,2) are liberties
	if len(libs) != 6 {
		t.Errorf("Expected 6 liberties, got %d", len(libs))
	}
}

func TestGroupAndLibertiesSurroundedGroup(t *testing.T) {
	var b Board
	b[2][2] = Black
	b[2][3] = Black
	b[1][2], b[1][3], b[2][1], b[2][4], b[3][2], b[3][3] = White, White, White, White, White, White
	stones, libs := Group(b, Point{2, 2})
	if len(stones) != 2 {
		t.Errorf("Expected 2 stones, got %d", len(stones))
	}
	if len(libs) != 0 {
		t.Errorf("Expected 0 liberties, got %d", len(libs))
	}
}
