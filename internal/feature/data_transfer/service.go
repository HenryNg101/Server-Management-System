package data_transfer

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/HenryNg101/server-management-system/internal/model"
	"github.com/google/uuid"
)

const (
	batchSize = 500
)

type Service interface {
	// importServers(ctx context.Context, r io.Reader) (*ImportServersResponse, error)
	importServers(ctx context.Context, r io.Reader, progressCallback func(processed, success, failed int)) ([]ImportFailure, error)

	// Handling jobs
	CreateImportJob(ctx context.Context, r io.Reader, fileSize int64) (*CreateImportJobResponse, error)
	ProcessImportJob(ctx context.Context, msg ImportJobMessage) error
	GetJob(ctx context.Context, id string) (*model.ImportJob, error)
	GenerateFileDownloadURL(ctx context.Context, objectKey string) (string, error)

	// For clean up frequently (TTL stuff)
	CleanupOldImportJobs(ctx context.Context, olderThan time.Duration) error
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

// Process a job
func (s *dataTransferService) ProcessImportJob(ctx context.Context, msg ImportJobMessage) error {
	job, err := s.GetJob(ctx, msg.JobID)
	if err != nil {
		return fmt.Errorf("failed to get job information: %v", err.Error())
	}

	// Ensure idempotency, in case where processing is done, but commit isn't done yet, and system crashed mid-way
	if job.Status == model.JobStatusDone {
		return nil
	}

	if err := s.repo.UpdateJobStatus(ctx, msg.JobID, model.JobStatusProcessing, nil, nil); err != nil {
		// log.Printf("failed to update job status: %v", err)
		return fmt.Errorf("failed to update job status: %v", err.Error())
	}

	// Download the file from MinIO and check if it's possible
	reader, err := s.blobStorage.Download(ctx, msg.ObjectKey)
	if err != nil {
		errMsg := err.Error()
		s.repo.UpdateJobStatus(ctx, msg.JobID, model.JobStatusFailed, &errMsg, nil)
		return err
	}
	defer reader.Close()

	// Pass progress callback into importServers
	failedRows, err := s.importServers(ctx, reader, func(processed, success, failed int) {
		_ = s.repo.UpdateJobProgress(ctx, msg.JobID, processed, success, failed)
	})

	if err != nil {
		errMsg := err.Error()
		s.repo.UpdateJobStatus(ctx, msg.JobID, model.JobStatusFailed, &errMsg, nil)
		return err
	}

	// Upload specific info on failed rows into a json ONCE
	resultPath := fmt.Sprintf("jobs/%s/failures.json", msg.JobID)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(failedRows); err != nil {
		errMsg := err.Error()
		s.repo.UpdateJobStatus(ctx, msg.JobID, model.JobStatusFailed, &errMsg, nil)
		return err
	}

	_, err = s.blobStorage.Upload(ctx, resultPath, &buf, int64(buf.Len()))
	if err != nil {
		errMsg := err.Error()
		s.repo.UpdateJobStatus(ctx, msg.JobID, model.JobStatusFailed, &errMsg, nil)
		return err
	}

	return s.repo.UpdateJobStatus(ctx, msg.JobID, model.JobStatusDone, nil, &resultPath)
}

func (s *dataTransferService) importServers(ctx context.Context, r io.Reader, progressCallback func(processed, success, failed int)) ([]ImportFailure, error) {
	reader := csv.NewReader(r)

	headers, err := reader.Read()
	if err != nil {
		return nil, errors.New("failed to read CSV headers")
	}

	var (
		batch         []*model.Server
		failedRows    []ImportFailure
		rowIndex      = 1
		processedRows int
		successCount  int
		failedCount   int
	)

	for {
		rowIndex++

		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		// Processed rows counts only represents the amount of rows it reads so far, so as long as there's a line that's read (Even if invalid), it's counted
		processedRows++

		// Validate parsing
		if err != nil {
			failedRows = append(failedRows, ImportFailure{
				Row:   rowIndex,
				Error: err.Error(),
			})
			failedCount++
			continue
		}

		record := mapRow(headers, row)

		server, err := parseServer(record)
		if err != nil {
			failedRows = append(failedRows, ImportFailure{
				Row:   rowIndex,
				Error: err.Error(),
			})
			failedCount++
			continue
		}

		// Push to batch after it's validated, and batch insert each time it reaches the size
		batch = append(batch, server)

		if len(batch) >= batchSize {
			insertedCount, failed := s.insertBatch(ctx, batch, rowIndex-len(batch)+1)

			successCount += insertedCount
			failedCount += len(failed)
			failedRows = append(failedRows, failed...)

			batch = batch[:0]

			if progressCallback != nil {
				progressCallback(processedRows, successCount, failedCount)
			}
		}
	}

	// final flush
	if len(batch) > 0 {
		insertedCount, failed := s.insertBatch(ctx, batch, rowIndex-len(batch)+1)

		successCount += insertedCount
		failedCount += len(failed)
		failedRows = append(failedRows, failed...)
	}

	if progressCallback != nil {
		progressCallback(processedRows, successCount, failedCount)
	}
	return failedRows, nil
}

func (s *dataTransferService) insertBatch(ctx context.Context, batch []*model.Server, startRow int) (int, []ImportFailure) {
	err := s.repo.BulkUpsertServers(ctx, batch)
	if err == nil {
		return len(batch), nil
	}

	// Fallback when there's a failure: try inserting one by one to identify failures
	successCount := 0
	var failedRows []ImportFailure
	for i, srv := range batch {
		err := s.repo.CreateServer(ctx, srv)
		if err != nil {
			failedRows = append(failedRows, ImportFailure{
				Row:   startRow + i,
				Error: err.Error(),
			})
			continue
		}
		successCount++
	}
	return successCount, failedRows
}

func (s *dataTransferService) CleanupOldImportJobs(ctx context.Context, olderThan time.Duration) error {
	jobs, err := s.repo.GetOldJobs(ctx, olderThan)
	if err != nil {
		return err
	}

	for _, job := range jobs {
		// delete input file
		_ = s.blobStorage.Delete(ctx, job.FilePath)

		// delete result file if exists
		if job.ResultPath != nil {
			_ = s.blobStorage.Delete(ctx, *job.ResultPath)
		}

		// delete DB record
		if err := s.repo.DeleteJob(ctx, job.ID); err != nil {
			// Log and continue, without stopping the whole cleanup
			fmt.Printf("failed to delete job %s: %v\n", job.ID, err)
		}
	}

	return nil
}

func (s *dataTransferService) GenerateFileDownloadURL(ctx context.Context, objectKey string) (string, error) {
	// e.g. 15 minutes expiry
	return s.blobStorage.GetPresignedURL(ctx, objectKey, 15*time.Minute)
}
