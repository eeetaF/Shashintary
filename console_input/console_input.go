package console_input

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func RunInput(host, port string) error {
	address := fmt.Sprintf("%s:%s", host, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("couldn't connect to %s: %v", address, err)
	}
	defer conn.Close()
	fmt.Printf("Connected to %s.\nYou are now able to send data (FEN, moves, etc.)\nUse 'exit' to safe exit\n", address)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text()) + "\n"

		if text == "\n" {
			continue
		}

		_, err = fmt.Fprintf(conn, text)
		if err != nil {
			return fmt.Errorf("couldn't send data: %v\n", err)
		}
		if text == "exit\n" {
			return nil
		}
	}

	return nil
}
