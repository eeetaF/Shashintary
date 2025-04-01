package commentary

import (
	"fmt"

	config_module "Shashintary/modules/config"
)

var funcMap = map[string]func(cfg *config_module.Config) string{
	"welcome_message": welcomeMessage,
}

func welcomeMessage(cfg *config_module.Config) string {
	return fmt.Sprintf("Today %s plays %s. This will be interesting!", cfg.PlayerWhite, cfg.PlayerBlack)
}

func CommentsHandler(cfg *config_module.Config, wantedMessage string) (string, error) {
	if messageFunc, exists := funcMap[wantedMessage]; exists {
		return messageFunc(cfg), nil
	}
	return "", fmt.Errorf("unknown message key")
}

func PromptInitialPosition(playerWhite, playerBlack, language string) string {
	return fmt.Sprintf("You are a professional chess commentator.\n\nGenerate a short and expressive introductory commentary for the start of a classical chess game.\n\nDetails:\n- Player with white pieces: %s\n- Player with black pieces: %s\n\nLanguage: %s\n\nGuidelines:\n- The tone should be engaging, lively, and human-like — as if a real chess streamer is talking.\n- Do not translate or modify player names.\n- Do not provide move suggestions or predictions yet — it's just the start of the game.\n- Limit the response to 1-3 short sentences (under 300 characters if possible).\n- Use chess terminology naturally, but avoid sounding robotic or scripted.\n", playerWhite, playerBlack, language)
}
