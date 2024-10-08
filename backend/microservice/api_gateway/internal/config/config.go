package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config defines the structure of the entire configuration file.
type Config struct {
	Port     string `json:"port"`
	Services struct {
		UserService         string `json:"user_service"`
		ChatService         string `json:"chat_service"`
		RedisService        string `json:"redis_service"`
		FileSystemService   string `json:"file_system_service"`
		NotificationService string `json:"notification_service"`
		AuthService         string `json:"auth_service"`
	} `json:"services"`
}

// LoadConfig reads and decodes the configuration file into a Config struct.
func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл конфигурации: %v", err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("не удалось декодировать конфигурацию: %v", err)
	}

	return &config, nil
}
