package bootstrap

import (
	"encoding/json"
	"net/http"
)

func RotateKey(cfg *Config, sec *Secret) error {
	req, _ := http.NewRequest("POST", cfg.APIURL+"/agents/rotate-key", nil)
	req.Header.Set("X-Agent-API-Key", sec.APIKeyEncrypted)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var res struct {
		APIKey string `json:"api_key"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	// Update in-memory, then persist to disk (overwrite old key)
	oldKey := sec.APIKeyEncrypted
	sec.APIKeyEncrypted = res.APIKey

	if err := SaveConfig(cfg, sec); err != nil {
		// rollback in memory to the old value, if it fails to save it
		sec.APIKeyEncrypted = oldKey
		return err
	}

	return SaveConfig(cfg, sec)
}
