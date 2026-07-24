package bootstrap

import (
	"encoding/json"
	"net/http"

	"github.com/HenryNg101/docker-monitoring-agent/internal/security"
)

func RotateKey(cfg *Config, sec *Secret) error {
	// Decrypt before sending
	key := security.DeriveKey(cfg.InstanceID)
	rawKey, err := security.Decrypt(sec.APIKeyEncrypted, key)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest("POST", cfg.APIURL+"/agents/rotate-key", nil)
	req.Header.Set("X-Agent-API-Key", rawKey)

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

	// 3. Encrypt NEW key before storing
	encryptedNewKey, err := security.Encrypt(res.APIKey, key)
	if err != nil {
		return err
	}
	// 4. Update ONLY with encrypted value
	sec.APIKeyEncrypted = encryptedNewKey

	return SaveConfig(cfg, sec)
}
