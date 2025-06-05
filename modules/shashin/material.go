package shashin

import (
	"github.com/notnil/chess"
)

// -2 - step towards Petrosian
// -1 - little step
// 0 - equal
// 1 - little step
// 2 - step towards Tal
func getMaterialFactor(pos *chess.Position) int8 {
	// index - piece type in notnil/chess, value - piece value
	pieceToMaterial := [7]int16{0, 0, 18, 10, 7, 6, 2}

	var materialDifference int16
	for _, piece := range pos.Board().SquareMap() {
		if piece.Color() == pos.Turn() {
			materialDifference += pieceToMaterial[piece.Type()]
		} else {
			materialDifference -= pieceToMaterial[piece.Type()]
		}
	}

	switch {
	case materialDifference > 2:
		return -2
	case materialDifference > 0:
		return -1
	case materialDifference < -2:
		return 2
	case materialDifference < 0:
		return 1
	default:
		return 0
	}
}
