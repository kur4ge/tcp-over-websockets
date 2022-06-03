package main

import (
	"github.com/kur4ge/tcp-over-websockets/pkg/client"
)

func main() {
	client.Client("ws://127.0.0.1:1337/ws?host=baidu.com&port=80")
}
