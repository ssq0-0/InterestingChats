package config

import (
	"encoding/json"
	"os"
)

// ConfigKafkaConsumer defines the structure for Kafka consumer settings.
type ConfigKafkaConsumer struct {
	EnableAutoCommit     bool `json:"enable_auto_commit"`
	AutoCommitIntervalMs int  `json:"auto_commit_interval_ms"`
	SessionTimeoutMs     int  `json:"session_timeout_ms"`
}

// ConfigKafka defines Kafka configuration structure.
type ConfigKafka struct {
	BootstrapServers string              `json:"bootstrap_servers"`
	GroupID          string              `json:"group_id"`
	AutoOffsetReset  string              `json:"auto_offset_reset"`
	Topics           []string            `json:"topics"`
	Consumer         ConfigKafkaConsumer `json:"consumer"`
}

// ConfigRedis defines Redis configuration structure.
type ConfigRedis struct {
	Addr       string `json:"addr"`
	Password   string `json:"password"`
	DB         int    `json:"db"`
	SessionTTL int    `json:"session_ttl"`
}

// Config defines the structure of the configuration file.
type Config struct {
	Kafka ConfigKafka `json:"kafka"`
	Redis ConfigRedis `json:"redis"`
}

// LoadConfig reads and decodes the configuration file into a Config struct.
func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
