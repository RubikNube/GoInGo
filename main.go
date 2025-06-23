package main

import "fmt"
import "github.com/RubikNube/GoInGo/cmd/game"
import "math/rand"

func main() {
	// This is a placeholder for the main function.
	// You can add your application logic here.
	fmt.Println("Welcome to the Game of the Gods!")
	gui := game.Gui{}
	// initialize the grid with some random values
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			// initialize each field with a random value
			gui.Grid[i][j] = game.FieldState(rand.Intn(3)) // Randomly assign 0, 1, or 2 as FieldState
		}
	}
	gui.DrawGrid()
}
