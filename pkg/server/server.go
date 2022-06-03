package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  BufferSize,
	WriteBufferSize: BufferSize,
}

var DNSResolver *net.Resolver

func Server(addr string) {
	log.SetFlags(0)
	http.HandleFunc("/ws", Handle)
	DNSResolver = InitDNSResolver("8.8.8.8:53")
	log.Fatal(http.ListenAndServe(addr, nil))
}

func Handle(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer wsConn.Close()
	host := r.URL.Query().Get("host")
	if host == "" {
		wsConn.WriteMessage(websocket.TextMessage, []byte("host cannot be empty"))
		return
	}

	port, err := strconv.Atoi(r.URL.Query().Get("port"))
	if err != nil || !(0 < port && port < 65536) {
		wsConn.WriteMessage(websocket.TextMessage, []byte("port format error"))
		return
	}

	ip := net.ParseIP(host)
	hostname := ""
	if ip == nil {
		rip, err := DNSResolver.LookupHost(context.Background(), host)
		hostname = host
		if err == nil && len(rip) != 0 {
			ip = net.ParseIP(rip[0])
		}
	}
	if ip == nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte("Temporary failure in name resolution"))
		return
	}
	Process(remoteAddr(r), wsConn, hostname, ip, port)
}
