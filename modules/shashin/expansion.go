package shashin

import (
	"github.com/notnil/chess"
)

// -1 - step towards Petrosian
// 0 - equal
// 1 - step towards Tal
func getExpansionFactor(pos *chess.Position) int8 {
	var whiteRankSum, blackRankSum, whiteNumPieces, blackNumPieces int

	for sq, pc := range pos.Board().SquareMap() {
		if pc.Color() == chess.White {
			whiteNumPieces++
			whiteRankSum += int(sq.Rank()) + 1
		} else {
			blackNumPieces++
			blackRankSum += 8 - int(sq.Rank())
		}
	}

	whiteExpansion := float32(whiteRankSum) / float32(whiteNumPieces)
	blackExpansion := float32(blackRankSum) / float32(blackNumPieces)

	var res int8 = 1.0
	if pos.Turn() == chess.Black {
		res = -1
	}

	switch {
	case whiteExpansion < blackExpansion:
		return res
	case whiteExpansion > blackExpansion:
		return -res
	}

	return 0
}
