package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	"unicode"

	"github.com/RubikNube/GoInGo/pkg/engine"
	"github.com/RubikNube/GoInGo/pkg/game"
	"github.com/jroimartin/gocui"
)

type Config struct {
	Keybindings map[string]string `json:"keybindings"`
}

var (
	cursorRow, cursorCol int8
	gui                  game.Gui
	keybindings          map[string]string
	prevBoard            *game.Board       // Track previous board for Ko rule
	currentPlayer        int8          = 1 // Track current player (1 or 2), start with Black
	koPoint              *game.Point       // Track Ko point (nil if no Ko)
	passCount            int8              // Track consecutive passes
	gameOver             bool              // Track if the game is over
	engineEnabled        bool              // Play against engine if true
	selectedEngine       engine.Engine     // The engine instance
)

func loadConfig(path string) (Config, error) {
	var cfg Config
	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&cfg)
	return cfg, err
}

func printMovePrompt(v *gocui.View) {
	fmt.Fprintf(v, "Move (%s/%s/%s/%s), %s to place stone, %s to pass, %s to quit", keybindings["moveLeft"], keybindings["moveDown"], keybindings["moveUp"], keybindings["moveRight"], keybindings["placeStone"], keybindings["passTurn"], keybindings["quit"])
}

func layout(g *gocui.Gui) error {
	maxX, _ := g.Size()
	boardHeight := 22 // enough for 9x9 grid with borders
	if v, err := g.SetView("board", 0, 0, maxX-1, boardHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Go (Baduk)"
		v.Wrap = false
	}
	if v, err := g.SetView("prompt", 0, boardHeight+1, maxX-1, boardHeight+3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = false
		printMovePrompt(v)
	}
	// Always redraw board
	if v, err := g.View("board"); err == nil {
		v.Clear()
		gui.DrawGridToWriter(v, cursorRow, cursorCol)
	}
	return nil
}

func moveCursor(dRow, dCol int8, jumpOverOccupied bool) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		nextRow, nextCol := cursorRow, cursorCol
		for {
			nextRow += dRow
			nextCol += dCol
			if nextRow < 0 {
				nextRow = 0
				break
			}
			if nextRow > 8 {
				nextRow = 8
				break
			}
			if nextCol < 0 {
				nextCol = 0
				break
			}
			if nextCol > 8 {
				nextCol = 8
				break
			}
			if jumpOverOccupied {
				if gui.Grid[nextRow][nextCol] == game.Empty {
					cursorRow, cursorCol = nextRow, nextCol
					break
				}
				// If we hit the edge and still not empty, stop
				if (dRow != 0 && (nextRow == 0 || nextRow == 8)) || (dCol != 0 && (nextCol == 0 || nextCol == 8)) {
					break
				}
			} else {
				cursorRow, cursorCol = nextRow, nextCol
				break
			}
		}
		return nil
	}
}

