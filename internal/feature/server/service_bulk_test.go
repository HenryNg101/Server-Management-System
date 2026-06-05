package server

import (
	"context"
	"testing"

	"github.com/HenryNg101/server-management-system/internal/model"
)

func TestBulkUpdateServersStatuses(t *testing.T) {
	svc, repo := setupService()

	var servers []*model.Server

	for i := 0; i < 1500; i++ {
		s, _ := repo.Create(context.Background(), &model.Server{
			Status: false,
		})
		servers = append(servers, &model.Server{
			ID:     s.ID,
			Status: true,
		})
	}

	err := svc.BulkUpdateServersStatuses(context.Background(), servers)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, s := range servers {
		if !s.Status {
			t.Fatalf("expected all servers to be updated")
		}
	}
}

func TestBulkUpdate_EmptyList(t *testing.T) {
	svc, _ := setupService()

	err := svc.BulkUpdateServersStatuses(context.Background(), []*model.Server{})
	if err != nil {
		t.Fatalf("should not fail on empty input")
	}
}

func TestBulkUpdate_SmallChunk(t *testing.T) {
	svc, repo := setupService()

	s, _ := repo.Create(context.Background(), &model.Server{Status: false})

	err := svc.BulkUpdateServersStatuses(context.Background(), []*model.Server{
		{ID: s.ID, Status: true},
	})

	if err != nil {
		t.Fatalf("unexpected error")
	}

	if !repo.servers[s.ID].Status {
		t.Fatalf("status not updated")
	}
}
