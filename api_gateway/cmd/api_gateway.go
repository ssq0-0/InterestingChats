package main

import (
	"InterestingChats/backend/api_gateway/internal/server"
)

func main() {
	server := server.NewServer()
	server.Start()
}
