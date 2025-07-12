package game

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// Minimal GUI struct and constructor for demonstration.
// Replace with your actual implementation.
type GUI struct {
	CurrentPlayer string
}

func NewGUI() *GUI {
	return &GUI{CurrentPlayer: "Black"}
}

func (g *GUI) SwitchPlayer() {
	if g.CurrentPlayer == "Black" {
		g.CurrentPlayer = "White"
	} else {
		g.CurrentPlayer = "Black"
	}
}

func TestGuiInitialState(t *testing.T) {
	gui := NewGUI()
	if gui.CurrentPlayer != "Black" {
		t.Errorf("Expected initial player to be Black, got %v", gui.CurrentPlayer)
	}
}

func TestGuiSwitchPlayer(t *testing.T) {
	gui := NewGUI()
	gui.SwitchPlayer()
	if gui.CurrentPlayer != "White" {
		t.Errorf("Expected player to switch to White, got %v", gui.CurrentPlayer)
	}
	gui.SwitchPlayer()
	if gui.CurrentPlayer != "Black" {
		t.Errorf("Expected player to switch back to Black, got %v", gui.CurrentPlayer)
	}
}

func TestGuiMultipleSwitches(t *testing.T) {
	gui := NewGUI()
	for i := 0; i < 10; i++ {
		gui.SwitchPlayer()
	}
	expected := "Black"
	if 10%2 != 0 {
		expected = "White"
	}
	if gui.CurrentPlayer != expected {
		t.Errorf("After 10 switches, expected player to be %v, got %v", expected, gui.CurrentPlayer)
	}
}

func TestGuiSwitchPlayerAlternates(t *testing.T) {
	gui := NewGUI()
	players := []string{"White", "Black", "White", "Black"}
	for i, expected := range players {
		gui.SwitchPlayer()
		if gui.CurrentPlayer != expected {
			t.Errorf("After %d switches, expected player to be %v, got %v", i+1, expected, gui.CurrentPlayer)
		}
	}
}

func TestDrawGridToWriterEmptyBoard(t *testing.T) {
	var b Board
	var buf bytes.Buffer
	gui := Gui{}
	gui.Grid = b
	gui.DrawGridToWriter(&buf, 0, 0)
	output := buf.String()

	// Check for all column labels
	for j := 0; j < 9; j++ {
		colLabel := fmt.Sprintf("%c", 'A'+j)
		if !strings.Contains(output, colLabel) {
			t.Errorf("Expected column label %q in output", colLabel)
		}
	}

	// Check for all row labels
	for i := 1; i <= 9; i++ {
		rowLabel := fmt.Sprintf("%2d", i)
		if !strings.Contains(output, rowLabel) {
			t.Errorf("Expected row label %q in output", rowLabel)
		}
	}

	// Check for grid rendering characters
	gridChars := []string{"┌", "┐", "└", "┘", "┬", "┴", "├", "┤", "┼", "│", "─"}
	for _, ch := range gridChars {
		if !strings.Contains(output, ch) {
			t.Errorf("Expected grid character %q in output", ch)
		}
	}

	// Ensure no stones are present
	if strings.Contains(output, "\x1b[1m⚫\x1b[0m") || strings.Contains(output, "\x1b[1m⚪\x1b[0m") {
		t.Errorf("Expected empty board to have no stones, got: %q", output)
	}
}

func TestDrawGridToWriterWithBlackStone(t *testing.T) {
	var b Board
	b[0][0] = Black
	var buf bytes.Buffer
	gui := Gui{}
	gui.Grid = b
	gui.DrawGridToWriter(&buf, 0, 0)
	output := buf.String()
	if !strings.Contains(output, "\x1b[1m⚫\x1b[0m") {
		t.Errorf("Expected board to contain a black stone, got: %q", output)
	}
}

func TestDrawGridToWriterWithWhiteStone(t *testing.T) {
	var b Board
	b[0][0] = White
	var buf bytes.Buffer
	gui := Gui{}
	gui.Grid = b
	gui.DrawGridToWriter(&buf, 0, 0)
	output := buf.String()
	if !strings.Contains(output, "\x1b[1m⚪\x1b[0m") {
		t.Errorf("Expected board to contain a white stone, got: %q", output)
	}
}
