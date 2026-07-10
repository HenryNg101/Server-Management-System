package config

import (
	"log"
	"os"
	"strings"

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

type ElasticSearchConfig struct {
	Host                  string
	User                  string
	Password              string
	Port                  string
	StatusDataStreamName  string
	MetricsDataStreamName string
}

type KafkaConfig struct {
	Brokers                    []string
	AgentMetricsTopic          string
	AgentMetricsConsumerGroup  string
	AgentMetricsDLQTopic       string
	ServersImportTopic         string
	ServersImportConsumerGroup string
}

type MinIOConfig struct {
	InternalEndpoint string
	PublicEndpoint   string
	AccessKey        string
	SecretKey        string
	Bucket           string
	UseSSL           bool
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func LoadAppConfig() *ApplicationConfig {
	return &ApplicationConfig{
		Host: getEnv("SERVER_SERVICE_HOST", "localhost"),
		Port: getEnv("SERVER_SERVICE_PORT", "8083"),
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
		Host:                  getEnv("ELASTIC_HOST", "localhost"),
		User:                  getEnv("ELASTIC_USER", "elastic"),
		Password:              getEnv("ELASTIC_PASSWORD", ""),
		Port:                  getEnv("ELASTIC_PORT", "9200"),
		StatusDataStreamName:  getEnv("ELASTIC_STATUS_DATA_STREAM_SOURCE", "server-status"),
		MetricsDataStreamName: getEnv("ELASTIC_METRICS_DATA_STREAM_SOURCE", "server-metrics"),
	}
}

func LoadJWTSecret() string {
	return getEnv("JWT_SECRET", "")
}

func LoadKafka() *KafkaConfig {
	return &KafkaConfig{
		Brokers:                   strings.Split(getEnv("KAFKA_BROKERS", "kafka:9092"), ","),
		AgentMetricsTopic:         getEnv("KAFKA_AGENT_METRICS_TOPIC", "agent-metrics"),
		AgentMetricsConsumerGroup: getEnv("KAFKA_AGENT_METRICS_GROUP", "agent-metrics-consumer-group"),
		AgentMetricsDLQTopic:      getEnv("KAFKA_AGENT_METRICS_DLQ_TOPIC", "agent-metrics-dlq"),
	}
}

func LoadMinIO() *MinIOConfig {
	return &MinIOConfig{
		InternalEndpoint: getEnv("MINIO_INTERNAL_ENDPOINT", "minio:9000"),
		PublicEndpoint:   getEnv("MINIO_PUBLIC_ENDPOINT", "host.docker.internal:9000"),
		AccessKey:        getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		SecretKey:        getEnv("MINIO_SECRET_KEY", "minioadmin"),
		Bucket:           getEnv("MINIO_BUCKET", "jobs"),
		UseSSL:           false,
	}
}
