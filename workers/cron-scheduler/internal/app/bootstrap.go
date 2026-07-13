package app

import (
	"errors"

	"github.com/HenryNg101/cron-scheduler/internal/agent"
	"github.com/HenryNg101/cron-scheduler/internal/config"
	jobs "github.com/HenryNg101/cron-scheduler/internal/jobs"
	"github.com/HenryNg101/cron-scheduler/internal/monitoring"
	"github.com/HenryNg101/cron-scheduler/internal/platform/blob_storage"
	"github.com/HenryNg101/cron-scheduler/internal/platform/elastic"
	"github.com/HenryNg101/cron-scheduler/internal/platform/mailer"
	"github.com/HenryNg101/cron-scheduler/internal/platform/postgres"
	"github.com/HenryNg101/cron-scheduler/internal/server"
	"github.com/HenryNg101/cron-scheduler/internal/user"
	"github.com/elastic/go-elasticsearch/v9"
	"gorm.io/gorm"
)

// TODO: Add Kafka dependency here somehow
type Application struct {
	PostgresSession   *gorm.DB
	ElasticSession    *elasticsearch.Client
	UserService       user.Service
	MonitoringService monitoring.Service
	JobService        jobs.Service
}

// Create a new app with services to be used
// Why I do this ? So that, when an app is created, it can use all internal functionalities without having to use other ways of communication like gRPC or HTTP
// If it's API app, you can add handlers, add API group, etc. to it. If it's a worker, it can still get access to DB query, services, etc. without it being separated entirely
// It might be tightly coupled, but it's good for now
func NewApp() (*Application, error) {
	// Load configs
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

	minIOConfig := config.LoadMinIO()
	if len(minIOConfig.AccessKey) == 0 {
		return nil, errors.New("Access key for MinIO is not configured. You have to set it in .env file in root folder using MINIO_ACCESS_KEY variable")
	}
	minIOInternalSession := blob_storage.NewMinIOSession(
		minIOConfig.InternalEndpoint,
		minIOConfig.AccessKey,
		minIOConfig.SecretKey,
		minIOConfig.UseSSL,
	)
	minIOExternalSession := blob_storage.NewMinIOSession(
		minIOConfig.PublicEndpoint,
		minIOConfig.AccessKey,
		minIOConfig.SecretKey,
		minIOConfig.UseSSL,
	)

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

	//
	// Create services
	serverRepo := server.NewRepository(postgresSession)
	elasticServerRepo := server.NewServerESRepository(esSession)

	userRepo := user.NewRepository(postgresSession)
	userService := user.NewService(userRepo)

	elasticAgentRepo := agent.NewAgentESRepository(esSession)
	monitoringService := monitoring.NewService(elasticAgentRepo, elasticServerRepo, serverRepo, mailerUtility)

	jobRepo := jobs.NewRepository(postgresSession)
	minIORepo := jobs.NewBlobStorage(minIOInternalSession, minIOExternalSession, minIOConfig.Bucket)
	jobService := jobs.NewService(jobRepo, minIORepo)

	return &Application{
		PostgresSession:   postgresSession,
		ElasticSession:    esSession,
		UserService:       userService,
		MonitoringService: monitoringService,
		JobService:        jobService,
	}, nil
}
