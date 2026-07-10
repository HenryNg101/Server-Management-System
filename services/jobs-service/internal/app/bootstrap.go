package app

import (
	"errors"

	// jobs "command-line-argumentsC:\\Users\\hungnq77\\Documents\\Server-Management-System\\services\\jobs-service\\internal\\jobs\\kafka_producer.go"
	"github.com/HenryNg101/jobs-service/internal/config"
	"github.com/HenryNg101/jobs-service/internal/jobs"
	"github.com/HenryNg101/jobs-service/internal/platform/blob_storage"
	"github.com/HenryNg101/jobs-service/internal/platform/elastic"
	kafkaClient "github.com/HenryNg101/jobs-service/internal/platform/kafka"
	"github.com/HenryNg101/jobs-service/internal/platform/postgres"
	"github.com/elastic/go-elasticsearch/v9"
	"gorm.io/gorm"
)

type Application struct {
	PostgresSession *gorm.DB
	ElasticSession  *elasticsearch.Client
	// MinIOClient     *minio.Client
	JobsService jobs.Service
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

	jobsRepo := jobs.NewRepository(postgresSession)
	kafkaImportProducer := jobs.NewKafkaProducer(kafkaClient.NewProducer(kafkaConfig.Brokers, kafkaConfig.ServersImportTopic))
	minIORepo := jobs.NewBlobStorage(minIOInternalSession, minIOExternalSession, minIOConfig.Bucket)
	jobsService := jobs.NewService(jobsRepo, kafkaImportProducer, minIORepo)

	return &Application{
		PostgresSession: postgresSession,
		ElasticSession:  esSession,
		JobsService:     jobsService,
	}, nil
}
