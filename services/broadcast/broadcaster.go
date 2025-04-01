package broadcast

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"Shashintary/modules/message"
)

func RunBroadcaster(host, port string, outputChannel <-chan []*message.OutputMessage) {
	if host == "" {
		log.Fatal("RunBroadcaster: host is empty")
	}
	if port == "" {
		log.Fatal("RunBroadcaster: port is empty")
	}
	address := fmt.Sprintf("%s:%s", host, port)

	for {
		var err error
		var conn net.Conn
		for {
			fmt.Printf("Broadcaster: waiting for '%s' to open connection\n", address)
			conn, err = net.Dial("tcp", address)
			if err == nil {
				fmt.Printf("Broadcaster: successfully connected to '%s' and ready to broadcast\n", address)
				break
			}
			time.Sleep(5 * time.Second)
		}

		if !broadcastMessages(outputChannel, &conn) {
			fmt.Printf("Broadcaster: exiting...")
			conn.Close()
			return
		}

		fmt.Printf("Broadcaster: %s closed connection\n", address)
		conn.Close()
	}
}

// broadcastMessages returns true if exits because of connection lost, returns false if outputChannel is closed.
func broadcastMessages(outputChannel <-chan []*message.OutputMessage, conn *net.Conn) bool {
	for msgs := range outputChannel {
		var s string

		for _, msg := range msgs {
			if msg.IsBoard {
				s = insertHashCharAfterNewlines(msg.Value)
			} else {
				s += msg.Value
			}
		}
		for {
			_, err := fmt.Fprintf(*conn, s)
			if err == nil {
				break
			}
			return true
			//fmt.Printf("Broadcaster: couldn't send message `%s` to `%s`, trying again in 5 sec.\n", s, (*conn).RemoteAddr())
			//time.Sleep(5 * time.Second)
		}
		fmt.Printf("Broadcaster: broadcasted: %s\n", s)
	}
	return false
}

func insertHashCharAfterNewlines(input string) string {
	var result strings.Builder
	n := len(input)
	for i := 0; i < n; i++ {
		result.WriteByte(input[i])
		if input[i] == '\n' && (i < n-1 && strings.ContainsRune(input[i+1:], '\n')) {
			result.WriteByte('#')
		}
	}
	return result.String()
}
