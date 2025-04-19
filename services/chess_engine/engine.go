package chess_engine

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// getChessEngine do once
func getChessEngine(engineName string) (*bufio.Scanner, io.WriteCloser, error) {
	foundEngine, err := findEngine(engineName)
	if err != nil {
		return nil, nil, err
	}

	engine := exec.Command("./" + foundEngine)
	stdin, err := engine.StdinPipe()
	if err != nil {
		return nil, nil, err
	}
	stdout, err := engine.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	scanner := bufio.NewScanner(stdout)

	if err = engine.Start(); err != nil {
		return nil, nil, err
	}

	fmt.Fprintln(stdin, "uci")
	for scanner.Scan() {
		if scanner.Text() == "uciok" {
			break
		}
	}
	fmt.Fprintln(stdin, "isready")
	for scanner.Scan() {
		if scanner.Text() == "readyok" {
			break
		}
	}

	return scanner, stdin, nil
}

func findEngine(input string) (string, error) {
	var result string

	err := filepath.Walk("./engines/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		lowerName := strings.ToLower(info.Name())
		lowerQuery := strings.ToLower(input)

		if strings.Contains(lowerName, lowerQuery) {
			result = path
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return "", err
	}
	if result == "" {
		return "", fmt.Errorf("chess engine not found. Check engines directory for '%s'", input)
	}

	return result, nil
}
