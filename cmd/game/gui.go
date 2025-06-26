// Package game provides the GUI for the game of the gods.
package game

import (
	"fmt"
	"io"
)

// FieldState A field an either empty or occupied by a black or white stone.
type FieldState int

// Constants for the FieldState type
const (
	Empty FieldState = iota // Empty field
	Black                   // Black stone
	White                   // White stone
)

type Gui struct {
	Grid [9][9]FieldState // 9x9 grid for the game
}

var FieldStateName = map[FieldState]string{
	Empty: "┼",
	Black: "○", // Unicode empty circle for black stone
	White: "●", // Unicode filled circle for white stone
}

func (fs FieldState) String() string {
	// This function returns the string representation of the FieldState.
	if name, ok := FieldStateName[fs]; ok {
		return name
	}
	return "Unknown"
}

func (g *Gui) ClearScreen() {
	// This function clears the terminal screen.
	// Its supported for bash, zsh, and other common shells.
	print("\033[H\033[2J")
}

func (g *Gui) PrintAt(row, col int, char string) {
	// This function prints a character at a specific position in the terminal.
	print("\033[", row+1, ";", col+1, "H", char)
}

func (g *Gui) Refresh() {
	// This function refreshes the terminal display.
	print("\033[0m")   // Reset terminal formatting
	print("\033[?25l") // Hide cursor
	print("\033[?25h") // Show cursor
}

func (g *Gui) DrawGridToWriter(w io.Writer, cursorRow, cursorCol int) {
	// Column labels
	fmt.Fprint(w, "   ")
	for j := 0; j < 9; j++ {
		fmt.Fprintf(w, " %c  ", 'A'+j)
	}
	fmt.Fprintln(w)

	for i := 0; i < 9; i++ {
		// Row label
		fmt.Fprintf(w, "%2d ", i+1)
		for j := 0; j < 9; j++ {
			stone := g.Grid[i][j].String()
			// Use box-drawing characters for borders and intersections
			if g.Grid[i][j] != Empty {
				stone = g.Grid[i][j].String()
			} else {
				switch {
				case i == 0 && j == 0:
					stone = "┌"
				case i == 0 && j == 8:
					stone = "┐"
				case i == 8 && j == 0:
					stone = "└"
				case i == 8 && j == 8:
					stone = "┘"
				case i == 0:
					stone = "┬"
				case i == 8:
					stone = "┴"
				case j == 0:
					stone = "├"
				case j == 8:
					stone = "┤"
				default:
					stone = g.Grid[i][j].String()
				}
			}
			var cell string
			if j == 0 {
				cell = fmt.Sprintf(" %s─", stone)
			} else if j == 8 {
				cell = fmt.Sprintf("─%s ", stone)
			} else {
				// Use a box-drawing character for the stone
				cell = fmt.Sprintf("─%s─", stone)
			}

			if i == cursorRow && j == cursorCol {
				// Use a different background or brackets, but keep width 3
				cell = fmt.Sprintf("[%s]", stone)
			}
			fmt.Fprint(w, cell)

			// Draw horizontal line except after last column
			if j < 8 {
				fmt.Fprint(w, "─")
			}
		}
		fmt.Fprintln(w)
		// Draw vertical lines except after last row
		if i < 8 {
			fmt.Fprint(w, "   ")
			for j := 0; j < 9; j++ {
				fmt.Fprint(w, " │ ")
				if j < 8 {
					fmt.Fprint(w, " ")
				}
			}
			fmt.Fprintln(w)
		}
	}
}
