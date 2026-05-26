package server

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"net"
	"strconv"

	"github.com/HenryNg101/server-management-system/internal/model"
)

type Service interface {
	GetServers(ctx context.Context) ([]model.Server, error)
	CreateServer(ctx context.Context, req CreateServerRequest) (*model.Server, error)
	GetServer(ctx context.Context, id uint, server *model.Server) (*model.Server, error)
	UpdateServer(ctx context.Context, id uint, req UpdateServerRequest) (*model.Server, error)
	DeleteServer(ctx context.Context, id uint) error
	ImportServers(ctx context.Context, r io.Reader) (*ImportServersResponse, error)
}

type serverService struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &serverService{repo: r}
}

func (s *serverService) GetServers(ctx context.Context) ([]model.Server, error) {
	return s.repo.FindAll(ctx)
}

func (s *serverService) CreateServer(ctx context.Context, req CreateServerRequest) (*model.Server, error) {
	server := model.Server{
		Name:        req.Name,
		Status:      req.Status,
		IPv4Address: req.IPv4Address,
		Port:        req.Port,
		Protocol:    req.Protocol,
	}
	return s.repo.Create(ctx, &server)
}

func (s *serverService) GetServer(ctx context.Context, id uint, server *model.Server) (*model.Server, error) {
	return s.repo.FindByID(ctx, id, server)
}

func (s *serverService) UpdateServer(ctx context.Context, id uint, req UpdateServerRequest) (*model.Server, error) {
	var server *model.Server

	server, err := s.repo.FindByID(ctx, id, server)
	if err != nil {
		return nil, err
	}

	// Apply updates ONLY if provided
	if req.Name != nil {
		server.Name = *req.Name
	}

	if req.Status != nil {
		server.Status = *req.Status
	}

	if req.IPv4Address != nil {
		server.IPv4Address = *req.IPv4Address
	}

	if req.Port != nil {
		server.Port = *req.Port
	}

	if req.Protocol != nil {
		server.Protocol = *req.Protocol
	}

	updated, err := s.repo.Update(ctx, server)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

var ErrNotFound = errors.New("server not found")

func (s *serverService) DeleteServer(ctx context.Context, id uint) error {
	exists, err := s.repo.ExistsByID(ctx, id)
	if err != nil {
		return err
	}

	if !exists {
		return ErrNotFound
	}

	return s.repo.Delete(ctx, id)
}

// TODO: Handle more edge cases of uploading
// TODO: Improve performance of this API (It took nearly 8 seconds to loaded 10k records)
func (s *serverService) ImportServers(ctx context.Context, r io.Reader) (*ImportServersResponse, error) {
	reader := csv.NewReader(r)

	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, errors.New("empty csv")
	}

	headers := rows[0]

	var successes []model.Server
	var failures []ImportFailure

	for i, row := range rows[1:] {
		record := mapRow(headers, row)

		server, err := parseServer(record)
		if err != nil {
			failures = append(failures, ImportFailure{
				Row:    i + 2,
				Error:  err.Error(),
				Record: record,
			})
			continue
		}

		created, err := s.repo.Create(ctx, server)
		if err != nil {
			failures = append(failures, ImportFailure{
				Row:    i + 2,
				Error:  err.Error(),
				Record: record,
			})
			continue
		}

		if created != nil {
			successes = append(successes, *created)
		}
	}

	return &ImportServersResponse{
		SuccessCount: len(successes),
		FailedCount:  len(failures),
		Successes:    successes,
		Failures:     failures,
	}, nil
}

func mapRow(headers, row []string) map[string]string {
	m := make(map[string]string)
	for i := range headers {
		m[headers[i]] = row[i]
	}
	return m
}

func parseServer(record map[string]string) (*model.Server, error) {
	name := record["name"]
	if name == "" {
		return nil, errors.New("name is required")
	}

	status, err := strconv.ParseBool(record["status"])
	if err != nil {
		return nil, errors.New("invalid status")
	}

	ip := record["ipv4_address"]
	if net.ParseIP(ip) == nil {
		return nil, errors.New("invalid IP")
	}

	port, err := strconv.Atoi(record["port"])
	if err != nil || port < 0 || port > 65535 {
		return nil, errors.New("invalid port")
	}

	protocol := record["protocol"]
	if protocol == "" {
		protocol = "tcp"
	}

	return &model.Server{
		Name:        name,
		Status:      status,
		IPv4Address: ip,
		Port:        uint(port),
		Protocol:    protocol,
	}, nil
}
