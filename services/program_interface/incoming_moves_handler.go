package program_interface

import (
	"context"
	"strings"

	"Shashintary/modules"
	"Shashintary/modules/commentary"
	config_module "Shashintary/modules/config"
)

type RequestForMove struct {
	Move        string // e2e4
	Side        string // white
	FEN         string // ....
	ShashinType int8   // -2 for Petrosian, -1 CP, 0 Capablanca, 1 CT, 2 Tal
}
type preCalculated struct {
	bestMove     modules.CalculatedMove
	moves        map[string]modules.CalculatedMove
	requestChan  chan RequestForMove // e2e4
	responseChan chan string         // Great move!..
	sideToMove   string

	cfg *config_module.Config
}

type readyPrompt struct {
	move   string
	prompt string
}

var pc preCalculated

func handleIncomingCalculatedMoves(cfg *config_module.Config, incomingMoves <-chan []modules.CalculatedMove) {
	pc.cfg = cfg
	pc.requestChan = make(chan RequestForMove)
	pc.responseChan = make(chan string)
	pc.sideToMove = "white" // TODO fix this for game that start with black's move

	for {
		calculatedMoves := <-incomingMoves

		pc.moves = make(map[string]modules.CalculatedMove, len(calculatedMoves))

		for i := 0; i < len(calculatedMoves); i++ {
			move := calculatedMoves[i].Move
			pc.moves[move] = calculatedMoves[i]
		}
		pc.bestMove = calculatedMoves[0]

		select {
		case move := <-pc.requestChan: // request already present
			if move.Side == "black" {
				pc.sideToMove = "white"
			} else {
				pc.sideToMove = "black"
			}
			resp := sendSinglePromptRequest(move.Move, move.Side, move.FEN)
			pc.responseChan <- resp
			continue

		default: // request not received yet, calculate for every move
			ctx, cancelFn := context.WithCancel(context.Background())
			promptsChan := make(chan readyPrompt)
			sendPromptForAllMoves(ctx, promptsChan, pc.sideToMove)

			move := <-pc.requestChan
			if move.Side == "black" {
				pc.sideToMove = "white"
			} else {
				pc.sideToMove = "black"
			}

			if _, ok := pc.moves[move.Move]; !ok { // we didn't face this move
				cancelFn()
				pc.responseChan <- sendSinglePromptRequest(move.Move, move.Side, move.FEN)
			}

			// we did face this move and the result is being calculated
			for {
				currRes := <-promptsChan
				if currRes.move == move.Move {
					cancelFn()
					pc.responseChan <- currRes.prompt
					break
				}
			}
			continue
		}
	}
}

func sendPromptForAllMoves(ctx context.Context, prompts chan<- readyPrompt, sideMoved string) {
	for _, move := range pc.moves {
		go sendPromptRequest(ctx, prompts, move.Move, sideMoved, "")
	}
}

func sendSinglePromptRequest(moveUCI, sideMoved, fen string) string {
	return commentary.SendPrompt(generateOpts(moveUCI, sideMoved, fen))
}

func sendPromptRequest(ctx context.Context, promptsChan chan<- readyPrompt, moveUCI, sideMoved, fen string) {
	result := commentary.SendPrompt(generateOpts(moveUCI, sideMoved, fen))
	select {
	case promptsChan <- readyPrompt{
		move:   moveUCI,
		prompt: result,
	}:
		return
	case <-ctx.Done():
		return
	}
}

func generateOpts(moveUCI, sideMoved, fen string) *commentary.PromptOpts {
	opts := &commentary.PromptOpts{
		BestMove:    pc.bestMove.Move,
		MoveMade:    moveUCI,
		EvalBefore:  pc.bestMove.ScoreInCP,
		SideMoved:   sideMoved,
		FEN:         fen,
		PlayerBlack: pc.cfg.PlayerBlack,
		PlayerWhite: pc.cfg.PlayerWhite,
		Language:    pc.cfg.Language,
	}
	if moveInfo, ok := pc.moves[moveUCI]; ok {
		opts.EvalAfter = moveInfo.ScoreInCP
		opts.Continuation = strings.Join(moveInfo.ContinuationMoves[1:], " ")
	} else {
		opts.EvalAfter = "much less that a move before"
	}
	return opts
}

func getPromptResult(moveUCI, side, fen string, shashType int8) string {
	pc.requestChan <- RequestForMove{Move: moveUCI, Side: side, FEN: fen, ShashinType: shashType}
	resp := <-pc.responseChan
	return resp
}
