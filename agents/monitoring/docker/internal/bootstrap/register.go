package bootstrap

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/HenryNg101/docker-monitoring-agent/internal/security"
)

func Register(cfg *Config, secret *Secret, token string) error {
	body := map[string]string{
		"server_name": cfg.ServerName,
		"instance_id": cfg.InstanceID,
	}

	b, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", cfg.APIURL+"/agents/register", bytes.NewBuffer(b))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var res struct {
		ServerID int    `json:"server_id"`
		APIKey   string `json:"api_key"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	// Encrypt NEW key before storing
	key := security.DeriveKey(cfg.InstanceID)
	encryptedNewKey, err := security.Encrypt(res.APIKey, key)
	if err != nil {
		return err
	}

	secret.APIKeyEncrypted = encryptedNewKey
	cfg.ServerID = res.ServerID

	return nil
}
