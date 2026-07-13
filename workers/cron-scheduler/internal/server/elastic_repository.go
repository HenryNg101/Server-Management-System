package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
)

type ElasticServerRepository interface {
	// BulkInsertStatus(ctx context.Context, results []*model.Server) error
	GetDailyUptime(ctx context.Context, startTime time.Time, endTime time.Time, topN int) (map[uint]*ServerPullStats, error)
}

type elasticServerRepository struct {
	es *elasticsearch.Client
}

func NewServerESRepository(es *elasticsearch.Client) ElasticServerRepository {
	return &elasticServerRepository{es: es}
}

func (r *elasticServerRepository) GetDailyUptime(ctx context.Context, startTime time.Time, endTime time.Time, topN int) (map[uint]*ServerPullStats, error) {
	query := fmt.Sprintf(`{
	  "size": 0,
	  "query": {
	    "range": {
	      "@timestamp": {
	        "gte": "%s",
	        "lt": "%s"
	      }
	    }
	  },
	  "aggs": {
	    "servers": {
	      "terms": {
	        "field": "server_id",
			"size": %d,
			"order": {
				"uptime": "asc"
			}
	      },
	      "aggs": {
	        "uptime": {
	          "avg": {
	            "field": "status"
	          }
	        }
	      }
	    }
	  }
	}`, startTime.Format(time.RFC3339), endTime.Format(time.RFC3339), topN)

	res, err := r.es.Search(
		r.es.Search.WithContext(ctx),
		r.es.Search.WithIndex("server-status"),
		r.es.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("ES error: %s", res.String())
	}

	//
	// Parse the result, only extract what we need
	var parsed map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	result := make(map[uint]*ServerPullStats)

	aggs, ok := parsed["aggregations"].(map[string]interface{})
	if !ok {
		return result, nil
	}

	buckets := aggs["servers"].(map[string]interface{})["buckets"].([]interface{})

	for _, b := range buckets {
		bucket := b.(map[string]interface{})

		serverID := uint(bucket["key"].(float64))
		uptime := bucket["uptime"].(map[string]interface{})["value"].(float64)

		result[serverID] = &ServerPullStats{
			Uptime: uptime,
		}
	}

	return result, nil
}