func placeStone(g *gocui.Gui, v *gocui.View) error {
	if gameOver {
		return nil
	}
	if gui.Grid[cursorRow][cursorCol] != game.Empty {
		return nil
	}
	// Ko rule: forbid move at koPoint
	if koPoint != nil && int8(cursorRow) == koPoint.Row && int8(cursorCol) == koPoint.Col {
		if v, err := g.View("prompt"); err == nil && v != nil {
			v.Clear()
			fmt.Fprint(v, "Illegal move! Ko rule.")
			go func() {
				time.Sleep(1 * time.Second)
				g.Update(func(g *gocui.Gui) error {
					if v, err := g.View("prompt"); err == nil && v != nil {
						v.Clear()
						printMovePrompt(v)
					}
					return nil
				})
			}()
		}
		return nil
	}

	stone := game.Black
	if currentPlayer == 2 {
		stone = game.White
	}
	// Track previous board for Ko rule (simple implementation: store last board)
	if prevBoard == nil {
		prev := game.Board{}
		copy(prev[:], gui.Grid[:])
		prevBoard = &prev
	}
	// Simulate the move and captures on a copy of the board
	var nextBoard game.Board
	copy(nextBoard[:], gui.Grid[:])
	nextBoard[cursorRow][cursorCol] = stone

	opp := game.Black
	if stone == game.Black {
		opp = game.White
	}
	captured := []game.Point{}
	for _, n := range game.Neighbors(game.Point{Row: int8(cursorRow), Col: int8(cursorCol)}) {
		if nextBoard[n.Row][n.Col] == opp {
			group, libs := game.Group(nextBoard, n)
			if len(libs) == 0 {
				for stonePt := range group {
					nextBoard[stonePt.Row][stonePt.Col] = game.Empty
					captured = append(captured, stonePt)
				}
			}
		}
	}

	// Check for liberties of the placed stone's group (suicide rule)
	_, libs := game.Group(nextBoard, game.Point{Row: int8(cursorRow), Col: int8(cursorCol)})
	if len(libs) == 0 {
		if v, err := g.View("prompt"); err == nil && v != nil {
			v.Clear()
			fmt.Fprint(v, "Illegal move! No liberties.")
			go func() {
				time.Sleep(1 * time.Second)
				g.Update(func(g *gocui.Gui) error {
					if v, err := g.View("prompt"); err == nil && v != nil {
						v.Clear()
						printMovePrompt(v)
					}
					return nil
				})
			}()
		}
		return nil
	}

	// Ko rule and legality check: resulting board must not match prevBoard
	if prevBoard != nil {
		same := true
		for i := int8(0); i < int8(len(nextBoard)); i++ {
			if nextBoard[i] != (*prevBoard)[i] {
				same = false
				break
			}
		}
		if same {
			if v, err := g.View("prompt"); err == nil && v != nil {
				v.Clear()
				fmt.Fprint(v, "Illegal move! Try again.")
				go func() {
					time.Sleep(1 * time.Second)
					g.Update(func(g *gocui.Gui) error {
						if v, err := g.View("prompt"); err == nil && v != nil {
							v.Clear()
							printMovePrompt(v)
						}
						return nil
					})
				}()
			}
			return nil
		}
	}

	// Place the stone and update the board
	gui.Grid[cursorRow][cursorCol] = stone
	for _, pt := range captured {
		gui.Grid[pt.Row][pt.Col] = game.Empty
	}

	// Ko rule: set koPoint if exactly one stone was captured and the group size is 1
	if len(captured) == 1 {
		koPoint = &captured[0]
	} else {
		koPoint = nil
	}

	// Update prevBoard for next move (for Ko rule)
	copy(prevBoard[:], nextBoard[:])

	passCount = 0 // Reset pass count on a move

	currentPlayer = 3 - currentPlayer // Switch player only after a legal move

	// If engine is enabled and it's the engine's turn, make engine move
	if engineEnabled && !gameOver && currentPlayer == 2 {
		go func() {
			time.Sleep(300 * time.Millisecond)
			g.Update(func(g *gocui.Gui) error {
				engineMove(g)
				return nil
			})
		}()
	}
	return nil
}

