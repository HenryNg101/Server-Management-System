package bootstrap

import (
	"encoding/json"
	"os"

	"github.com/HenryNg101/docker-monitoring-agent/internal/security"
)

// const configPath = "/app/data/config.json"
const (
	configDir  = "/app/data"
	configPath = "/app/data/.agent_config.json"
	secretPath = "/app/data/.agent_secret"
)

type Config struct {
	APIURL     string `json:"api_url"`
	ServerName string `json:"server_name"`

	InstanceID string `json:"instance_id"`
	ServerID   int    `json:"server_id"`
}

// secrets are saved in seperate files
type Secret struct {
	APIKeyEncrypted string `json:"api_key"`
}

func LoadConfig() (*Config, *Secret, error) {
	_ = os.MkdirAll(configDir, 0755)

	cfg := &Config{}
	sec := &Secret{}

	// load config
	if data, err := os.ReadFile(configPath); err == nil {
		if err := json.Unmarshal(data, cfg); err != nil {
			return nil, nil, err
		}
	}

	// load secret
	if data, err := os.ReadFile(secretPath); err == nil {
		if err := json.Unmarshal(data, sec); err != nil {
			return nil, nil, err
		}
	}

	// Decrypt secret after load in here
	if sec.APIKeyEncrypted != "" && cfg.InstanceID != "" {
		key := security.DeriveKey(cfg.InstanceID)

		decrypted, err := security.Decrypt(sec.APIKeyEncrypted, key)
		if err != nil {
			return nil, nil, err
		}

		sec.APIKeyEncrypted = decrypted
	}

	return cfg, sec, nil
}

func SaveConfig(cfg *Config, sec *Secret) error {
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// save config
	cfgData, _ := json.MarshalIndent(cfg, "", "  ")
	if err := os.WriteFile(configPath, cfgData, 0644); err != nil {
		return err
	}

	//
	// Encrypt the content
	key := security.DeriveKey(cfg.InstanceID)
	encrypted, err := security.Encrypt(sec.APIKeyEncrypted, key)
	if err != nil {
		return err
	}
	sec.APIKeyEncrypted = encrypted

	// save secret (stricter permission)
	secData, _ := json.MarshalIndent(sec, "", "  ")
	if err := os.WriteFile(secretPath, secData, 0600); err != nil {
		return err
	}

	return nil
}
