package bootstrap

import (
	"encoding/json"
	"os"
)

const configPath = "/app/data/config.json"

type Config struct {
	APIURL     string `json:"api_url"`
	ServerName string `json:"server_name"`

	InstanceID string `json:"instance_id"`
	APIKey     string `json:"api_key"`
	ServerID   int    `json:"server_id"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return &Config{}, nil // first run
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func SaveConfig(cfg *Config) error {
	if err := os.MkdirAll("/app/data", 0755); err != nil {
		return err
	}

	data, _ := json.MarshalIndent(cfg, "", "  ")
	return os.WriteFile(configPath, data, 0600)
}
