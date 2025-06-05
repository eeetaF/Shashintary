package shashin

import (
	"strings"
)

func switchTurn(fen string) string {
	parts := strings.Split(fen, " ")
	if parts[1] == "w" {
		parts[1] = "b"
	} else if parts[1] == "b" {
		parts[1] = "w"
	}

	return strings.Join(parts, " ")
}
