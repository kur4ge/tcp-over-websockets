package server

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func Process(src string, wsConn *websocket.Conn, hostname string, ip net.IP, port int) {
	var dest string
	var err error
	uuid := uuid.New().String()

	if hostname != "" {
		dest = fmt.Sprintf("%s(%s):%d", hostname, ip, port)
	} else {
		dest = buildAddr(ip, port)
	}
	fmt.Printf("%s [%s] %s Conect to %s\n", src, datetime(), uuid, dest)

	d := net.Dialer{Timeout: 3 * time.Second} // 3 秒超时
	sConn, err := d.Dial("tcp", buildAddr(ip, port))
	if err != nil {
		log.Printf("%s [%s] %s Conect Err: %v", src, datetime(), uuid, err)
		wsConn.WriteMessage(websocket.TextMessage, []byte("Connection failed."))
		return
	}
	defer sConn.Close()

	sRecv := make(chan []byte) // 创建 chan 监听来自 socket 的数据
	go socketReader(sConn, sRecv)

	wsRecv := make(chan []byte) // 创建 chan 监听来自 webSocket 的数据
	go wsReader(wsConn, wsRecv)

	for {
		select {
		case data, ok := <-sRecv:
			if !ok { // 代表 socket 端断开
				log.Printf("%s [%s] %s %s Close", src, datetime(), uuid, dest)
				wsConn.WriteMessage(websocket.TextMessage, []byte("Connection closed."))
				return
			}
			fmt.Printf("%s [%s] %s Recv(%d): %v\n", src, datetime(), uuid, len(data), data)
			err = wsConn.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				log.Printf("%s [%s] %s Sendto %s Err: %v", src, datetime(), uuid, src, err)
				return
			}
		case data, ok := <-wsRecv:
			if !ok { // 代表 ws 端断开
				log.Printf("%s [%s] %s %s Close", src, datetime(), uuid, src)
				return
			}
			fmt.Printf("%s [%s] %s Send(%d): %v\n", src, datetime(), uuid, len(data), data)
			_, err = sConn.Write(data)
			if err != nil {
				log.Printf("%s [%s] %s Sendto %s Err: %v", src, datetime(), uuid, dest, err)
				return
			}
		}
	}
}
