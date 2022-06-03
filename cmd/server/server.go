package main

import "github.com/kur4ge/tcp-over-websockets/pkg/server"

func main() {
	server.Server("0.0.0.0:1337")
}
