package config_module

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Engine      string // "shashchess", "stockfish"
	DebugMode   *bool  // false, true, nil
	Language    string // "russian", "english", "turkish"
	PlayerWhite string // "Alex", "Hikaru Nakamura", "White"
	PlayerBlack string // "Ben", "Magnus Carlsen", "Black"

	SendBoard *bool // true, false, nil

	SelfHost string // "localhost"
	SelfPort string // "53002"

	DisplayHost string // "localhost"
	DisplayPort string // "53003"
}

func (cfg *Config) LoadConfig() error {
	file, err := os.Open("config.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Expect key: value format
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "engine":
			cfg.Engine = strings.ToLower(value)
		case "debug":
			dm := strings.ToLower(value) == "true"
			cfg.DebugMode = &dm
		case "language":
			cfg.Language = strings.ToLower(value)
		case "player_white":
			cfg.PlayerWhite = value
		case "player_black":
			cfg.PlayerBlack = value
		case "output_host":
			cfg.SelfHost = value
		case "output_port":
			cfg.SelfPort = value
		case "display_host":
			cfg.DisplayHost = value
		case "display_port":
			cfg.DisplayPort = value
		case "send_board":
			dm := strings.ToLower(value) == "true"
			cfg.SendBoard = &dm
		}
	}
	if cfg.PlayerWhite == "" || cfg.PlayerBlack == "" {
		cfg.PlayerWhite = ""
		cfg.PlayerBlack = ""
	}

	return nil
}

func (cfg *Config) CompleteEmptyConfigFields() {
	reader := bufio.NewReader(os.Stdin)

	// DebugMode
	if cfg.DebugMode == nil {
		fmt.Print("Use debug mode (press Enter to default: no) y / n: ")
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))
		dm := input == "y" || input == "yes"
		cfg.DebugMode = &dm
	}

	// Host
	if cfg.SelfHost == "" {
		fmt.Print("Host to run on (press Enter to default: localhost): ")
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))
		// todo add check if host is valid
		if input == "" {
			input = "localhost"
		}
		cfg.SelfHost = input
	}

	// Port
	if cfg.SelfPort == "" {
		fmt.Print("Port to run on (press Enter to default: 53002): ")
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))
		// todo add check if port is valid
		if input == "" {
			input = "53002"
		}
		cfg.SelfPort = input
	}

	// DisplayHost
	if cfg.DisplayHost == "" {
		fmt.Print("Host to broadcast to (press Enter to default: localhost): ")
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))
		// todo add check if host is valid
		if input == "" {
			input = "localhost"
		}
		cfg.DisplayHost = input
	}

	// DisplayPort
	if cfg.DisplayPort == "" {
		fmt.Print("Port to broadcast to (press Enter to default: 53003): ")
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))
		// todo add check if port is valid
		if input == "" {
			input = "53003"
		}
		cfg.DisplayPort = input
	}

	// SendBoard
	if cfg.SendBoard == nil {
		fmt.Print("Broadcast board (press Enter to default: yes) y / n: ")
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))
		dm := input == "y" || input == "yes" || input == ""
		cfg.SendBoard = &dm
	}

	// Engine
	if cfg.Engine == "" {
		fmt.Print("Engine for chess calculations (press Enter to default: ShashChess): ")
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))
		if input == "" {
			input = "shashchess"
		}
		cfg.Engine = input
	}

	// Language
	if cfg.Language == "" {
		fmt.Print("Commentary language (press Enter to default: english): ")
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))
		if input == "" {
			input = "english"
		}
		cfg.Language = input
	}

	// Players
	if cfg.PlayerWhite == "" {
		fmt.Print("Player playing white (press Enter to default: White): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			input = "White"
		}
		cfg.PlayerWhite = input
		fmt.Print("Player playing black (press Enter to default: Black): ")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			input = "Black"
		}
		cfg.PlayerBlack = input
	}
}
