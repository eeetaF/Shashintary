package commentary

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type PromptRequest struct {
	Prompt      string  `json:"prompt"`
	MinTokens   int     `json:"min_tokens"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float64 `json:"temperature"`
}

type PromptOpts struct {
	ReadyPrompt  string
	BestMove     string // "Nf3"
	MoveMade     string // "d4"
	EvalBefore   string // "+0.3"
	EvalAfter    string // "+0.1"
	SideMoved    string // "white" / "black"
	Continuation string // "Nf3 Nc6 Bb5"
	FEN          string // current position
	PlayerBlack  string // Hikaru Nakamura
	PlayerWhite  string // Magnus Carlsen
	Language     string // Russian
}

type PromptResponse struct {
	Result string `json:"result"`
}

func generatePromptText(r *PromptOpts) string {
	return fmt.Sprintf(
		`You are a chess commentator. Comment a move made by (%s), as a real commentator would. Answer strictly in the following language: %s 
- Move: %s
- Best move in position: %s
- Evaluation before this move: %s
- Evaluation after this move: %s
- Best continuation: %s
- FEN of current position: %s
- Player playing white: %s
- Player playing black: %s
- 

Your comment should feel real, creative and have 400-600 symbols in it. Dont mention that its engine response. Answer as you were a professional chess commentator. Never mention exact evaluation. If evaluation changes are minor: consider the move good enough. If evaluation changes are big (2 or more points), comment negatively on this move and explain why it's bad. If the move made and best move are the same, consider this move the best in the position and explain why it's so good. If you see "MX" in evaluation, it means it's mate in X. If it's "-MX", it means it's mate in X against us. `,
		r.SideMoved,
		r.Language,
		r.MoveMade,
		r.BestMove,
		r.EvalBefore,
		r.EvalAfter,
		r.Continuation,
		r.FEN,
		r.PlayerWhite,
		r.PlayerBlack,
	)
}

func SendPrompt(opts *PromptOpts) string {
	request := &PromptRequest{}
	if opts.ReadyPrompt == "" {
		request.Prompt = generatePromptText(opts)
	} else {
		request.Prompt = opts.ReadyPrompt
	}
	fmt.Println("prompt: ")
	fmt.Println(request.Prompt)
	fmt.Println()
	request.MinTokens = 200
	request.MaxTokens = 400
	request.Temperature = 0.7

	payload, _ := json.Marshal(*request)

	resp, err := http.Post("http://localhost:53004/generate", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var result PromptResponse
	json.Unmarshal(body, &result)

	return result.Result + "\n"
}
