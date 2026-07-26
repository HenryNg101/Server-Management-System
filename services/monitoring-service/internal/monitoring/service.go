package monitoring

import (
	"cmp"
	"context"
	"maps"
	"slices"
	"time"

	"github.com/HenryNg101/monitoring-service/internal/agent"
	"github.com/HenryNg101/monitoring-service/internal/platform/mailer"
	"github.com/HenryNg101/monitoring-service/internal/server"
)

type Service interface {
	GetServersOverview(ctx context.Context, start, end time.Time, topN int) ([]*ServerOverview, error)
	SendReports(startTime time.Time, endTime time.Time, topN int, emailsList *[]string, ctx context.Context) (*Report, error)
}

type monitoringService struct {
	agentElasticRepo agent.ElasticAgentRepository
	serverRepo       server.Repository
	mailUtility      mailer.MailingUtility
}

func NewService(
	agentElastic agent.ElasticAgentRepository,
	serverRepo server.Repository,
	mailer mailer.MailingUtility,
) Service {
	return &monitoringService{
		agentElasticRepo: agentElastic,
		serverRepo:       serverRepo,
		mailUtility:      mailer,
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

	total, err := s.serverRepo.GetStats(ctx)
	if err != nil {
		return nil, err
	}

	overview, err := s.GetServersOverview(ctx, startTime, endTime, topN)
	if err != nil {
		return nil, err
	}

	// In push-only mode:
	// - any server that has agent data in the window is treated as "up"
	// - servers with no data are "down / no data"
	up := int64(len(overview))
	down := total - up
	if down < 0 {
		down = 0
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

	// TODO: Add async here, could be through Kafka or some sort of job system
	go func() {
		_ = s.mailUtility.Send(*emailsList, "Server Statuses Report", html)
	}()

	return report, nil
}

func (s *monitoringService) GetServersOverview(ctx context.Context, start, end time.Time, topN int) ([]*ServerOverview, error) {
	metrics, err := s.agentElasticRepo.GetServersStats(ctx, start, end, topN)
	if err != nil {
		return nil, err
	}

	result := make(map[uint]*ServerOverview)

	for id, m := range metrics {
		if _, ok := result[id]; !ok {
			result[id] = &ServerOverview{ServerID: id}
		}
		result[id].ServerPushStats = *m
	}

	// Sorting the result by uptime (worst first) and limiting to topN
	values := slices.Collect(maps.Values(result))
	slices.SortFunc(values, func(a, b *ServerOverview) int {
		return cmp.Compare(a.ServerPushStats.Uptime, b.ServerPushStats.Uptime) // worst first
	})
	if len(values) > topN {
		values = values[:topN]
	}
	return values, nil
}
