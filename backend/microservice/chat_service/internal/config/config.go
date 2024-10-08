package config

import (
	"encoding/json"
	"os"
)

// KafkaConfig holds the configuration details for Kafka, including bootstrap servers and topic mappings.
type KafkaConfig struct {
	BootstrapServers string   `json:"bootstrapServers"`
	Topics           []string `json:"topics"`
}

// ServerConfig holds the server-related configurations, including the port number.
type ServerConfig struct {
	Port string `json:"port"`
}

// LoggingConfig holds the logging configuration, specifically the log level.
type LoggingConfig struct {
	Level string `json:"level"`
}

// Config is the main configuration struct that aggregates server, Kafka, and logging configurations.
type Config struct {
	Server   ServerConfig  `json:"server"`
	Kafka    KafkaConfig   `json:"kafka"`
	Logging  LoggingConfig `json:"logging"`
	Services struct {
		DBService string `json:"db_service"`
	} `json:"services"`
}

// LoadConfig reads and parses the configuration from a given JSON file.
func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
