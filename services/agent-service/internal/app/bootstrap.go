package app

import (
	"errors"

	"github.com/HenryNg101/agent-service/internal/agent"
	"github.com/HenryNg101/agent-service/internal/config"
	kafkaClient "github.com/HenryNg101/agent-service/internal/platform/kafka"
	"github.com/HenryNg101/agent-service/internal/platform/postgres"
	"gorm.io/gorm"
)

type Application struct {
	PostgresSession *gorm.DB
	AgentService    agent.Service
}

func NewApp() (*Application, error) {
	// Load config
	postgresConfig := config.LoadPostgres()
	if len(postgresConfig.Password) == 0 {
		return nil, errors.New("Password for Postgres connection is not set. You have to set it in .env file in root folder using POSTGRES_PASSWORD variable")
	}
	postgresSession := postgres.NewPostgresSession(postgresConfig)

	kafkaConfig := config.LoadKafka()
	if len(kafkaConfig.Brokers) == 0 {
		return nil, errors.New("Kafka brokers not configured. You have to set it in .env file in root folder using KAFKA_BROKERS variable")
	}

	agentRepo := agent.NewRepository(postgresSession)
	kafkaAgentProducer := agent.NewKafkaProducer(kafkaClient.NewProducer(kafkaConfig.Brokers, kafkaConfig.AgentMetricsTopic))
	agentService := agent.NewService(agentRepo, kafkaAgentProducer)

	return &Application{
		PostgresSession: postgresSession,
		AgentService:    agentService,
	}, nil
}
