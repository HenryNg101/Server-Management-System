package server

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/HenryNg101/server-management-system/internal/model"
)

type Service interface {
	GetServers(ctx context.Context, q GetServersQuery) (*PaginatedServers, error)
	CreateServer(ctx context.Context, req CreateServerRequest) (*model.Server, error)
	GetServer(ctx context.Context, id uint, server *model.Server) (*model.Server, error)
	UpdateServer(ctx context.Context, id uint, req UpdateServerRequest) (*model.Server, error)
	DeleteServer(ctx context.Context, id uint) error
	ImportServers(ctx context.Context, r io.Reader) (*ImportServersResponse, error)
	BulkUpdateServersStatuses(ctx context.Context, servers []*model.Server) error

	// Elastic services
	ElasticBulkInsert(ctx context.Context, serversResults []*model.Server) error

	// Reporting
	SendReports(startTime time.Time, endTime time.Time, topN int, emailsList *[]string, ctx context.Context) (*Report, error)
}

type serverService struct {
	repo        Repository
	elasticRepo ElasticServerRepository
	mailUtility MailingUtility
}

func NewService(r Repository, elastic ElasticServerRepository, mailer MailingUtility) Service {
	return &serverService{repo: r, elasticRepo: elastic, mailUtility: mailer}
}

func (s *serverService) GetServers(ctx context.Context, q GetServersQuery) (*PaginatedServers, error) {
	servers, total, err := s.repo.FindAll(ctx, q)
	if err != nil {
		return nil, err
	}

	var totalPages int
	if q.PageSize != nil {
		totalPages = int(math.Ceil(float64(total) / float64(*q.PageSize)))
	}

	return &PaginatedServers{
		Servers:    servers,
		Total:      total,
		Page:       q.Page,
		PageSize:   q.PageSize,
		TotalPages: &totalPages,
	}, nil
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
	server := &model.Server{}

	server, err := s.repo.FindByID(ctx, id, server)
	if err != nil {
		return nil, err
	}

	// Apply updates ONLY if provided
	isUpdated := false
	if req.Name != nil {
		isUpdated = true
		server.Name = *req.Name
	}

	if req.Status != nil {
		isUpdated = true
		server.Status = *req.Status
	}

	if req.IPv4Address != nil {
		isUpdated = true
		server.IPv4Address = *req.IPv4Address
	}

	if req.Port != nil {
		isUpdated = true
		server.Port = *req.Port
	}

	if req.Protocol != nil {
		isUpdated = true
		server.Protocol = *req.Protocol
	}

	// If anything changes, it means that, it's actually updated
	// TODO: Handle the case where nothing is updated
	if !isUpdated {
		return nil, errors.New("Nothing is updated")
	}

	server.LastUpdated = time.Now()
	updated, err := s.repo.Update(ctx, server)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

// Bulk update statuses to be used by worker
func (s *serverService) BulkUpdateServersStatuses(ctx context.Context, servers []*model.Server) error {
	const chunkSize = 1000

	for i := 0; i < len(servers); i += chunkSize {
		end := i + chunkSize
		if end > len(servers) {
			end = len(servers)
		}

		err := s.repo.BulkUpdateStatus(ctx, servers[i:end])
		if err != nil {
			return err
		}
	}

	return nil
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

func (s *serverService) ElasticBulkInsert(ctx context.Context, serversResults []*model.Server) error {
	return s.elasticRepo.BulkInsertStatus(ctx, serversResults)
}

// TODO: Let the send email happens here
func (s *serverService) SendReports(startTime time.Time, endTime time.Time, topN int, emailsList *[]string, ctx context.Context) (*Report, error) {
	total, up, down, err := s.repo.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	uptime, err := s.elasticRepo.GetDailyUptime(ctx, startTime, endTime, topN)
	if err != nil {
		return nil, err
	}

	serversReport := &Report{
		TotalServers: total,
		ServersUp:    up,
		ServersDown:  down,
		Uptime:       uptime,
	}

	if emailsList == nil || len(*emailsList) == 0 {
		return serversReport, nil
	}

	mailHtmlContent := buildReportHTML(serversReport, startTime, endTime)

	subject := fmt.Sprintf(
		"Server Report (%s → %s)",
		startTime.Format("2006-01-02"),
		endTime.Add(-24*time.Hour).Format("2006-01-02"),
	)

	// TODO: Make this async -> No blocking of waiting for sending all emails
	err = s.mailUtility.Send(*emailsList, subject, mailHtmlContent)
	if err != nil {
		return nil, err
	}
	return serversReport, nil
}

func buildReportHTML(report *Report, startTime, endTime time.Time) string {
	var b strings.Builder

	b.WriteString("<h2>Server Report</h2>")
	b.WriteString(fmt.Sprintf("<p><b>Period:</b> %s → %s</p>",
		startTime.Format("2006-01-02"),
		endTime.Add(-24*time.Hour).Format("2006-01-02"),
	))

	// Summary
	b.WriteString("<h3>Summary</h3>")
	b.WriteString(fmt.Sprintf(`
	<ul>
		<li>Total Servers: %d</li>
		<li>Up: %d</li>
		<li>Down: %d</li>
	</ul>
	`, report.TotalServers, report.ServersUp, report.ServersDown))

	// Table
	b.WriteString("<h3>Uptime per Server</h3>")
	b.WriteString(`<table border="1" cellpadding="5" cellspacing="0">
	<tr><th>Server ID</th><th>Uptime</th></tr>`)

	for id, uptime := range report.Uptime {
		b.WriteString(fmt.Sprintf(
			"<tr><td>%d</td><td>%.2f%%</td></tr>",
			id, uptime*100,
		))
	}

	b.WriteString("</table>")

	return b.String()
}
