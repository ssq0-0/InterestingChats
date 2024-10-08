package main

import (
	"InterestingChats/backend/microservice/redis/internal/config"
	"InterestingChats/backend/microservice/redis/internal/rdb"
	"InterestingChats/backend/microservice/redis/internal/server"
	"InterestingChats/backend/microservice/redis/internal/services"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
)

func main() {
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	var redisClient *rdb.RedisClient
	var client *redis.Client

	// Попытка подключения к Redis с задержкой
	for i := 0; i < 10; i++ {
		client, err = rdb.Connect(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to Redis, retrying in 5 seconds: %v", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("Failed to connect to Redis after retries: %v", err)
	}
	log.Println("Successfully connected to Redis")

	// Создаем экземпляр RedisClient
	redisClient = rdb.NewRedisClient(client)

	consumer, err := services.NewConsumer(redisClient, cfg)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v\n", err)
	}

	if err = consumer.Subscriber(cfg.Kafka.Topics); err != nil {
		log.Fatalf("Failed to subscribe to Kafka topics: %v\n", err)
	}

	log.Println("Successfully subscribed to Kafka topics")
	go func() {
		err = consumer.Reader(cfg)
		if err != nil {
			log.Fatalf("Error in Kafka reader: %v", err)
		}
	}()

	srv := server.NewServer(redisClient)
	go srv.Start()
	log.Println("Server started successfully")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("Received signal: %s, shutting down...", sig)
}
