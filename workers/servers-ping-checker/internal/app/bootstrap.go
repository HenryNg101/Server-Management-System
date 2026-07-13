package app

import (
	"errors"

	"github.com/HenryNg101/servers-ping-checker/internal/config"
	"github.com/HenryNg101/servers-ping-checker/internal/platform/elastic"
	"github.com/HenryNg101/servers-ping-checker/internal/platform/postgres"
	"github.com/HenryNg101/servers-ping-checker/internal/server"
	"github.com/elastic/go-elasticsearch/v9"
	"gorm.io/gorm"
)

type Application struct {
	PostgresSession *gorm.DB
	ElasticSession  *elasticsearch.Client
	ServerService   server.Service
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

	serverRepo := server.NewRepository(postgresSession)
	elasticServerRepo := server.NewServerESRepository(esSession)
	serverService := server.NewService(serverRepo, elasticServerRepo)

	return &Application{
		PostgresSession: postgresSession,
		ElasticSession:  esSession,
		ServerService:   serverService,
	}, nil
}
