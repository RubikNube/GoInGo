package game

type FieldState uint8

const (
	Black FieldState = iota
	White
)

type GameState struct {
	LastMoveRow   uint8 // 0-8 for rows
	LastMoveCol   uint8 // 0-8 for columns
	LastMoveColor FieldState
	SpecialCase   uint8 // 0 = none, 1 = ko, 2 = seki, 3 = pass, etc.
}

// Encode encodes the game state as a single uint8 value.
// 0-80: black stones, 81-161: white stones, 162-255: special cases.
func (gs *GameState) Encode() uint8 {
	if gs.SpecialCase > 0 {
		return 161 + gs.SpecialCase
	}
	if gs.LastMoveColor == Black {
		return uint8(gs.LastMoveRow*9 + gs.LastMoveCol)
	} else if gs.LastMoveColor == White {
		return uint8(81 + gs.LastMoveRow*9 + gs.LastMoveCol)
	}
	return 255 // fallback for undefined state
}

// Decode decodes a uint8 value into the game state.
func (gs *GameState) Decode(val uint8) {
	switch {
	case val <= 80:
		gs.LastMoveColor = Black
		gs.LastMoveRow = val / 9
		gs.LastMoveCol = val % 9
		gs.SpecialCase = 0
	case val <= 161:
		gs.LastMoveColor = White
		adj := val - 81
		gs.LastMoveRow = adj / 9
		gs.LastMoveCol = adj % 9
		gs.SpecialCase = 0
	default:
		gs.SpecialCase = val - 161
		gs.LastMoveColor = 0
		gs.LastMoveRow = 0
		gs.LastMoveCol = 0
	}
}

