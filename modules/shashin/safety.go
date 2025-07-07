package shashin

import (
	"github.com/notnil/chess"
)

// -3 - step towards Petrosian
// -1 - little step towards Petrosian
// 0 - equal
// 1 - little step towards Tal
// 3 - step towards Tal
func getSafetyFactor(game *chess.Game) int8 {
	board := game.Position().Board()

	newFen, err := chess.FEN(switchTurn(game.FEN()))
	if err != nil {
		// we're under check most probably
		return -3
	}

	selfField, enemyField := getKingsRadius(board, game.Position().Turn())
	selfMoves, enemyMoves := game.ValidMoves(), chess.NewGame(newFen).ValidMoves()
	var totalAttackingMoves, diff int

	for _, move := range selfMoves {
		if _, ok := enemyField[move.S2()]; ok {
			totalAttackingMoves++
			diff++
		}
	}

	for _, move := range enemyMoves {
		if _, ok := selfField[move.S2()]; ok {
			totalAttackingMoves++
			diff--
		}
	}

	if diff == 0 {
		return 0
	}

	var power int8 = 3

	if float32(abs(diff))/float32(totalAttackingMoves) < 0.4 || totalAttackingMoves < 4 {
		power = 1
	}

	if diff > 0 {
		return power
	}

	return -power
}

func abs(val int) int {
	if val < 0 {
		return -val
	}
	return val
}

func getKingsRadius(board *chess.Board, selfColor chess.Color) (map[chess.Square]struct{}, map[chess.Square]struct{}) {
	var i int8
	var selfField, enemyField map[chess.Square]struct{}

	for i = 0; i < 64; i++ {
		sq := chess.Square(i)

		piece := board.Piece(sq)

		if piece.Type() == chess.King {
			field := generateFieldAround(sq)

			if piece.Color() == selfColor {
				selfField = field
			} else {
				enemyField = field
			}
		}
	}

	return selfField, enemyField
}

func generateFieldAround(sq chess.Square) map[chess.Square]struct{} {
	field := make(map[chess.Square]struct{}, 8)

	rank := int8(sq) / 8
	file := int8(sq) % 8

	directions := [8][2]int8{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for _, d := range directions {
		newRank := rank + d[0]
		newFile := file + d[1]

		if newRank >= 0 && newRank < 8 && newFile >= 0 && newFile < 8 {
			neighbor := chess.Square(newRank*8 + newFile)
			field[neighbor] = struct{}{}
		}
	}

	return field
}
