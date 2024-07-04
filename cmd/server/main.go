package main

import (
	"InterestingChats/internal/db"
	"InterestingChats/internal/services"
	"log"
)

func main() {
	cfg := db.DataBaseConfig{Host: "localhost", Port: 5432, User: "ssq", Password: "ScR26011161", Dbname: "petproject"}
	conn, err := db.Connect(cfg)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close()

	server := services.NewServer(conn)
	server.Start()
}
