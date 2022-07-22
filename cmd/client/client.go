package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kur4ge/tcp-over-websockets/pkg/client"
)

func main() {
	wsConn, _, err := Connect("ws://127.0.0.1:1337/ws?host=baidu.com&port=80")
	if err != nil {
		log.Fatal(err)
	}
	wsRecv := make(chan []byte)
	go client.WSReader(wsConn, wsRecv)
	stdRecv := make(chan []byte)
	go client.StdinReader(stdRecv)

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
