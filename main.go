package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"Shashintary/console_display"
	"Shashintary/console_input"
	"Shashintary/services/config"
	"Shashintary/services/program_interface"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Not enough arguments. Choose \"input\", \"output\" or \"display\"")
		return
	}
	mode := strings.ToLower(os.Args[1])

	switch mode {
	case "input":
		runInput()
	case "output":
		runOutput()
	case "display":
		runDisplay()
	default:
		fmt.Println("Unknown mode. Choose \"input\", \"output\" or \"display\"")
	}
}

func runInput() {
	host, port := config.LoadHostPort(true)
	err := console_input.RunInput(host, port)
	finishProgram(err)
}

func runOutput() {
	cfg := config.LoadConfig(true)
	err := program_interface.HandleProgram(cfg)
	finishProgram(err)
}

func runDisplay() {
	host, port := config.LoadHostPort(false)
	err := console_display.RunDisplay(host, port)
	finishProgram(err)
}

func finishProgram(err error) {
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Program finished")
}
