package program_interface

import (
	"time"

	"github.com/notnil/chess"

	"Shashintary/modules/commentary"
	config_module "Shashintary/modules/config"
	"Shashintary/modules/message"
)

func HandleGame(cfg *config_module.Config, inputChannel <-chan string, outputChannel chan<- []*message.OutputMessage) {
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
		if sendBoard {
			sendBoardToChannel(outputChannel, game)
		}
		sendInitialPosition(outputChannel, cfg)
		//outputChannel <- stringToSliceOfOutputMessages(commentary.SendPrompt(&commentary.PromptRequest{
		//			BestMove:     "it's a starting position, any book move is good",
		//			MoveMade:     "none yet",
		//			EvalBefore:   "0",
		//			EvalAfter:    "0",
		//			SideMoved:    "none yet",
		//			Continuation: "any book move",
		//			FEN:          "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		//			PlayerBlack:  cfg.PlayerBlack,
		//			PlayerWhite:  cfg.PlayerWhite,
		//			Language:     cfg.Language,
		//		}), false)
		//}
		time.Sleep(time.Hour)
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
