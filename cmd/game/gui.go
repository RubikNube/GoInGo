// Package game provides the GUI for the game of the gods.
package game

import (
	"fmt"
	"io"
)

// FieldState A field an either empty or occupied by a black or white stone.
type FieldState uint8

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
	Black: "\033[1m⚫\033[0m",
	White: "\033[1m⚪\033[0m",
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
	for j := range g.Grid[0] {
		fmt.Fprintf(w, " %c  ", 'A'+j)
	}
	fmt.Fprintln(w)

	for i, row := range g.Grid {
		// Row label
		fmt.Fprintf(w, "%2d ", i+1)
		for j, cellVal := range row {
			stone := cellVal.String()
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
					// Mark the center intersection bold if empty
					if i == 2 && j == 2 || i == 4 && j == 4 || i == 6 && j == 6 || i == 2 && j == 6 || i == 6 && j == 2 {

						stone = "\033[1m" + g.Grid[i][j].String() + "\033[0m"
					} else {
						stone = g.Grid[i][j].String()
					}
				}
			}
			cell := stone

			if j == 0 {
				cell = fmt.Sprintf(" %s─", stone)
			} else if j == 8 {
				cell = fmt.Sprintf("─%s ", stone)
			} else {
				// Use a box-drawing character for the stone
				cell = fmt.Sprintf("─%s─", stone)
			}
			if i == cursorRow && j == cursorCol {
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
			for j := range g.Grid[i] {
				fmt.Fprint(w, " │ ")
				if j < 8 {
					fmt.Fprint(w, " ")
				}
			}
			fmt.Fprintln(w)
		}
	}
}
