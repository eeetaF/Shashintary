package program_interface

import (
	"fmt"

	"github.com/notnil/chess"

	"Shashintary/modules"
	"Shashintary/modules/commentary"
	config_module "Shashintary/modules/config"
	"Shashintary/modules/message"
	"Shashintary/modules/shashin"
)

func HandleGame(cfg *config_module.Config, inputChannel <-chan string, validInputMovesChannel chan<- modules.Input,
	calculatedMovesChannel <-chan []modules.CalculatedMove, outputChannel chan<- []*message.OutputMessage) {

	go handleIncomingCalculatedMoves(cfg, calculatedMovesChannel)
	sendBoard := *cfg.SendBoard
	outputChannel <- stringToSliceOfOutputMessages("@", false)
	for {
		outputChannel <- stringToSliceOfOutputMessages("Ready to comment game. Provide the FEN or send \"default\" to start from the initial chess position\n", false)
		fen, ok := <-inputChannel
		if !ok {
			return
		}
		game, err := initGameByFen(fen)
		if err != nil {
			outputChannel <- stringToSliceOfOutputMessages("@Couldn't initialize the game with this FEN. Try again.\n", false)
			continue
		}
		validInputMovesChannel <- modules.Input{Move: game.FEN(), IsFEN: true}
		if sendBoard {
			sendBoardToChannel(outputChannel, game)
		}
		sendInitialPosition(outputChannel, cfg)

		for game.Outcome() == chess.NoOutcome {
			move := <-inputChannel
			err = game.MoveStr(move)
			if err != nil {
				fmt.Printf("received invalid move: %s\n", move)
				continue
			}
			move = game.Moves()[len(game.Moves())-1].String() // always in UCI notation

			validInputMovesChannel <- modules.Input{Move: move, IsFEN: false}

			sendBoardToChannel(outputChannel, game)
			color := game.Position().Turn().String()
			if color == "w" {
				color = "black"
			} else {
				color = "white"
			}
			shashType := shashin.GetPositionType(game)
			outputChannel <- stringToSliceOfOutputMessages(getPromptResult(move, color, game.FEN(), shashType), false)
		}
		outputChannel <- stringToSliceOfOutputMessages("Game is finished. Write anything to continue...", false)
		<-inputChannel
	}
}

func stringToSliceOfOutputMessages(s string, isBoard bool) []*message.OutputMessage {
	return []*message.OutputMessage{{Value: s, IsBoard: isBoard}}
}

func initGameByFen(fen string) (*chess.Game, error) {
	if fen == "default" || fen == "" {
		fen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	}
	chessFen, err := chess.FEN(fen)
	if err != nil {
		return nil, err
	}
	return chess.NewGame(chessFen), nil
}

func sendBoardToChannel(outputChannel chan<- []*message.OutputMessage, game *chess.Game) {
	outputChannel <- stringToSliceOfOutputMessages("@"+game.Position().Board().Draw(), true)
}

func sendInitialPosition(outputChannel chan<- []*message.OutputMessage, cfg *config_module.Config) {
	outputChannel <- stringToSliceOfOutputMessages(commentary.SendPrompt(&commentary.PromptOpts{
		ReadyPrompt: commentary.PromptInitialPosition(cfg.PlayerWhite, cfg.PlayerBlack, cfg.Language),
	}), false)
}
