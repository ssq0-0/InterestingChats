package main

import (
	"chat_service/internal/config"
	"chat_service/internal/consts"
	"chat_service/internal/server"
	"log"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	consts.InitConstants(cfg)

	// Create and start the server
	srv := server.NewServer(cfg)
	srv.Start()
}
