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

type PostgresConfig struct {
	Host         string
	User         string
	Password     string
	DatabaseName string
	Port         string
}

type ElasticSearchConfig struct {
	Host           string
	User           string
	Password       string
	Port           string
	DataStreamName string
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

func LoadPostgres() *PostgresConfig {
	return &PostgresConfig{
		Host:         getEnv("POSTGRES_HOST", "localhost"),
		User:         getEnv("POSTGRES_USER", "postgres"),
		Password:     getEnv("POSTGRES_PASSWORD", "password"),
		DatabaseName: getEnv("POSTGRES_DB", "postgres"),
		Port:         getEnv("POSTGRES_PORT", "5432"),
	}
}

func LoadElasticsearch() *ElasticSearchConfig {
	return &ElasticSearchConfig{
		Host:           getEnv("ELASTIC_HOST", "localhost"),
		User:           getEnv("ELASTIC_USER", "elastic"),
		Password:       getEnv("ELASTIC_PASSWORD", "password"),
		Port:           getEnv("ELASTIC_PORT", "9200"),
		DataStreamName: getEnv("ELASTIC_DATA_STREAM_SOURCE", "data_stream"),
	}
}
