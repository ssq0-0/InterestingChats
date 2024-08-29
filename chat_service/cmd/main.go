package main

import "chat_service/internal/server"

func main() {
	srv := server.NewServer()
	srv.Start()
}
