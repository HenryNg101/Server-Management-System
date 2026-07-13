package server

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/HenryNg101/servers-ping-checker/internal/model"
	"github.com/elastic/go-elasticsearch/v9"
)

type ElasticServerRepository interface {
	BulkInsertStatus(ctx context.Context, results []*model.Server) error
}

type elasticServerRepository struct {
	es *elasticsearch.Client
}

func NewServerESRepository(es *elasticsearch.Client) ElasticServerRepository {
	return &elasticServerRepository{es: es}
}

// TODO: Finish this function
func (r *elasticServerRepository) BulkInsertStatus(ctx context.Context, results []*model.Server) error {
	var buf bytes.Buffer

	for _, r := range results {
		meta := `{"create":{"_index":"server-status"}}` + "\n"

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
