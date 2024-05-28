package main

import (
	"LingerAPI/internal/server"
)

func main() {
	server := server.NewLingerServer()
	server.StartServer()
}
