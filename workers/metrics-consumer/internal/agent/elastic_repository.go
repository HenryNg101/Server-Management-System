package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v9"
)

type ElasticAgentRepository interface {
	BulkInsertStatus(ctx context.Context, results []MetricMessage) error
}

type elasticAgentRepository struct {
	es *elasticsearch.Client
}

func NewAgentESRepository(es *elasticsearch.Client) ElasticAgentRepository {
	return &elasticAgentRepository{es: es}
}

func (r *elasticAgentRepository) BulkInsertStatus(ctx context.Context, messages []MetricMessage) error {
	var buf bytes.Buffer

	for _, m := range messages {
		// Deterministic document IDs, to prevent duplicates in case of Kafka retries. Format: serverID-containerName-timestamp
		docID := fmt.Sprintf("%d-%s-%d",
			m.ServerID,
			m.ContainerName,
			m.Timestamp.UnixNano(),
		)
		meta := fmt.Sprintf(`{"index":{"_index":"server-metrics","_id":"%s"}}`, docID) + "\n"

		// marshal actual document
		docBytes, err := json.Marshal(m)
		if err != nil {
			return err
		}

		buf.WriteString(meta)
		buf.Write(docBytes)
		buf.WriteString("\n")
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
