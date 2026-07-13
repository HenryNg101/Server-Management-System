package monitoring

import (
	"cmp"
	"context"
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/HenryNg101/cron-scheduler/internal/agent"
	"github.com/HenryNg101/cron-scheduler/internal/platform/mailer"
	"github.com/HenryNg101/cron-scheduler/internal/server"
)

type Service interface {
	GetServersOverview(ctx context.Context, start, end time.Time, topN int) ([]*ServerOverview, error)
	SendReports(startTime time.Time, endTime time.Time, topN int, emailsList *[]string, ctx context.Context) (*Report, error)
}

type monitoringService struct {
	agentElasticRepo  agent.ElasticAgentRepository
	serverRepo        server.Repository
	serverElasticRepo server.ElasticServerRepository
	mailUtility       mailer.MailingUtility
}

func NewService(
	agentElastic agent.ElasticAgentRepository,
	serverElastic server.ElasticServerRepository,
	serverRepo server.Repository,
	mailer mailer.MailingUtility,
) Service {
	return &monitoringService{
		agentElasticRepo:  agentElastic,
		serverElasticRepo: serverElastic,
		serverRepo:        serverRepo,
		mailUtility:       mailer,
	}
}

// TODO: Add proper sorting for top N servers based on uptime and other metrics, which we needs to decide further on it, maybe some scoring would be needed
func (s *monitoringService) SendReports(
	startTime time.Time,
	endTime time.Time,
	topN int,
	emailsList *[]string,
	ctx context.Context,
) (*Report, error) {

	total, up, down, err := s.serverRepo.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	overview, err := s.GetServersOverview(ctx, startTime, endTime, topN)
	if err != nil {
		return nil, err
	}

	report := &Report{
		TotalServers: total,
		ServersUp:    up,
		ServersDown:  down,
		Stats:        overview,
	}

	if emailsList == nil || len(*emailsList) == 0 {
		return report, nil
	}

	html := buildReportHTML(report, startTime, endTime)

	subject := fmt.Sprintf(
		"Server Report (%s → %s)",
		startTime.Format("2006-01-02"),
		endTime.Format("2006-01-02"),
	)

	// TODO: Add async here, could be through Kafka or some sort of job system
	go func() {
		_ = s.mailUtility.Send(*emailsList, subject, html)
	}()

	return report, nil
}

func (s *monitoringService) GetServersOverview(ctx context.Context, start, end time.Time, topN int) ([]*ServerOverview, error) {
	metrics, err := s.agentElasticRepo.GetServersStats(ctx, start, end, topN)
	if err != nil {
		return nil, err
	}

	uptime, err := s.serverElasticRepo.GetDailyUptime(ctx, start, end, topN)
	if err != nil {
		return nil, err
	}

	result := make(map[uint]*ServerOverview)

	for id, u := range uptime {
		result[id] = &ServerOverview{ServerID: id}
		result[id].ServerPullStats = *u
	}
	for id, m := range metrics {
		if _, ok := result[id]; !ok {
			result[id] = &ServerOverview{ServerID: id}
		}
		result[id].ServerPushStats = *m
	}

	// Sorting the result by uptime (worst first) and limiting to topN
	values := slices.Collect(maps.Values(result))
	slices.SortFunc(values, func(a, b *ServerOverview) int {
		return cmp.Compare(a.Uptime, b.Uptime) // worst first
	})
	if len(values) > topN {
		values = values[:topN]
	}
	return values, nil
}
