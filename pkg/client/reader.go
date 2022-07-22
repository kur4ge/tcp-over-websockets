package client

import (
	"bufio"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

func WSReader(conn *websocket.Conn, recv chan []byte) {
	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			close(recv)
			return
		}
		if len(data) == 0 {
			continue
		}
		switch messageType {
		case websocket.TextMessage:
			log.Println(string(data))
			continue
		}
		recv <- data
	}
}

func StdinReader(recv chan []byte) {
	input := bufio.NewReader(os.Stdin)
	for {
		data, err := input.ReadBytes('\n')
		if err != nil {
			close(recv)
			return
		}
		recv <- data
	}
}
