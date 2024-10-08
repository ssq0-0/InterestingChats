package config

import (
	"encoding/json"
	"os"
)

// Configuration holds the application configuration data
type Configuration struct {
	Server   ServerConfig             `json:"server"`
	Services map[string]ServiceConfig `json:"services"` // Map to hold multiple external services
	Kafka    KafkaConfig              `json:"kafka"`
}

// ServerConfig holds the server configuration
type ServerConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// ServiceConfig holds the configuration for an external service
type ServiceConfig struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

// KafkaConfig holds the Kafka configuration
type KafkaConfig struct {
	BootstrapServers     string   `json:"bootstrap_servers"`
	GroupID              string   `json:"group_id"`
	AutoOffsetReset      string   `json:"auto_offset_reset"`
	EnableAutoCommit     bool     `json:"enable_auto_commit"`
	AutoCommitIntervalMs int      `json:"auto_commit_interval_ms"`
	Topics               []string `json:"topics"`
}

// LoadConfig loads the configuration from a JSON file
func LoadConfig(filePath string) (Configuration, error) {
	var config Configuration

	// Open the JSON file
	file, err := os.Open(filePath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	// Decode the JSON file into the Configuration struct
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}
