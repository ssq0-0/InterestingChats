package main

import (
	"InterestingChats/internal/db"
	"fmt"
	"log"
	"net/http"
)

func main() {
	cfg := db.DataBaseConfig{Host: "localhost", Port: 5432, User: "***", Password: "***", Dbname: "***"}
	conn, err := db.Connect(cfg)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close()

	fmt.Println("Successfully connected to the database!")
	log.Println("Starting server on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Unable to start server: %v\n", err)
	}
}
