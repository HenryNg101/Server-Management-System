package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

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

		// Data streams only allow "create"
		meta := fmt.Sprintf(`{"create":{"_index":"server-metrics","_id":"%s"}}`, docID) + "\n"

		// marshal actual document
		docBytes, err := json.Marshal(m)
		if err != nil {
			return err
		}

		buf.WriteString(meta)
		buf.Write(docBytes)
		buf.WriteString("\n")
	}

	res, err := r.es.Bulk(
		bytes.NewReader(buf.Bytes()),
		r.es.Bulk.WithContext(ctx),
		r.es.Bulk.WithRefresh("wait_for"), // good for demo visibility
	)
	// res, err := r.es.Bulk(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.IsError() {
		return fmt.Errorf("bulk insert error: %s", string(body))
	}

	var bulkResp struct {
		Errors bool `json:"errors"`
		Items  []map[string]struct {
			Status int `json:"status"`
			Error  *struct {
				Type   string `json:"type"`
				Reason string `json:"reason"`
			} `json:"error,omitempty"`
		} `json:"items"`
	}

	if err := json.Unmarshal(body, &bulkResp); err != nil {
		return err
	}

	if bulkResp.Errors {
		for _, item := range bulkResp.Items {
			for action, result := range item {
				// 409 is normal if the same docID gets retried
				if result.Status == 409 {
					continue
				}
				if result.Error != nil {
					return fmt.Errorf("bulk item failed (%s): %s: %s", action, result.Error.Type, result.Error.Reason)
				}
			}
		}
		// return fmt.Errorf("bulk insert had errors")
	}

	return nil
}
