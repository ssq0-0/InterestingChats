package main

import (
	"InterestingChats/backend/api_gateway/internal/config"
	"InterestingChats/backend/api_gateway/internal/consts"
	"InterestingChats/backend/api_gateway/internal/server"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	consts.InitConstants(cfg)

	srv := server.NewServer(cfg)
	srv.Start()
}