func passTurn(g *gocui.Gui, v *gocui.View) error {
	// Update prevBoard to current board for Ko rule
	if prevBoard == nil {
		prev := game.Board{}
		copy(prev[:], gui.Grid[:])
		prevBoard = &prev
	} else {
		copy(prevBoard[:], gui.Grid[:])
	}
	koPoint = nil // Passing clears Ko

	passCount++
	if passCount >= 2 {
		gameOver = true
		if v, err := g.View("prompt"); err == nil {
			v.Clear()
			blackScore, whiteScore := game.CalculateScore(gui.Grid)
			winner := "Black"
			if whiteScore > blackScore {
				winner = "White"
			} else if whiteScore == blackScore {
				winner = "Draw"
			}
			fmt.Fprintf(v, "Game Over! Black: %d, White: %d. Winner: %s", blackScore, whiteScore, winner)
		}
		return nil
	}

	// Show "Turn passed." message briefly
	if v, err := g.View("prompt"); err == nil {
		v.Clear()
		fmt.Fprint(v, "Turn passed.")
		go func() {
			time.Sleep(1 * time.Second)
			g.Update(func(_ *gocui.Gui) error {
				if v, err := g.View("prompt"); err == nil && v != nil {
					v.Clear()
					printMovePrompt(v)
				}
				return nil
			})
		}()
	}

	currentPlayer = 3 - currentPlayer
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func engineMove(g *gocui.Gui) {
	// Use the engine interface to get a move for White
	if selectedEngine == nil {
		return
	}
	move := selectedEngine.Move(gui.Grid, game.White, koPoint)
	if move != nil {
		// Do not move the cursor for the engine, just place the stone directly
		row, col := move.Row, move.Col
		gui.Grid[row][col] = game.White

		// Simulate captures and ko logic as in placeStone
		var nextBoard game.Board
		copy(nextBoard[:], gui.Grid[:])
		opp := game.Black
		captured := []game.Point{}
		for _, n := range game.Neighbors(game.Point{Row: row, Col: col}) {
			if nextBoard[n.Row][n.Col] == opp {
				group, libs := game.Group(nextBoard, n)
				if len(libs) == 0 {
					for stonePt := range group {
						nextBoard[stonePt.Row][stonePt.Col] = game.Empty
						captured = append(captured, stonePt)
					}
				}
			}
		}
		for _, pt := range captured {
			gui.Grid[pt.Row][pt.Col] = game.Empty
		}
		// Ko rule: set koPoint if exactly one stone was captured and the group size is 1
		if len(captured) == 1 {
			koPoint = &captured[0]
		} else {
			koPoint = nil
		}
		// Update prevBoard for next move (for Ko rule)
		if prevBoard == nil {
			prev := game.Board{}
			copy(prev[:], gui.Grid[:])
			prevBoard = &prev
		}
		copy(prevBoard[:], gui.Grid[:])
		passCount = 0
		currentPlayer = 1 // Switch back to player
	} else {
		_ = passTurn(g, nil)
	}
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	cfg, err := loadConfig("config.json")
	if err != nil {
		log.Panicln("Failed to load config:", err)
	}
	keybindings = cfg.Keybindings

	// selectedEngine = &engine.RandomEngine{}
	selectedEngine = engine.NewAlphaBetaEngine()
	engineEnabled = true // Enable engine by default

	defer g.Close()

	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	// Keybindings
	moveLeftKey := []rune(keybindings["moveLeft"])[0]
	moveRightKey := []rune(keybindings["moveRight"])[0]
	moveUpKey := []rune(keybindings["moveUp"])[0]
	moveDownKey := []rune(keybindings["moveDown"])[0]
	quitKey := []rune(keybindings["quit"])[0]
	placeStoneKey := keybindings["placeStone"]
	passTurnKey := keybindings["passTurn"]
	var placeKey rune
	if placeStoneKey == "space" {
		placeKey = ' '
	} else {
		placeKey = []rune(placeStoneKey)[0]
	}
	var passKey rune
	if passTurnKey == "space" {
		passKey = ' '
	} else {
		passKey = []rune(passTurnKey)[0]
	}

	// Lowercase: move regardless of occupation
	if err := g.SetKeybinding("", moveLeftKey, gocui.ModNone, moveCursor(0, -1, false)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", moveRightKey, gocui.ModNone, moveCursor(0, 1, false)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", moveUpKey, gocui.ModNone, moveCursor(-1, 0, false)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", moveDownKey, gocui.ModNone, moveCursor(1, 0, false)); err != nil {
		log.Panicln(err)
	}
	// Uppercase: jump over occupied intersections (Shift+key)
	if err := g.SetKeybinding("", rune(unicode.ToUpper(moveLeftKey)), gocui.ModNone, moveCursor(0, -1, true)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", rune(unicode.ToUpper(moveRightKey)), gocui.ModNone, moveCursor(0, 1, true)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", rune(unicode.ToUpper(moveUpKey)), gocui.ModNone, moveCursor(-1, 0, true)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", rune(unicode.ToUpper(moveDownKey)), gocui.ModNone, moveCursor(1, 0, true)); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", placeKey, gocui.ModNone, placeStone); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", passKey, gocui.ModNone, passTurn); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", quitKey, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	// Engine toggle: enable for second player
	enableEngineKey := []rune(keybindings["enableEngine"])[0]
	if err := g.SetKeybinding("", enableEngineKey, gocui.ModNone, toggleEngine); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func toggleEngine(g *gocui.Gui, v *gocui.View) error {
	engineEnabled = !engineEnabled
	if v, err := g.View("prompt"); err == nil && v != nil {
		v.Clear()
		printMovePrompt(v)
	}
	// If toggled on and it's engine's turn (player 2/White), make engine move
	if engineEnabled && !gameOver && currentPlayer == 2 {
		go func() {
			time.Sleep(300 * time.Millisecond)
			g.Update(func(g *gocui.Gui) error {
				engineMove(g)
				return nil
			})
		}()
	}
	return nil
}
