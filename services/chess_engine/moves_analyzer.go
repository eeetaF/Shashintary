package chess_engine

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"Shashintary/modules"
	config_module "Shashintary/modules/config"
)

type moveEval struct {
	PV    []string
	CP    string
	MPV   int
	Score int
}

type scoredMove struct {
	Move string
	Eval moveEval
}

func RunMovesAnalyzer(cfg *config_module.Config, inputChannel <-chan modules.Input, outputChannel chan<- []modules.CalculatedMove, maxBestMoves, maxContinuationMoves int) {
	if maxBestMoves <= 0 {
		maxBestMoves = 256
	} else {
		maxBestMoves = min(256, maxBestMoves)
	}
	moves := make([]string, 0, 100)

	scanner, stdin, err := getChessEngine(cfg.Engine)
	if err != nil {
		log.Fatalf("Moves Analyzer: Cannot initialize chess engine: %s", err)
	}

	var cancelCurrent context.CancelFunc
	var currentEval map[int]moveEval
	var currentFEN string
	dataRead := make(chan struct{})
	previousFinished := make(chan struct{})

	for {
		input := <-inputChannel

		if cancelCurrent != nil { // only enters here when it's not the first evaluation
			cancelCurrent()
			fmt.Fprintln(stdin, "stop")
			<-previousFinished // wait for previous eval to finish
		}

		if input.IsFEN {
			fmt.Fprintln(stdin, "ucinewgame")
			fmt.Fprintf(stdin, "position fen %s\n", input.Move)
			currentFEN = input.Move
			moves = moves[:0]
			fmt.Printf("Moves Analyzer: Started evaluation for new position. FEN: %s\n", input.Move)
		} else {
			moves = append(moves, input.Move)
			if currentFEN != "" {
				fmt.Fprintf(stdin, "position fen %s moves %s\n", currentFEN, strings.Join(moves, " "))
			} else {
				fmt.Fprintf(stdin, "position startpos moves %s\n", strings.Join(moves, " "))
			}
			fmt.Printf("Moves Analyzer: Started evaluation for new position. Move: %s\n", input.Move)
		}
		fmt.Fprintln(stdin, "setoption name MultiPV value "+strconv.Itoa(maxBestMoves))
		fmt.Fprintln(stdin, "go infinite")

		ctx, cancelFn := context.WithCancel(context.Background())
		cancelCurrent = cancelFn

		currentEval = make(map[int]moveEval, maxBestMoves)

		go func(localCtx context.Context) {
			for scanner.Scan() {
				select {
				case <-localCtx.Done():
					dataRead <- struct{}{} // signal done scanning
					return
				default:
				}
				line := scanner.Text()
				if eval, ok := parseEvalLine(line); ok {
					if eval.MPV >= 1 && eval.MPV <= maxBestMoves && len(eval.PV) > 0 {
						currentEval[eval.MPV] = eval
					}
				}
			}
			dataRead <- struct{}{} // signal done scanning
		}(ctx)

		go func(localCtx context.Context) {
			// wait for cancel by new move arrival or a minute of calculation
			select {
			case <-localCtx.Done():
			case <-time.After(time.Minute):
				cancelCurrent()
				fmt.Fprintln(stdin, "stop")
			}

			<-dataRead // wait for upper goroutine to collect data

			if len(currentEval) > 0 {
				fmt.Println("Moves Analyzer: Evaluation after initial delay:")
				var sorted []scoredMove
				for _, eval := range currentEval {
					sorted = append(sorted, scoredMove{Move: eval.PV[0], Eval: eval})
				}
				sort.Slice(sorted, func(i, j int) bool {
					return sorted[i].Eval.Score > sorted[j].Eval.Score
				})
				output := make([]modules.CalculatedMove, 0, len(sorted))
				for _, sm := range sorted {
					pv := sm.Eval.PV[:min(maxContinuationMoves, len(sm.Eval.PV))]
					fmt.Printf("MA: %q : %s | %s\n", sm.Move, sm.Eval.CP, strings.Join(pv, " "))
					output = append(output, modules.CalculatedMove{
						Move:              sm.Move,
						ScoreInCP:         sm.Eval.CP,
						ContinuationMoves: pv,
					})
				}
				outputChannel <- output
			} else {
				log.Fatalf("Moves Analyzer: No evaluation received after delay")
			}
			previousFinished <- struct{}{} // signal of finish to run a new evaluation
		}(ctx)

		time.Sleep(2 * time.Second) // Minimum time for engine to calculate at least something...
	}
}

func parseEvalLine(line string) (moveEval, bool) {
	if !strings.Contains(line, " pv ") || !strings.Contains(line, "score ") {
		return moveEval{}, false
	}

	var eval moveEval
	fields := strings.Fields(line)
	scoreIndex, pvIndex := -1, -1
	for i, f := range fields {
		if f == "multipv" && i+1 < len(fields) {
			mpv, err := strconv.Atoi(fields[i+1])
			if err == nil {
				eval.MPV = mpv
			}
		}
		if f == "score" && i+2 < len(fields) {
			scoreIndex = i + 1
		}
		if f == "pv" && i+1 < len(fields) {
			pvIndex = i + 1
			break
		}
	}

	if scoreIndex == -1 || pvIndex == -1 {
		return moveEval{}, false
	}

	scoreType := fields[scoreIndex]
	scoreValue := fields[scoreIndex+1]
	switch scoreType {
	case "cp":
		val, err := strconv.Atoi(scoreValue)
		if err != nil {
			return moveEval{}, false
		}
		eval.CP = fmt.Sprintf("%.2f", float64(val)/100.0)
		eval.Score = val
	case "mate":
		val, err := strconv.Atoi(scoreValue)
		if err != nil {
			return moveEval{}, false
		}
		if val > 0 {
			eval.Score = 100000 - val
			eval.CP = fmt.Sprintf("M%d", val)
		} else {
			eval.Score = -100000 - val
			eval.CP = fmt.Sprintf("-M%d", -val)
		}
	default:
		return moveEval{}, false
	}

	eval.PV = fields[pvIndex:]
	if eval.MPV == 0 {
		eval.MPV = 1
	}
	return eval, true
}
