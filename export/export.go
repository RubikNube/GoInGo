package main

/*
#include <stdint.h>
*/
import "C"
import (
	"sync"

	"github.com/RubikNube/GoInGo/pkg/compareengines"
	"github.com/RubikNube/GoInGo/pkg/engine"
	"github.com/RubikNube/GoInGo/pkg/game"
)

var (
	engineRegistry = struct {
		sync.Mutex
		nextID  uint64
		objects map[uint64]engine.Engine
	}{objects: make(map[uint64]engine.Engine)}
	boardRegistry = struct {
		sync.Mutex
		nextID  uint64
		objects map[uint64]*game.Board
	}{objects: make(map[uint64]*game.Board)}
)

//export NewAlphaBetaEngine
func NewAlphaBetaEngine() C.uint64_t {
	engineRegistry.Lock()
	defer engineRegistry.Unlock()
	e := engine.NewAlphaBetaEngine()
	id := engineRegistry.nextID
	engineRegistry.nextID++
	engineRegistry.objects[id] = e
	return C.uint64_t(id)
}

//export NewRandomEngine
func NewRandomEngine() C.uint64_t {
	engineRegistry.Lock()
	defer engineRegistry.Unlock()
	e := engine.NewRandomEngine()
	id := engineRegistry.nextID
	engineRegistry.nextID++
	engineRegistry.objects[id] = e
	return C.uint64_t(id)
}

//export NewBoard
func NewBoard(size C.int) C.uint64_t {
	boardRegistry.Lock()
	defer boardRegistry.Unlock()
	b := game.NewBoard()
	id := boardRegistry.nextID
	boardRegistry.nextID++
	boardRegistry.objects[id] = &b
	return C.uint64_t(id)
}

//export CompareEngines
func CompareEngines(engineAID, engineBID, boardID C.uint64_t, firstPlayer C.int, maxMoves C.int) C.int {
	engineRegistry.Lock()
	engineA := engineRegistry.objects[uint64(engineAID)]
	engineB := engineRegistry.objects[uint64(engineBID)]
	engineRegistry.Unlock()
	boardRegistry.Lock()
	board := boardRegistry.objects[uint64(boardID)]
	boardRegistry.Unlock()
	if engineA == nil || engineB == nil || board == nil {
		return -1 // error code
	}
	result := compareengines.CompareEngines(engineA, engineB, *board, game.FieldState(firstPlayer), int(maxMoves))
	if result == 0 {
		return 0 // draw
	} else if result > 0 {
		return 1 // engineA wins
	} else {
		return 2 // engineB wins
	}
}

func main() {}
