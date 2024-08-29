package rdb

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(client *redis.Client) *RedisClient {
	return &RedisClient{
		client: client,
	}
}

func Readfile(filepath string) (map[string]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed read file path: %v", err)
	}
	defer file.Close()

	data := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key, value := parts[0], parts[1]
			data[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error read file: %v", err)
	}

	return data, nil
}

func Connect(filepath string) (*redis.Client, error) {
	fileData, err := Readfile(filepath)
	if err != nil {
		return nil, fmt.Errorf("func 'ReadFile' exec error: %v", err)
	}

	db, err := strconv.Atoi(fileData["db"])
	if err != nil {
		return nil, fmt.Errorf("can't find db integer: %v", err)
	}

	rdbClient := redis.NewClient(&redis.Options{
		Addr:     fileData["addr"],
		Password: fileData["password"],
		DB:       db,
	})

	return rdbClient, nil
}
