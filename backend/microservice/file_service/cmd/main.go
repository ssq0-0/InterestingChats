package main

import (
	"file_service/internal/config"
	"file_service/internal/server"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("config.json") // Load configuration
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	srv := server.NewServer(cfg) // Pass config to the server
	srv.Start(cfg)               // Pass config to start the server
}
