package main

import "fmt"
import "github.com/RubikNube/GoInGo/cmd/game"

func main() {
	// This is a placeholder for the main function.
	fmt.Println("Welcome to the Game of the Gods!")
	gui := game.Gui{}
	cursorRow, cursorCol := 0, 0

	gui.DrawGridWithCursor(cursorRow, cursorCol)

	for {
		var input string
		fmt.Print("Move (h/j/k/l, p to place stone, q to quit): ")
		fmt.Scanln(&input)

		switch input {
		case "h":
			if cursorCol > 0 {
				cursorCol--
			}
		case "l":
			if cursorCol < 8 {
				cursorCol++
			}
		case "k":
			if cursorRow > 0 {
				cursorRow--
			}
		case "j":
			if cursorRow < 8 {
				cursorRow++
			}
		case "p":
			// Place a stone if the cell is empty
			if gui.Grid[cursorRow][cursorCol] == game.Empty {
				// Alternate between Black and White stones
				stone := game.Black
				stoneCount := 0
				for i := 0; i < 9; i++ {
					for j := 0; j < 9; j++ {
						if gui.Grid[i][j] == game.Black || gui.Grid[i][j] == game.White {
							stoneCount++
						}
					}
				}
				if stoneCount%2 == 1 {
					stone = game.White
				}
				gui.Grid[cursorRow][cursorCol] = stone
			}
		case "q":
			fmt.Println("Quitting game.")
			return
		}
		gui.DrawGridWithCursor(cursorRow, cursorCol)
	}
}
