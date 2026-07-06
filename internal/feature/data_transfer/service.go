package data_transfer

import (
	"context"
	"encoding/csv"
	"errors"
	"io"

	"github.com/HenryNg101/server-management-system/internal/model"
)

const (
	batchSize = 500
)

type Service interface {
	ImportServers(ctx context.Context, r io.Reader) (*ImportServersResponse, error)
}

type dataTransferService struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &dataTransferService{repo: r}
}

// TODO: Handle more edge cases of uploading
// TODO: Improve performance of this API (It took nearly 8 seconds to loaded 10k records)
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
	err := s.repo.BulkUpsert(ctx, batch)
	if err == nil {
		return len(batch), nil
	}

	// Fallback when there's a failure: try inserting one by one to identify failures
	successCount := 0
	var failures []ImportFailure
	for i, srv := range batch {
		_, e := s.repo.Create(ctx, srv)
		if e != nil {
			failures = append(failures, ImportFailure{
				Row:   startRow + i,
				Error: e.Error(),
			})
			continue
		}
		successCount++
	}
	return successCount, failures
}
