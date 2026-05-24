package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ApplicationConfig struct {
	Host string
	Port string
}

func getEnv(key, fallback string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func LoadAppConfig() *ApplicationConfig {
	return &ApplicationConfig{
		Host: getEnv("APP_HOST", "localhost"),
		Port: getEnv("APP_PORT", "8080"),
	}
}
