package client

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func Client(url string) {
	wsConn, _, err := Connect(url)
	if err != nil {
		log.Fatal(err)
	}
	wsRecv := make(chan []byte)
	go wsReader(wsConn, wsRecv)
	stdRecv := make(chan []byte)
	go stdinReader(stdRecv)

	for {
		select {
		case data, ok := <-stdRecv:
			if !ok {
				return
			}
			err = wsConn.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				log.Printf("Send err: %v", err)
				return
			}
		case data, ok := <-wsRecv:
			if !ok { // 代表 ws 端断开
				log.Printf("WebSocket Close")
				return
			}
			fmt.Print(string(data))
		}
	}
}

func Connect(url string) (*websocket.Conn, *http.Response, error) {
	return websocket.DefaultDialer.Dial(url, nil)
}
