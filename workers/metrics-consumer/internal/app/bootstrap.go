package app

import (
	"errors"

	"github.com/HenryNg101/metrics-consumer/internal/agent"
	"github.com/HenryNg101/metrics-consumer/internal/config"
	"github.com/HenryNg101/metrics-consumer/internal/platform/elastic"
	"github.com/elastic/go-elasticsearch/v9"
)

type Application struct {
	ElasticSession *elasticsearch.Client
	AgentService   agent.Service
}

func NewApp() (*Application, error) {
	// Load config
	elasticConfig := config.LoadElasticsearch()
	if len(elasticConfig.Password) == 0 {
		return nil, errors.New("Password for Elasticsearch connection is not set. You have to set it in .env file in root folder using ELASTIC_PASSWORD variable")
	}
	esSession := elastic.NewElasticsearchSession(elasticConfig)

	kafkaConfig := config.LoadKafka()
	if len(kafkaConfig.Brokers) == 0 {
		return nil, errors.New("Kafka brokers not configured. You have to set it in .env file in root folder using KAFKA_BROKERS variable")
	}

	elasticAgentRepo := agent.NewAgentESRepository(esSession)
	agentService := agent.NewService(elasticAgentRepo)

	return &Application{
		ElasticSession: esSession,
		AgentService:   agentService,
	}, nil
}
