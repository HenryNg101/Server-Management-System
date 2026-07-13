package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Load the stuff first
func init() {
	// Prevent accidental loading of env variables to docker containers for services
	if os.Getenv("APP_ENV") != "docker" {
		err := godotenv.Load()
		if err != nil {
			log.Println("No .env file found, using environment variables")
		}
	}
}

type ApplicationConfig struct {
	Host string
	Port string
}

type PostgresConfig struct {
	Host         string
	User         string
	Password     string
	DatabaseName string
	Port         string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func LoadAppConfig() *ApplicationConfig {
	return &ApplicationConfig{
		Host: getEnv("HOST", "localhost"),
		Port: getEnv("AUTH_SERVICE_PORT", "8087"),
	}
}

func LoadPostgres() *PostgresConfig {
	return &PostgresConfig{
		Host:         getEnv("POSTGRES_HOST", "localhost"),
		User:         getEnv("POSTGRES_USER", "postgres"),
		Password:     getEnv("POSTGRES_PASSWORD", ""),
		DatabaseName: getEnv("POSTGRES_DB", "postgres"),
		Port:         getEnv("POSTGRES_PORT", "5432"),
	}
}

func LoadJWTSecret() string {
	return getEnv("JWT_SECRET", "")
}

func LoadRedis() *RedisConfig {
	return &RedisConfig{
		Host:     getEnv("REDIS_HOST", "redis"),
		Port:     getEnv("REDIS_PORT", "6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
	}
}
