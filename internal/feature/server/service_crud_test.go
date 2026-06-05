package server

import (
	"context"
	"testing"

	"github.com/HenryNg101/server-management-system/internal/model"
)

func setupService() (Service, *mockRepo) {
	repo := NewFakeRepo()
	svc := NewService(repo, nil, nil)
	return svc, repo
}

func TestCreateServer(t *testing.T) {
	svc, _ := setupService()

	res, err := svc.CreateServer(context.Background(), CreateServerRequest{
		Name:        "test",
		Status:      true,
		IPv4Address: "127.0.0.1",
		Port:        80,
		Protocol:    "tcp",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.ID == 0 {
		t.Fatalf("expected ID to be set")
	}
}

func TestGetServer(t *testing.T) {
	svc, repo := setupService()

	created, _ := repo.Create(context.Background(), &model.Server{
		Name: "test",
	})

	res, err := svc.GetServer(context.Background(), created.ID, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.ID != created.ID {
		t.Fatalf("expected same ID")
	}
}

func TestGetServers(t *testing.T) {
	svc, repo := setupService()

	repo.Create(context.Background(), &model.Server{Name: "a"})
	repo.Create(context.Background(), &model.Server{Name: "b"})

	page := 1
	size := 10

	res, err := svc.GetServers(context.Background(), GetServersQuery{
		Page:     &page,
		PageSize: &size,
	})

	if err != nil {
		t.Fatalf("unexpected error")
	}

	if res.Total != 2 {
		t.Fatalf("expected total 2")
	}
}

func TestGetServer_NotFound(t *testing.T) {
	svc, _ := setupService()

	_, err := svc.GetServer(context.Background(), 999, nil)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestGetServers_FilteringAndPagination(t *testing.T) {
	repo := NewFakeRepo()

	err := repo.Seed(
		&model.Server{Name: "alpha", Status: true, Protocol: "tcp"},
		&model.Server{Name: "beta", Status: false, Protocol: "udp"},
		&model.Server{Name: "gamma", Status: true, Protocol: "tcp"},
	)

	svc := NewService(repo, nil, nil)

	status := true
	protocol := "tcp"
	page := 1
	size := 1

	res, err := svc.GetServers(context.Background(), GetServersQuery{
		Status:   &status,
		Protocol: &protocol,
		Page:     &page,
		PageSize: &size,
		SortBy:   "name",
		Order:    "asc",
	})

	if err != nil {
		t.Fatalf("unexpected error")
	}

	if res.Total != 2 {
		t.Fatalf("expected 2 filtered results")
	}

	if len(res.Servers) != 1 {
		t.Fatalf("expected pagination to limit results")
	}
}

func TestUpdateServer(t *testing.T) {
	svc, repo := setupService()

	created, _ := repo.Create(context.Background(), &model.Server{
		Name:   "old",
		Status: true,
	})

	newName := "new"

	res, err := svc.UpdateServer(context.Background(), created.ID, UpdateServerRequest{
		Name: &newName,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.Name != "new" {
		t.Fatalf("expected name updated")
	}
}

func TestUpdateServer_NoChanges(t *testing.T) {
	svc, repo := setupService()

	created, _ := repo.Create(context.Background(), &model.Server{})

	_, err := svc.UpdateServer(context.Background(), created.ID, UpdateServerRequest{})
	if err == nil {
		t.Fatalf("expected error when nothing updated")
	}
}

func TestUpdateServer_AllFields(t *testing.T) {
	svc, repo := setupService()

	created, _ := repo.Create(context.Background(), &model.Server{
		Name:        "old",
		Status:      false,
		IPv4Address: "127.0.0.1",
		Port:        80,
		Protocol:    "tcp",
	})

	name := "new"
	status := true
	ip := "192.168.1.1"
	port := uint(443)
	protocol := "udp"

	res, err := svc.UpdateServer(context.Background(), created.ID, UpdateServerRequest{
		Name:        &name,
		Status:      &status,
		IPv4Address: &ip,
		Port:        &port,
		Protocol:    &protocol,
	})

	if err != nil {
		t.Fatalf("unexpected error")
	}

	if res.Name != name || res.Protocol != protocol || res.Port != port {
		t.Fatalf("fields not updated correctly")
	}
}

func TestUpdateServer_NotFound(t *testing.T) {
	svc, _ := setupService()

	name := "new"

	_, err := svc.UpdateServer(context.Background(), 999, UpdateServerRequest{
		Name: &name,
	})

	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestDeleteServer(t *testing.T) {
	svc, repo := setupService()

	created, _ := repo.Create(context.Background(), &model.Server{})

	err := svc.DeleteServer(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteServer_NotFound(t *testing.T) {
	svc, _ := setupService()

	err := svc.DeleteServer(context.Background(), 999)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestDeleteServer_RepoError(t *testing.T) {
	repo := NewFakeRepo()
	repo.SetDbExists(false)

	svc := NewService(repo, nil, nil)

	err := svc.DeleteServer(context.Background(), 1)
	if err == nil {
		t.Fatalf("expected error")
	}
}
