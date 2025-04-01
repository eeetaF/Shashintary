package shashins_parameters

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/notnil/chess"
)

// Evaluate position clearly according to Shashin's theory
func evaluatePosition(fen string) (float64, float64, float64, string) {
	game, err := chess.FEN(fen)
	if err != nil {
		log.Fatal(err)
	}
	currGame := chess.NewGame(game)

	material := CalculateMaterial(currGame.Position())
	mobility := CalculateMobility(currGame)
	safety := CalculateSafety(currGame.Position())

	return material, mobility, safety, DeterminePositionType(material, mobility, safety)
}

func EvaluatePositionShashin() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter FEN of the chess position:")
	fen, _ := reader.ReadString('\n')
	fen = strings.TrimSpace(fen)

	m, t, s, positionType := evaluatePosition(fen)

	fmt.Println("\nPosition Analysis based on Shashin's Theory:")
	fmt.Printf("Material Parameter (m): %.2f\n", m)
	fmt.Printf("Time Parameter (t): %.2f\n", t)
	fmt.Printf("Safety Parameter (s): %.2f\n", s)
	fmt.Printf("Position Type: %s\n", positionType)
}
