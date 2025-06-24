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
		fmt.Print("Move (h/j/k/l, q to quit): ")
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
		case "q":
			fmt.Println("Quitting game.")
			return
		}
		gui.DrawGridWithCursor(cursorRow, cursorCol)
	}
}
