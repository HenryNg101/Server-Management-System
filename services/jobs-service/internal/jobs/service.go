package jobs

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/HenryNg101/jobs-service/internal/model"
	"github.com/google/uuid"
)

const (
	batchSize = 500
)

type Service interface {
	// Handling jobs
	CreateImportJob(ctx context.Context, r io.Reader, fileSize int64) (*CreateImportJobResponse, error)
	GetJob(ctx context.Context, id string) (*model.ImportJob, error)
	GenerateFileDownloadURL(ctx context.Context, objectKey string) (string, error)
}

type dataTransferService struct {
	repo          Repository
	kafkaProducer KafkaProducer
	blobStorage   BlobStorage
}

func NewService(r Repository, k KafkaProducer, s BlobStorage) Service {
	return &dataTransferService{repo: r, kafkaProducer: k, blobStorage: s}
}

func (s *dataTransferService) GetJob(ctx context.Context, id string) (*model.ImportJob, error) {
	return s.repo.GetJob(ctx, id)
}

// Create new job for CSV file importing
func (s *dataTransferService) CreateImportJob(ctx context.Context, r io.Reader, fileSize int64) (*CreateImportJobResponse, error) {
	jobID := uuid.New().String() // Generate unique ID

	// upload to MinIO
	objectKey := fmt.Sprintf("jobs/%s/input.csv", jobID)
	_, err := s.blobStorage.Upload(ctx, objectKey, r, fileSize)
	if err != nil {
		return nil, err
	}

	// Create new job in the DB
	job := &model.ImportJob{
		ID:       jobID,
		FilePath: objectKey,
		Status:   model.JobStatusPending,
	}
	if err := s.repo.CreateJob(ctx, job); err != nil {
		return nil, err
	}

	// Publish Kafka message
	if err := s.kafkaProducer.PublishImportJob(ctx, jobID, objectKey); err != nil {
		return nil, err
	}
	return &CreateImportJobResponse{
		JobID:  jobID,
		Status: string(model.JobStatusPending),
	}, nil
}

func (s *dataTransferService) GenerateFileDownloadURL(ctx context.Context, objectKey string) (string, error) {
	// e.g. 15 minutes expiry
	return s.blobStorage.GetPresignedURL(ctx, objectKey, 15*time.Minute)
}
