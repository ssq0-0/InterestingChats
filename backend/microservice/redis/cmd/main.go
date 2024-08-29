package main

import (
	"InterestingChats/backend/microservice/redis/internal/rdb"
	"InterestingChats/backend/microservice/redis/internal/server"
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Write filepath!")
		return
	}

	redisClient, err := rdb.Connect(args[1])
	if err != nil {
		log.Fatalf("Failed to connect Redis: %v\n", err)
	}

	rdbClient := rdb.NewRedisClient(redisClient)
	srv := server.NewServer(rdbClient)
	srv.Start()
}
