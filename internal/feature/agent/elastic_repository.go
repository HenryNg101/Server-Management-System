package agent

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
)

type ElasticAgentRepository interface {
	BulkInsertStatus(ctx context.Context, results []MetricMessage) error
	GetDailyStats(ctx context.Context, startTime time.Time, endTime time.Time, topN int) (map[uint]float64, error)
}

type elasticAgentRepository struct {
	es *elasticsearch.Client
}

func NewServerESRepository(es *elasticsearch.Client) ElasticAgentRepository {
	return &elasticAgentRepository{es: es}
}

func (r *elasticAgentRepository) BulkInsertStatus(ctx context.Context, messages []MetricMessage) error {
	var buf bytes.Buffer

	for _, m := range messages {
		meta := `{"create":{"_index":"server-metrics"}}` + "\n"

		doc := fmt.Sprintf(
			`{"@timestamp":"%s","server_id":%d,"status":%t}`+"\n",
			r.LastUpdated.UTC().Format(time.RFC3339),
			r.ID,
			r.Status,
		)

		buf.WriteString(meta)
		buf.WriteString(doc)
	}

	res, err := r.es.Bulk(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk insert error: %s", res.String())
	}

	return nil
}

// TODO: Implement later
func (r *elasticAgentRepository) GetDailyStats(ctx context.Context, startTime time.Time, endTime time.Time, topN int) (map[uint]float64, error) {

}
