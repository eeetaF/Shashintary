package shashin

import (
	"github.com/notnil/chess"
)

// -2 - Petrosian
// -1 - Capablanca-Petrosian
// 0 - equal
// 1 - Capablanca-Tal
// 2 - Tal
func getTimeFactor(pos *chess.Position) int8 {
	numValidMoves := len(pos.ValidMoves())

	definedFen, err := chess.FEN(switchTurn(pos.String()))
	if err != nil {
		// player's most probably in check
		return -2
	}

	newGame := chess.NewGame(definedFen)

	numValidMovesOther := len(newGame.Position().ValidMoves())

	ratio := float32(numValidMoves) / float32(numValidMovesOther)
	if ratio < 0.8 {
		return -2
	}
	if ratio < 1 {
		return -1
	}
	if ratio == 1 {
		return 0
	}
	if ratio < 1.25 {
		return 1
	}
	return 2
}
