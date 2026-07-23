package bootstrap

import (
	"log"
	"os"

	"github.com/google/uuid"
)

// Only load the first time
func InitiateAgent() (*Config, *Secret) {
	cfg, secret, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Inject env ONLY if missing
	if cfg.APIURL == "" {
		cfg.APIURL = os.Getenv("API_URL")
		cfg.ServerName = os.Getenv("SERVER_NAME")
	}

	// Generate instance ID if first time
	if cfg.InstanceID == "" {
		cfg.InstanceID = uuid.New().String()
	}

	// Register if not done yet
	if secret.APIKeyEncrypted == "" {
		log.Println("First run: registering agent...")

		accessToken, refreshToken, err := Login(cfg)
		if err != nil {
			log.Fatal("login failed:", err)
		}

		if err := Register(cfg, secret, accessToken); err != nil {
			log.Fatal("register failed:", err)
		}

		if err := SaveConfig(cfg, secret); err != nil {
			log.Fatal("save config failed:", err)
		}

		// logout immediately after it's done
		if err := Logout(cfg.APIURL, refreshToken); err != nil {
			log.Println("logout failed (non-fatal):", err)
		}

		log.Println("Agent registered successfully")
	}
	return cfg, secret
}
