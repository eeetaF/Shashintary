package config

import (
	"bufio"
	"log"
	"os"
	"strings"

	config_module "Shashintary/modules/config"
)

// LoadConfig loads config only for output mode.
func LoadConfig(loadFromConfig bool) *config_module.Config {
	cfg := &config_module.Config{}

	if loadFromConfig {
		if err := cfg.LoadConfig(); err != nil {
			log.Fatalf("load config: %v", err)
		}
	}
	cfg.CompleteEmptyConfigFields()

	return cfg
}

func LoadHostPort(isOutput bool) (host, port string) {
	host = "localhost"
	port = "53003"
	lookingForHost := "display_host"
	lookingForPort := "display_port"
	if isOutput {
		port = "53002"
		lookingForHost = "output_host"
		lookingForPort = "output_port"
	}

	file, err := os.Open("config.txt")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case lookingForHost:
			if value != "" {
				host = value
			}
		case lookingForPort:
			if value != "" {
				port = value
			}
		}
	}

	return
}
