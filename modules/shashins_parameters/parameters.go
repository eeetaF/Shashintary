package shashins_parameters

import (
	"strings"

	"github.com/notnil/chess"
)

func DeterminePositionType(m, t, safety float64) string {
	score := map[string]int{"Tal": 0, "Capablanca": 0, "Petrosian": 0}

	if m > 1.0 {
		score["Tal"]++
	} else if m < 1.0 {
		score["Petrosian"]++
	} else {
		score["Capablanca"]++
	}

	if t >= 1.25 {
		score["Tal"]++
	} else if t <= 0.80 {
		score["Petrosian"]++
	} else {
		score["Capablanca"]++
	}

	if safety <= -2.0 {
		score["Tal"]++
	} else if safety >= 2.0 {
		score["Petrosian"]++
	} else {
		score["Capablanca"]++
	}

	maxScore := "Tal"
	for k, v := range score {
		if v > score[maxScore] {
			maxScore = k
		}
	}
	return maxScore
}

func CalculateMaterial(pos *chess.Position) float64 {
	pieceValues := map[chess.PieceType]float64{
		chess.Pawn:   1,
		chess.Knight: 3,
		chess.Bishop: 3,
		chess.Rook:   5,
		chess.Queen:  9,
	}
	balance := 0.0
	for _, piece := range pos.Board().SquareMap() {
		value := pieceValues[piece.Type()]
		if piece.Color() == pos.Turn() {
			balance += value
		} else {
			balance -= value
		}
	}
	return balance
}

func CalculateMobility(currGame *chess.Game) float64 {
	activePlayerMoves := float64(len(currGame.Position().ValidMoves()))
	var opponentPlayerMoves float64

	// copy the game with switching color to count opponent's activity
	parts := strings.Fields(currGame.FEN())
	if parts[1] == "w" {
		parts[1] = "b"
	} else {
		parts[1] = "w"
	}
	game, err := chess.FEN(strings.Join(parts, " "))
	if err == nil {
		newGame := chess.NewGame(game)
		opponentPlayerMoves = float64(len(newGame.ValidMoves()))
	}
	if opponentPlayerMoves == 0.0 {
		return activePlayerMoves
	}

	return activePlayerMoves / opponentPlayerMoves
}

func CalculateSafety(pos *chess.Position) float64 {
	color := pos.Turn()
	safetyScore := 0.0

	// 1. King Safety Evaluation
	kingSquare := findKingSquare(pos, color)
	if kingSquare != chess.NoSquare {
		safetyScore += evaluateKingSafety(pos, kingSquare)
	}

	// 2. Major Piece Safety Evaluation (Queens and Rooks)
	safetyScore += evaluateMajorPieceSafety(pos)

	return safetyScore
}

func CalculateCompactness(pos *chess.Position, color chess.Color) float64 {
	return 0.0
	//centerSquares := []chess.Square{chess.D4, chess.D5, chess.E4, chess.E5}
	//compactness := 0.0
	//for sq, piece := range pos.Board().SquareMap() {
	//	if piece.Color() == color {
	//		for _, center := range centerSquares {
	//			distance := chess.Distance(sq, center)
	//			compactness += (4.0 - float64(distance)) / 4.0
	//		}
	//	}
	//}
	//return compactness
}

// findKingSquare locates the king's position on the board.
func findKingSquare(pos *chess.Position, color chess.Color) chess.Square {
	for sq, piece := range pos.Board().SquareMap() {
		if piece.Type() == chess.King && piece.Color() == color {
			return sq
		}
	}
	return chess.NoSquare
}

// evaluateKingSafety calculates king's safety for current player.
// +1 per friendly piece around the king
// -2 per enemy piece around the king
// -1.5 per attacked square around the king
// -7.0 if king is in check
func evaluateKingSafety(pos *chess.Position, kingSquare chess.Square) float64 {
	safetyScore := 0.0
	board := pos.Board()
	color := pos.Turn()

	file := int(kingSquare.File())
	rank := int(kingSquare.Rank())

	directions := [][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for _, dir := range directions {
		f := file + dir[0]
		r := rank + dir[1]
		if f < 0 || f > 7 || r < 0 || r > 7 {
			continue
		}
		sq := chess.Square(r*8 + f)
		piece := board.Piece(sq)

		if piece != chess.NoPiece {
			if piece.Color() == color {
				safetyScore += 1.0
			} else {
				safetyScore -= 2.0
			}
		} else {
			if isSquareAttacked(pos, sq) {
				safetyScore -= 1.5
			}
		}
	}

	// If king is currently in check
	if isSquareAttacked(pos, kingSquare) {
		safetyScore -= 7.0
	}

	return safetyScore
}

// evaluateMajorPieceSafety calculates major pieces' safety for current player.
// -4.5 per attack on the queen
// -2.5 per attack on the rook
func evaluateMajorPieceSafety(pos *chess.Position) float64 {
	safetyScore := 0.0
	pieceValues := map[chess.PieceType]float64{
		chess.Queen: 9.0,
		chess.Rook:  5.0,
	}

	for sq, piece := range pos.Board().SquareMap() {
		if piece.Color() == pos.Turn() && (piece.Type() == chess.Queen || piece.Type() == chess.Rook) {
			if isSquareAttacked(pos, sq) {
				safetyScore -= pieceValues[piece.Type()] * 0.5 // Penalty for threatened major piece
			}
		}
	}

	return safetyScore
}

func isSquareAttacked(pos *chess.Position, sq chess.Square) bool {
	// Switch the FEN to give the move to the opposing side
	fenParts := strings.Fields(pos.String())
	if len(fenParts) < 2 {
		return false
	}
	if fenParts[1] == "w" {
		fenParts[1] = "b"
	} else {
		fenParts[1] = "w"
	}
	newFEN := strings.Join(fenParts, " ")
	game, err := chess.FEN(newFEN)
	if err != nil {
		return false
	}
	tempGame := chess.NewGame(game)

	for _, move := range tempGame.ValidMoves() {
		if move.S2() == sq {
			return true
		}
	}
	return false
}
