package app

import (
	"errors"

	"github.com/HenryNg101/monitoring-service/internal/agent"
	"github.com/HenryNg101/monitoring-service/internal/config"
	"github.com/HenryNg101/monitoring-service/internal/monitoring"
	"github.com/HenryNg101/monitoring-service/internal/platform/elastic"
	"github.com/HenryNg101/monitoring-service/internal/platform/mailer"
	"github.com/HenryNg101/monitoring-service/internal/platform/postgres"
	"github.com/HenryNg101/monitoring-service/internal/server"
	"github.com/elastic/go-elasticsearch/v9"
	"gorm.io/gorm"
)

type Application struct {
	PostgresSession   *gorm.DB
	ElasticSession    *elasticsearch.Client
	MonitoringService monitoring.Service
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

	mailerConfig := config.LoadMailer()
	if len(mailerConfig.Password) == 0 {
		return nil, errors.New("Password for email server is not set. You have to set it in .env file in root folder using MAIL_PASSWORD variable")
	}
	mailerUtility := mailer.NewMailer(
		mailerConfig.Server,
		mailerConfig.Port,
		mailerConfig.UserName,
		mailerConfig.Password,
		mailerConfig.FromEmail,
	)

	serverRepo := server.NewRepository(postgresSession)

	elasticAgentRepo := agent.NewAgentESRepository(esSession)

	monitoringService := monitoring.NewService(elasticAgentRepo, serverRepo, mailerUtility)

	return &Application{
		PostgresSession:   postgresSession,
		ElasticSession:    esSession,
		MonitoringService: monitoringService,
	}, nil
}
