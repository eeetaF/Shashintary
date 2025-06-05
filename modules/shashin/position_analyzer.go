package shashin

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/notnil/chess"
)

// GetPositionType returns type of position
// 2 - Tal
// 1 - Capablanca-Tal
// 0 - Capablanca
// -1 - Capablanca-Petrosian
// -2 - Petrosian
func GetPositionType(game *chess.Game) int8 {
	fenParts := strings.Split(game.FEN(), " ")
	movesMade, _ := strconv.Atoi(fenParts[5])
	if movesMade < 7 {
		return 0
	}
	// res raises -> steps towards Tal
	// res decreases -> steps towards Petrosian
	// res around 0 -> Capablanca
	var res int8

	matFactor := getMaterialFactor(game.Position())
	res += matFactor
	fmt.Printf("materialFactor: %d\n", matFactor)

	timeFactor := getTimeFactor(game.Position())
	res += timeFactor
	fmt.Printf("timeFactor: %d\n", timeFactor)

	safetyFactor := getSafetyFactor(game)
	res += safetyFactor
	fmt.Printf("safetyFactor: %d\n", safetyFactor)

	if res <= 1 {
		compactness := getCompactnessFactor(game.Position())
		res += compactness
		fmt.Printf("compactnessFactor: %d\n", compactness)

		expansion := getExpansionFactor(game.Position())
		res += expansion
		fmt.Printf("expansionFactor: %d\n", expansion)
	}

	fmt.Println("final:", res)

	switch {
	case res > 1:
		return 2
	case res == 1:
		return 1
	case res == -1:
		return -1
	case res < -1:
		return -2
	}

	return 0
}
