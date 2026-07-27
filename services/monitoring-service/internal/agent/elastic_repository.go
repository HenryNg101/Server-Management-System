package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
)

const reportBucketInterval = 30 * time.Second

type ElasticAgentRepository interface {
	GetServersStats(ctx context.Context, startTime time.Time, endTime time.Time, topN int, serverCreatedAt map[uint]time.Time) (map[uint]*ServerPushStats, error)
}

type elasticAgentRepository struct {
	es *elasticsearch.Client
}

func NewAgentESRepository(es *elasticsearch.Client) ElasticAgentRepository {
	return &elasticAgentRepository{es: es}
}

// Basically, the idea is to calculate sum/avg of the stats of all containers in each time bucket at first
// Then, do it the same, but sum/avg of the stats of all buckets (And yes, min/max for the case of OOM events)
// And finally, count the number of buckets, divided by the amount of expected buckets -> Uptime percentage
func (r *elasticAgentRepository) GetServersStats(
	ctx context.Context,
	startTime time.Time,
	endTime time.Time,
	topN int,
	serverCreatedAt map[uint]time.Time,
) (map[uint]*ServerPushStats, error) {
	query := fmt.Sprintf(`
		{
		"size": 0,
		"query": {
			"range": {
			"@timestamp": {
				"gte": "%s",
				"lte": "%s"
			}
			}
		},
		"aggs": {
			"servers": {
			"terms": {
				"field": "server_id",
				"size": %d
			},
			"aggs": {
				"per_time": {
					"date_histogram": {
						"field": "@timestamp",
						"fixed_interval": "%s",
						"min_doc_count": 1
					},
					"aggs": {
						"cpu_usage_sum": { "sum": { "field": "cpu.usage" } },
						"cpu_throttling_avg": { "avg": { "field": "cpu.throttling" } },
						"cpu_pressure_avg": { "avg": { "field": "cpu.pressure" } },
						"memory_usage_avg": { "avg": { "field": "memory.usage" } },
						"memory_ws_sum": { "sum": { "field": "memory.working_set" } },
						"memory_rss_sum": { "sum": { "field": "memory.rss" } },
						"memory_pressure_avg": { "avg": { "field": "memory.pressure" } },
						"read_bps_sum": { "sum": { "field": "io.read_bps" } },
						"write_bps_sum": { "sum": { "field": "io.write_bps" } },
						"io_pressure_avg": { "avg": { "field": "io.pressure" } },
						"pids_sum": { "sum": { "field": "pids" } }
					}
				},

				"cpu_usage_overall": { "avg_bucket": { "buckets_path": "per_time>cpu_usage_sum" } },
				"cpu_throttling_overall": { "avg_bucket": { "buckets_path": "per_time>cpu_throttling_avg" } },
				"cpu_pressure_overall": { "avg_bucket": { "buckets_path": "per_time>cpu_pressure_avg" } },
				"memory_usage_overall": { "avg_bucket": { "buckets_path": "per_time>memory_usage_avg" } },
				"memory_ws_overall": { "avg_bucket": { "buckets_path": "per_time>memory_ws_sum" } },
				"memory_rss_overall": { "avg_bucket": { "buckets_path": "per_time>memory_rss_sum" } },
				"memory_pressure_overall": { "avg_bucket": { "buckets_path": "per_time>memory_pressure_avg" } },
				"read_bps_overall": { "avg_bucket": { "buckets_path": "per_time>read_bps_sum" } },
				"write_bps_overall": { "avg_bucket": { "buckets_path": "per_time>write_bps_sum" } },
				"io_pressure_overall": { "avg_bucket": { "buckets_path": "per_time>io_pressure_avg" } },
				"pids_overall": { "avg_bucket": { "buckets_path": "per_time>pids_sum" } },
				"oom_events_max": { "max": { "field": "oom.events" } },
				"oom_events_min": { "min": { "field": "oom.events" } },
				"oom_kills_max": { "max": { "field": "oom.kills" } },
				"oom_kills_min": { "min": { "field": "oom.kills" } }
			}
			}
		}
		}
		`,
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
		topN,
		reportBucketInterval,
	)

	res, err := r.es.Search(
		r.es.Search.WithContext(ctx),
		r.es.Search.WithIndex("server-metrics"),
		r.es.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var raw map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		return nil, err
	}

	result := make(map[uint]*ServerPushStats)

	aggs, ok := raw["aggregations"].(map[string]interface{})
	if !ok {
		return result, nil
	}

	servers, ok := aggs["servers"].(map[string]interface{})
	if !ok {
		return result, nil
	}

	buckets, ok := servers["buckets"].([]interface{})
	if !ok {
		return result, nil
	}

	// The agent sends every 30s now, so keep this aligned with the agent interval.
	// expectedBuckets := int(endTime.Sub(startTime) / reportBucketInterval)
	// if expectedBuckets < 1 {
	// 	expectedBuckets = 1
	// }

	for _, b := range buckets {
		bucket, ok := b.(map[string]interface{})
		if !ok {
			continue
		}

		serverID := uint(bucket["key"].(float64))

		getVal := func(name string) float64 {
			if v, ok := bucket[name].(map[string]interface{})["value"].(float64); ok {
				return v
			}
			return 0
		}

		// Count how many time buckets actually contain docs.
		// This is the push-based uptime definition.
		actualBuckets := 0
		if perTime, ok := bucket["per_time"].(map[string]interface{}); ok {
			if tb, ok := perTime["buckets"].([]interface{}); ok {
				for _, item := range tb {
					if m, ok := item.(map[string]interface{}); ok {
						if dc, ok := m["doc_count"].(float64); ok && int(dc) > 0 {
							actualBuckets++
						}
					}
				}
			}
		}

		// Clamp the start time for this server to its own creation time.
		effectiveStart := startTime.UTC()
		if createdAt, ok := serverCreatedAt[serverID]; ok {
			if createdAt.After(effectiveStart) {
				effectiveStart = createdAt.UTC()
			}
		}

		uptime := 0.0
		if effectiveStart.Before(endTime.UTC()) || effectiveStart.Equal(endTime.UTC()) {
			expectedBuckets := int(endTime.UTC().Sub(effectiveStart) / reportBucketInterval)
			if expectedBuckets < 1 {
				expectedBuckets = 1
			}

			uptime = float64(actualBuckets) / float64(expectedBuckets) * 100
			if uptime > 100 {
				uptime = 100
			}
		}

		stats := &ServerPushStats{
			CPUUsageAvg:      getVal("cpu_usage_overall"),
			CPUThrottlingAvg: getVal("cpu_throttling_overall"),
			CPUPressureAvg:   getVal("cpu_pressure_overall"),

			MemoryUsageAvg:      getVal("memory_usage_overall"),
			MemoryWorkingSetAvg: getVal("memory_ws_overall"),
			MemoryRSSAvg:        getVal("memory_rss_overall"),
			MemoryPressureAvg:   getVal("memory_pressure_overall"),

			ReadBPSAvg:    getVal("read_bps_overall"),
			WriteBPSAvg:   getVal("write_bps_overall"),
			IOPressureAvg: getVal("io_pressure_overall"),

			PIDsAvg: getVal("pids_overall"),

			OOMEventsTotal: getVal("oom_events_max") - getVal("oom_events_min"),
			OOMKillsTotal:  getVal("oom_kills_max") - getVal("oom_kills_min"),

			Uptime: uptime,
		}

		result[serverID] = stats
	}

	return result, nil
}
