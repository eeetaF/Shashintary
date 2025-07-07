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
	ReadyPrompt        string
	BestMove           string // "Nf3"
	MoveMade           string // "d4"
	PieceMakingTheMove string // "Queen"
	EvalTrend          string // slightly improved
	SideMoved          string // "white" / "black"
	Continuation       string // "Nf3 Nc6 Bb5"
	FEN                string // current position
	PlayerBlack        string // Hikaru Nakamura
	PlayerWhite        string // Magnus Carlsen
	Language           string // Russian
	Shashin            int8   // -2 Petrosian, -1 CP, 0 Capablanca, 1 TC, 2 Tal
}

type PromptResponse struct {
	Result string `json:"result"`
}

func generatePromptText(r *PromptOpts) string {
	return fmt.Sprintf(
		`You are a professional chess commentator.

Current **position style**: %s   ← use this to set your narrative tone.

Comment on %s’s last move as if broadcasting live.
Answer strictly in: %s
Length: **exactly 500–600 symbols** (not words). 
Never reveal engine scores; if "MX" or "-MX" appears, only hint at mate.

Context
• Move played: %s   (piece: %s)
• Best engine move: %s
• Eval trend: %s
• Best continuation: %s
• Players: White – %s, Black – %s

Guidelines
1. If the played move equals the best move, praise it and explain *why* it fits the %s style.  
2. If it differs, compare the ideas: why is it weaker/stronger *in this style*.  
3. Weave one colourful adjective from the mood list %v but avoid cliches. Use analogues from the given language.
4. Conclude with a forward-looking hook (“Let’s see if…”, “Will (Black/White) seize the initiative?”).

Dont mention that its engine response.
Now deliver the commentary.`,
		toShashin(r.Shashin),
		r.SideMoved,
		r.Language,
		r.MoveMade,
		r.PieceMakingTheMove,
		r.BestMove,
		r.EvalTrend,
		r.Continuation,
		r.PlayerWhite,
		r.PlayerBlack,
		toShashin(r.Shashin),
		getAdjectives(r.Shashin),
	)
}

func getAdjectives(shashin int8) []string {
	switch shashin {
	case -2:
		return []string{"prophylactic", "fortress-building", "surgical"}
	case -1:
		return []string{"resourceful", "stubborn", "consolidating"}
	case 1:
		return []string{"combustible", "on the brink", "momentum"}
	case 2:
		return []string{"blistering", "sacrificial", "initiative-driven"}
	}
	return []string{"harmonious", "prophylactic", "squeeze-type"}
}

func toShashin(shashin int8) string {
	switch shashin {
	case -2:
		return "Deep-Strategic-Defense"
	case -1:
		return "Sturdy-Defensive"
	case 1:
		return "Enterprising-Attack"
	case 2:
		return "Sharp-Attacking"
	}
	return "Balanced-Positional"
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
