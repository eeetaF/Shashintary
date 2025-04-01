package console_display

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

func RunDisplay(host, port string) error {
	address := fmt.Sprintf("%s:%s", host, port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("couldn't run display on %s: %v", address, err)
	}
	defer listener.Close()

	for {
		clearConsole()
		fmt.Printf("Waiting for output device to connect on %s\n", address)
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("couldn't accept connection: %v", err)
		}
		fmt.Printf("%s successfully connected\n", conn.RemoteAddr())

		receiveData(conn)

		fmt.Printf("%s closed connection\n", conn.RemoteAddr())
		conn.Close()
	}

	return nil
}

func receiveData(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	var stopPrint chan struct{}
	var printMu sync.Mutex
	var printing sync.WaitGroup
	lastBoard := false

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		if line[0] == '@' {
			clearConsole()
			line = line[1:]
			if line == "" {
				continue
			}
		}

		if line[0] == '#' {
			fmt.Printf(" %s\n", line[1:])
			lastBoard = true

			printMu.Lock()
			if stopPrint != nil {
				close(stopPrint)
				stopPrint = nil
			}
			printMu.Unlock()

			continue
		}

		if lastBoard {
			fmt.Println()
			lastBoard = false
		}

		printMu.Lock()
		if stopPrint != nil {
			close(stopPrint)
		}
		stopPrint = make(chan struct{})
		printMu.Unlock()

		printing.Add(1)
		go func(line string, stop <-chan struct{}) {
			defer printing.Done()
			fmt.Print("\t")
			printWithDelay(line, stop)
		}(line, stopPrint)
	}

	printing.Wait()
}

func printWithDelay(line string, stop <-chan struct{}) {
	for _, ch := range line {
		select {
		case <-stop:
			return
		default:
			fmt.Printf("%c", ch)
			time.Sleep(50 * time.Millisecond)
		}
	}
	fmt.Println()
}

func clearConsole() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default: // Linux, macOS
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
