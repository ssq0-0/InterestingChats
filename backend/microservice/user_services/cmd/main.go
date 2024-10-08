package main

import (
	"InterestingChats/backend/user_services/internal/config"
	"InterestingChats/backend/user_services/internal/consts"
	"InterestingChats/backend/user_services/internal/server"
	"log"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	consts.InitConstants(cfg)

	// Pass the configuration to the server
	srv := server.NewServer(cfg)
	srv.Start()
}
