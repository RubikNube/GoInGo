// Package game provides the GUI for the game of the gods.
package game

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
	Empty: " ",
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

func (g *Gui) DrawGrid() {
	g.ClearScreen()
	// Draw column labels
	print("   ")
	for j := range [9]int{} {
		print(j + 1)
		if j < 8 {
			print("   ")
		}
	}
	println()
	// Draw top border
	print("  ┌")
	for j := range [9]int{} {
		print("───")
		if j < 8 {
			print("┬")
		}
	}
	println("┐")
	for i := range [9]int{} {
		// Draw row label
		print(i+1, " │")
		for j := range [9]int{} {
			char := g.Grid[i][j].String()
			print(" ", char, " ")
			if j < 8 {
				print("│")
			}
		}
		println("│")
		// Draw row separator or bottom border
		if i < 8 {
			print("  ├")
			for j := range [9]int{} {
				print("───")
				if j < 8 {
					print("┼")
				}
			}
			println("┤")
		} else {
			print("  └")
			for j := range [9]int{} {
				print("───")
				if j < 8 {
					print("┴")
				}
			}
			println("┘")
		}
	}
	g.Refresh()
}

func (g *Gui) DrawGridWithCursor(cursorRow, cursorCol int) {
	g.ClearScreen()
	// Draw column labels
	print("   ")
	for j := range [9]int{} {
		print(j + 1)
		if j < 8 {
			print("   ")
		}
	}
	println()
	// Draw top border
	print("  ┌")
	for j := range [9]int{} {
		print("───")
		if j < 8 {
			print("┬")
		}
	}
	println("┐")
	for i := range [9]int{} {
		// Draw row label
		print(i+1, " │")
		for j := range [9]int{} {
			char := g.Grid[i][j].String()
			if i == cursorRow && j == cursorCol {
				print("[", char, "]")
			} else {
				print(" ", char, " ")
			}
			if j < 8 {
				print("│")
			}
		}
		println("│")
		// Draw row separator or bottom border
		if i < 8 {
			print("  ├")
			for j := range [9]int{} {
				print("───")
				if j < 8 {
					print("┼")
				}
			}
			println("┤")
		} else {
			print("  └")
			for j := range [9]int{} {
				print("───")
				if j < 8 {
					print("┴")
				}
			}
			println("┘")
		}
	}
	g.Refresh()
}
