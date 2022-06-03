package server

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"net"

	"github.com/gorilla/websocket"
)

func wsReader(conn *websocket.Conn, recv chan []byte) {
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
			if len(data)%2 == 0 && reHex.Match(data) { // 判断是 hex字符串
				n, err := hex.Decode(data, data)
				if err == nil {
					data = data[:n]
				}
			} else if len(data)%4 == 0 && reBase64.Match(data) { // 判断是 base64串
				n, err := base64.StdEncoding.Decode(data, data)
				if err == nil {
					data = data[:n]
				}
			}
		}
		recv <- data
	}
}

func socketReader(conn net.Conn, recv chan []byte) {
	connReader := bufio.NewReader(conn)
	data := make([]byte, BufferSize)
	for {
		length, err := connReader.Read(data)
		if err != nil {
			close(recv)
			return
		}
		recv <- data[:length]
	}
}
