package main

import (
	"auth_service/internal/config"
	"auth_service/internal/server"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Unable to load config: %v", err)
	}

	srv := server.NewServer(cfg)
	srv.Start()
}
