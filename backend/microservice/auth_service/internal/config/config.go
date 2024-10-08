package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config defines the structure of the configuration file, containing server and secret key for generate tokens settings.
type Config struct {
	Server struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"server"`
	JWT struct {
		Secret string `json:"secret"`
	} `json:"jwt"`
}

// LoadConfig reads and decodes the configuration file into a Config struct.
func LoadConfig(file string) (*Config, error) {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %v", err)
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("could not decode config JSON: %v", err)
	}

	return &config, nil
}
