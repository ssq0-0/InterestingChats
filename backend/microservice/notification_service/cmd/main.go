package main

import (
	"log"
	"notifications/internal/config"
	"notifications/internal/consts"
	"notifications/internal/server"
)

func main() {
	// Load configuration from the JSON file
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err) // Exit if config loading fails
	}

	consts.InitServiceURLs(
		cfg.Services["db_service"].Protocol, cfg.Services["db_service"].Host, cfg.Services["db_service"].Port,
	)

	srv := server.NewServer(cfg.Server.Host, cfg.Server.Port) // Pass host and port from config
	go func() {
		log.Println("Starting server...")
		srv.Start()
	}()

	cnsm, err := server.NewConsumer(cfg.Kafka) // Pass Kafka configuration from config
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	log.Printf("Start subscribing topics...")
	go func() {
		log.Printf("Subscribing to topics...")
		if err := cnsm.Subscriber(cfg.Kafka.Topics); err != nil {
			log.Fatalf("Failed to subscribe: %v", err)
		}
		log.Printf("Start reading topics...")
		cnsm.Reader()
	}()

	select {}
}
