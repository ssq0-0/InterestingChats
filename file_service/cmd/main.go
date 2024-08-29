package main

import "file_service/internal/server"

func main() {
	srv := server.NewServer()
	srv.Start()
}
