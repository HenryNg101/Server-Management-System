package app

import (
	"errors"

	"github.com/HenryNg101/agent-service/internal/agent"
	"github.com/HenryNg101/agent-service/internal/config"
	"github.com/HenryNg101/agent-service/internal/platform/elastic"
	kafkaClient "github.com/HenryNg101/agent-service/internal/platform/kafka"
	"github.com/HenryNg101/agent-service/internal/platform/postgres"
	"github.com/elastic/go-elasticsearch/v9"
	"gorm.io/gorm"
)

type Application struct {
	PostgresSession *gorm.DB
	ElasticSession  *elasticsearch.Client
	AgentService    agent.Service
}

func NewApp() (*Application, error) {
	// Load config
	postgresConfig := config.LoadPostgres()
	if len(postgresConfig.Password) == 0 {
		return nil, errors.New("Password for Postgres connection is not set. You have to set it in .env file in root folder using POSTGRES_PASSWORD variable")
	}
	postgresSession := postgres.NewPostgresSession(postgresConfig)

	elasticConfig := config.LoadElasticsearch()
	if len(elasticConfig.Password) == 0 {
		return nil, errors.New("Password for Elasticsearch connection is not set. You have to set it in .env file in root folder using ELASTIC_PASSWORD variable")
	}
	esSession := elastic.NewElasticsearchSession(elasticConfig)

	kafkaConfig := config.LoadKafka()
	if len(kafkaConfig.Brokers) == 0 {
		return nil, errors.New("Kafka brokers not configured. You have to set it in .env file in root folder using KAFKA_BROKERS variable")
	}

	agentRepo := agent.NewRepository(postgresSession)
	elasticAgentRepo := agent.NewAgentESRepository(esSession)
	kafkaAgentProducer := agent.NewKafkaProducer(kafkaClient.NewProducer(kafkaConfig.Brokers, kafkaConfig.AgentMetricsTopic))
	agentService := agent.NewService(agentRepo, kafkaAgentProducer, elasticAgentRepo)

	return &Application{
		PostgresSession: postgresSession,
		ElasticSession:  esSession,
		AgentService:    agentService,
	}, nil
}
