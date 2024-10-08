package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type ConfigKafkaTopics struct {
	FriendsOperation string `json:"friends_operation"`
	SessionSet       string `json:"session_set"`
	SessionUpdate    string `json:"session_update"`
	PushFriends      string `json:"push_friends"`
	PushSubscribers  string `json:"push_subscribers"`
	UpdateSubscriber string `json:"update_subscriber"`
	UpdateFriend     string `json:"update_friend"`
}

type Config struct {
	Server struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"server"`
	Kafka struct {
		BootstrapServers string            `json:"bootstrap_servers"`
		Topics           ConfigKafkaTopics `json:"topics"`
	} `json:"kafka"`
	Services struct {
		DBService    string `json:"db_service"`
		FileService  string `json:"file_service"`
		AuthService  string `json:"auth_service"`
		CacheService string `json:"cache_service"`
	} `json:"services"`
}

func LoadConfig(file string) (*Config, error) {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл конфигурации: %v", err)
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("не удалось декодировать конфигурацию: %v", err)
	}

	return &config, nil
}
