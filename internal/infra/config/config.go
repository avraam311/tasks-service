package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Server *Server
}

type Server struct {
	Port string `json:"port" validate:"required"`
}

func New() (*Config, error) {
	return &Config{}, nil
}

func (c *Config) LoadJSON(fileName string) error {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("config/config.go - failed to read json - %w", err)
	}

	if err := json.Unmarshal(data, &c); err != nil {
		return fmt.Errorf("config/config.go - failed to unmarshal data - %w", err)
	}

	return nil
}
