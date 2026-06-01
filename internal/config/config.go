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

type MailerConfig struct {
	Server    string
	Port      string
	UserName  string
	Password  string
	FromEmail string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
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
		Password:     getEnv("POSTGRES_PASSWORD", ""),
		DatabaseName: getEnv("POSTGRES_DB", "postgres"),
		Port:         getEnv("POSTGRES_PORT", "5432"),
	}
}

func LoadElasticsearch() *ElasticSearchConfig {
	return &ElasticSearchConfig{
		Host:           getEnv("ELASTIC_HOST", "localhost"),
		User:           getEnv("ELASTIC_USER", "elastic"),
		Password:       getEnv("ELASTIC_PASSWORD", ""),
		Port:           getEnv("ELASTIC_PORT", "9200"),
		DataStreamName: getEnv("ELASTIC_DATA_STREAM_SOURCE", "data_stream"),
	}
}

func LoadMailer() *MailerConfig {
	return &MailerConfig{
		Server:    getEnv("MAIL_SERVER", "smtp.gmail.com"),
		Port:      getEnv("MAIL_PORT", "587"),
		UserName:  getEnv("MAIL_USER", "example@gmail.com"), // This one won't work, you have to do use your own actual email for auth
		Password:  getEnv("MAIL_PASSWORD", ""),              // Same idea. You have to generate from your auth email
		FromEmail: getEnv("MAIL_FROM_USER", "example@gmail.com"),
	}
}

func LoadJWTSecret() string {
	return getEnv("JWT_SECRET", "")
}

func LoadRedis() *RedisConfig {
	return &RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     getEnv("REDIS_PORT", "localhost"),
		Password: getEnv("REDIS_PASSWORD", ""),
	}
}
