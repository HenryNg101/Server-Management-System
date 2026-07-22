package bootstrap

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

func getHostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return h
}

func Register(cfg *Config, token string) error {
	body := map[string]string{
		"server_name": cfg.ServerName,
		"hostname":    getHostname(),
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

	cfg.APIKey = res.APIKey
	cfg.ServerID = res.ServerID

	return nil
}
