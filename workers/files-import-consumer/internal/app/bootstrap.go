package app

import (
	"errors"

	"github.com/HenryNg101/files-import-consumer/internal/config"
	"github.com/HenryNg101/files-import-consumer/internal/jobs"
	"github.com/HenryNg101/files-import-consumer/internal/platform/blob_storage"
	"github.com/HenryNg101/files-import-consumer/internal/platform/postgres"
	"gorm.io/gorm"
)

type Application struct {
	PostgresSession *gorm.DB
	JobsService     jobs.Service
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
	minIORepo := jobs.NewBlobStorage(minIOInternalSession, minIOExternalSession, minIOConfig.Bucket)
	jobsService := jobs.NewService(jobsRepo, minIORepo)

	return &Application{
		PostgresSession: postgresSession,
		JobsService:     jobsService,
	}, nil
}
