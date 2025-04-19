package program_interface

import (
	"bufio"
	"fmt"
	"net"

	"Shashintary/modules"
	config_module "Shashintary/modules/config"
	"Shashintary/modules/message"
	"Shashintary/services/broadcast"
	"Shashintary/services/chess_engine"
)

func HandleProgram(cfg *config_module.Config) error {
	address := fmt.Sprintf("%s:%s", cfg.SelfHost, cfg.SelfPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("couldn't run server on %s: %v", address, err)
	}
	defer listener.Close()

	inputChannel := make(chan string, 50)
	validInputMovesChannel := make(chan modules.Input, 50)
	calculatedMovesChannel := make(chan []modules.CalculatedMove, 1)
	outputChannel := make(chan []*message.OutputMessage, 50)

	go chess_engine.RunMovesAnalyzer(cfg, validInputMovesChannel, calculatedMovesChannel, 0, 8)
	go broadcast.RunBroadcaster(cfg.DisplayHost, cfg.DisplayPort, outputChannel)
	go HandleGame(cfg, inputChannel, validInputMovesChannel, calculatedMovesChannel, outputChannel)

	for {
		fmt.Printf("Waiting for input device to connect on %s\n", address)
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("couldn't accept connection: %v", err)
		}
		fmt.Printf("%s successfully connected\n", conn.RemoteAddr())

		receiveData(conn, inputChannel)

		fmt.Printf("%s closed connection\n", conn.RemoteAddr())
		conn.Close()
	}
}

func receiveData(conn net.Conn, inputChannel chan<- string) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Printf("Received: %s\n", text)
		if text == "exit" {
			return
		}
		inputChannel <- text
	}
}
