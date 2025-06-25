package main

import (
	"fmt"
	"log"

	"encoding/json"
	"os"
	"github.com/RubikNube/GoInGo/cmd/game"
	"github.com/jroimartin/gocui"
)

type Config struct {
	Keybindings map[string]string `json:"keybindings"`
}

var (
	cursorRow, cursorCol int
	gui                  game.Gui
	keybindings          map[string]string
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
		fmt.Fprintf(v, "Move (%s/%s/%s/%s), %s to place stone, %s to quit",keybindings["moveLeft"], keybindings["moveDown"], keybindings["moveUp"], keybindings["moveRight"], keybindings["placeStone"], keybindings["quit"])
	}
	// Always redraw board
	if v, err := g.View("board"); err == nil {
		v.Clear()
		gui.DrawGridToWriter(v, cursorRow, cursorCol)
	}
	return nil
}

func moveCursor(dRow, dCol int) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		cursorRow += dRow
		cursorCol += dCol
		if cursorRow < 0 {
			cursorRow = 0
		}
		if cursorRow > 8 {
			cursorRow = 8
		}
		if cursorCol < 0 {
			cursorCol = 0
		}
		if cursorCol > 8 {
			cursorCol = 8
		}
		return nil
	}
}

func placeStone(g *gocui.Gui, v *gocui.View) error {
	if gui.Grid[cursorRow][cursorCol] == game.Empty {
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
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	cfg, err := loadConfig("config.json")
	if err != nil {
		log.Panicln("Failed to load config:", err)
	}
	keybindings = cfg.Keybindings

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
	var placeKey rune
	if placeStoneKey == "space" {
		placeKey = ' '
	} else {
		placeKey = []rune(placeStoneKey)[0]
	}

	if err := g.SetKeybinding("", moveLeftKey, gocui.ModNone, moveCursor(0, -1)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", moveRightKey, gocui.ModNone, moveCursor(0, 1)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", moveUpKey, gocui.ModNone, moveCursor(-1, 0)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", moveDownKey, gocui.ModNone, moveCursor(1, 0)); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", placeKey, gocui.ModNone, placeStone); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", quitKey, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
