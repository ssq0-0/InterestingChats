package main

import (
	"InterestingChats/backend/user_services/internal/db"
	"InterestingChats/backend/user_services/internal/server"
	"log"
)

func main() {
	cfg := db.DataBaseConfig{Host: "localhost", Port: 5432, User: "ssq", Password: "ScR26011161", Dbname: "petproject"}
	conn, err := db.Connect(cfg)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close()

	srv := server.NewServer(conn)
	srv.Start()
}
