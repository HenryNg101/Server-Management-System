package data_transfer

import (
	"context"
	"fmt"
	"time"
)

const (
	batchSize = 500
)

type Service interface {
	// For clean up frequently (TTL stuff)
	CleanupOldImportJobs(ctx context.Context, olderThan time.Duration) error
}

type dataTransferService struct {
	repo        Repository
	blobStorage BlobStorage
}

func NewService(r Repository, s BlobStorage) Service {
	return &dataTransferService{repo: r, blobStorage: s}
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
