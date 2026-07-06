package data_transfer

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/HenryNg101/server-management-system/internal/model"
	"github.com/google/uuid"
)

const (
	batchSize = 500
)

type Service interface {
	ImportServers(ctx context.Context, r io.Reader) (*ImportServersResponse, error)

	// Handling jobs
	CreateImportJob(ctx context.Context, r io.Reader) (string, error)
	ProcessImportJob(ctx context.Context, msg ImportJobMessage) error
}

type dataTransferService struct {
	repo          Repository
	kafkaProducer KafkaProducer
	blobStorage   BlobStorage
}

func NewService(r Repository, k KafkaProducer, s BlobStorage) Service {
	return &dataTransferService{repo: r, kafkaProducer: k, blobStorage: s}
}

// Create new job for CSV file importing
func (s *dataTransferService) CreateImportJob(ctx context.Context, r io.Reader) (string, error) {
	jobID := uuid.New().String() // Generate unique ID

	// upload to MinIO
	objectKey := fmt.Sprintf("imports/%s.csv", jobID)
	_, err := s.blobStorage.Upload(ctx, objectKey, r, -1)
	if err != nil {
		return "", err
	}

	// Create new job in the DB
	job := &model.ImportJob{
		ID:       jobID,
		FilePath: objectKey,
		Status:   model.JobStatusPending,
	}
	if err := s.repo.CreateJob(ctx, job); err != nil {
		return "", err
	}

	// Publish Kafka message
	msg := ImportJobMessage{
		JobID:     jobID,
		ObjectKey: objectKey,
	}
	payload, _ := json.Marshal(msg)
	if err := s.kafkaProducer.WriteOne(ctx, []byte(jobID), payload); err != nil {
		return "", err
	}
	return jobID, nil
}

// Process a job
func (s *dataTransferService) ProcessImportJob(ctx context.Context, msg ImportJobMessage) error {
	_ = s.repo.UpdateJobStatus(ctx, msg.JobID, string(model.JobStatusProcessing), "")

	reader, err := s.blobStorage.Download(ctx, msg.ObjectKey)
	if err != nil {
		s.repo.UpdateJobStatus(ctx, msg.JobID, string(model.JobStatusFailed), err.Error())
		return err
	}
	defer reader.Close()

	_, err = s.ImportServers(ctx, reader)
	if err != nil {
		s.repo.UpdateJobStatus(ctx, msg.JobID, string(model.JobStatusFailed), err.Error())
		return err
	}

	return s.repo.UpdateJobStatus(ctx, msg.JobID, string(model.JobStatusDone), "")
}

// TODO: Handle more edge cases of uploading
func (s *dataTransferService) ImportServers(ctx context.Context, r io.Reader) (*ImportServersResponse, error) {
	reader := csv.NewReader(r)

	// Read headers
	headers, err := reader.Read()
	if err != nil {
		return nil, errors.New("failed to read CSV headers")
	}

	var (
		batch        []*model.Server
		successCount int
		failures     []ImportFailure
		rowIndex     = 1 // header is row 1
	)

	// Read line by line instead of whole file -> Not make this a blocking operation
	for {
		rowIndex++

		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			failures = append(failures, ImportFailure{
				Row:   rowIndex,
				Error: err.Error(),
			})
			continue
		}

		record := mapRow(headers, row)

		server, err := parseServer(record)
		if err != nil {
			failures = append(failures, ImportFailure{
				Row:   rowIndex,
				Error: err.Error(),
			})
			continue
		}

		batch = append(batch, server)

		// Flush batch
		if len(batch) >= batchSize {
			inserted, failed := s.insertBatch(ctx, batch, rowIndex-len(batch)+1)
			successCount += inserted
			failures = append(failures, failed...)

			batch = batch[:0]
		}
	}

	// Final flush
	if len(batch) > 0 {
		inserted, failed := s.insertBatch(ctx, batch, rowIndex-len(batch)+1)
		successCount += inserted
		failures = append(failures, failed...)
	}

	return &ImportServersResponse{
		SuccessCount: successCount,
		FailedCount:  len(failures),
		Failures:     failures,
	}, nil
}

func (s *dataTransferService) insertBatch(ctx context.Context, batch []*model.Server, startRow int) (int, []ImportFailure) {
	err := s.repo.BulkUpsertServers(ctx, batch)
	if err == nil {
		return len(batch), nil
	}

	// Fallback when there's a failure: try inserting one by one to identify failures
	successCount := 0
	var failures []ImportFailure
	for i, srv := range batch {
		err := s.repo.CreateServer(ctx, srv)
		if err != nil {
			failures = append(failures, ImportFailure{
				Row:   startRow + i,
				Error: err.Error(),
			})
			continue
		}
		successCount++
	}
	return successCount, failures
}
