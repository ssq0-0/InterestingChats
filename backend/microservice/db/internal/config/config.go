package config

import (
	"encoding/json"
	"os"
)

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

type KafkaConfig struct {
	BootstrapServers     string   `json:"bootstrap_servers"`
	GroupID              string   `json:"group_id"`
	AutoOffsetReset      string   `json:"auto_offset_reset"`
	EnableAutoCommit     bool     `json:"enable_auto_commit"`
	AutoCommitIntervalMs int      `json:"auto_commit_interval_ms"`
	Topics               []string `json:"topics"`
}

type ServerConfig struct {
	Port int `json:"port"`
}

type Config struct {
	Database DatabaseConfig `json:"database"`
	Kafka    KafkaConfig    `json:"kafka"`
	Server   ServerConfig   `json:"server"`
}

// LoadConfig to load configuration from a JSON file
func LoadConfig(filePath string) (Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
