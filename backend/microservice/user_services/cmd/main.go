package main

import (
	"InterestingChats/backend/user_services/internal/server"
)

func main() {
	srv := server.NewServer()
	srv.Start()
}
