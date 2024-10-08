package main

import (
	"InterestingChats/backend/microservice/db/internal/config"
	"InterestingChats/backend/microservice/db/internal/db"
	"InterestingChats/backend/microservice/db/internal/server"
	"database/sql"
	"log"
	"time"
)

// main to initialize configuration, database, server, and Kafka consumer
func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	var database *sql.DB
	for i := 0; i < 10; i++ {
		database, err = db.Connect(config.DatabaseConfig{
			Host:     cfg.Database.Host,
			Port:     cfg.Database.Port,
			User:     cfg.Database.User,
			Password: cfg.Database.Password,
			Dbname:   cfg.Database.Dbname,
		})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database, retrying in 5 seconds: %v", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to database after retries: %v", err)
	}
	defer database.Close()

	srv := server.NewServer(database, cfg.Server.Port)
	go func() {
		log.Println("Starting server...")
		srv.Start()
	}()

	cnsm, err := server.NewConsumer(srv.Handler.GetDBService(), cfg.Kafka)
	if err != nil {
		log.Printf("err: %v", err)
		return
	}
	log.Printf("start subscribe topics...")
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
