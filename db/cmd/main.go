package main

import (
	"InterestingChats/backend/microservice/db/internal/db"
	"InterestingChats/backend/microservice/db/internal/server"
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		log.Println("Write filepath!")
		return
	}

	database, err := db.Connect(args[1])
	if err != nil {
		log.Fatalf("Failed to connect to database: %v\n", err)
	}
	defer database.Close()

	srv := server.NewServer(database)
	srv.Start()
}
