package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
)

type ElasticAgentRepository interface {
	BulkInsertStatus(ctx context.Context, results []MetricMessage) error
	GetServersStats(ctx context.Context, startTime time.Time, endTime time.Time, topN int) (map[uint]*ServerPushStats, error)
	GetServerStats(ctx context.Context, serverID int, startTime time.Time, endTime time.Time) (*ServerPushStats, error)
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
		meta := `{"create":{"_index":"server-metrics"}}` + "\n"

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

func (r *elasticAgentRepository) GetServersStats(ctx context.Context, startTime time.Time, endTime time.Time, topN int) (map[uint]*ServerPushStats, error) {
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
					"fixed_interval": "10s",
					"min_doc_count": 1
				},
				"aggs": {
					"cpu_usage_sum": {
					"sum": { "field": "cpu.usage" }
					},
					"cpu_throttling_avg": {
					"avg": { "field": "cpu.throttling" }
					},
					"cpu_pressure_avg": {
					"avg": { "field": "cpu.pressure" }
					},

					"memory_usage_avg": {
					"avg": { "field": "memory.usage" }
					},
					"memory_ws_sum": {
					"sum": { "field": "memory.working_set" }
					},
					"memory_rss_sum": {
					"sum": { "field": "memory.rss" }
					},
					"memory_pressure_avg": {
					"avg": { "field": "memory.pressure" }
					},

					"read_bps_sum": {
					"sum": { "field": "io.read_bps" }
					},
					"write_bps_sum": {
					"sum": { "field": "io.write_bps" }
					},
					"io_pressure_avg": {
					"avg": { "field": "io.pressure" }
					},

					"pids_sum": {
					"sum": { "field": "pids" }
					}
				}
				},

				"cpu_usage_overall": {
				"avg_bucket": {
					"buckets_path": "per_time>cpu_usage_sum"
				}
				},
				"cpu_throttling_overall": {
				"avg_bucket": {
					"buckets_path": "per_time>cpu_throttling_avg"
				}
				},
				"cpu_pressure_overall": {
				"avg_bucket": {
					"buckets_path": "per_time>cpu_pressure_avg"
				}
				},

				"memory_usage_overall": {
				"avg_bucket": {
					"buckets_path": "per_time>memory_usage_avg"
				}
				},
				"memory_ws_overall": {
				"avg_bucket": {
					"buckets_path": "per_time>memory_ws_sum"
				}
				},
				"memory_rss_overall": {
				"avg_bucket": {
					"buckets_path": "per_time>memory_rss_sum"
				}
				},
				"memory_pressure_overall": {
				"avg_bucket": {
					"buckets_path": "per_time>memory_pressure_avg"
				}
				},

				"read_bps_overall": {
				"avg_bucket": {
					"buckets_path": "per_time>read_bps_sum"
				}
				},
				"write_bps_overall": {
				"avg_bucket": {
					"buckets_path": "per_time>write_bps_sum"
				}
				},
				"io_pressure_overall": {
				"avg_bucket": {
					"buckets_path": "per_time>io_pressure_avg"
				}
				},

				"pids_overall": {
				"avg_bucket": {
					"buckets_path": "per_time>pids_sum"
				}
				},

				"oom_events_max": {
				"max": { "field": "oom.events" }
				},
				"oom_events_min": {
				"min": { "field": "oom.events" }
				},
				"oom_kills_max": {
				"max": { "field": "oom.kills" }
				},
				"oom_kills_min": {
				"min": { "field": "oom.kills" }
				}
			}
			}
		}
		}
		`,
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
		topN,
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

	aggs := raw["aggregations"].(map[string]interface{})
	servers := aggs["servers"].(map[string]interface{})
	buckets := servers["buckets"].([]interface{})

	for _, b := range buckets {
		bucket := b.(map[string]interface{})

		serverID := uint(bucket["key"].(float64))

		getVal := func(name string) float64 {
			if v, ok := bucket[name].(map[string]interface{})["value"].(float64); ok {
				return v
			}
			return 0
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
		}

		result[serverID] = stats
	}

	return result, nil
}

func (r *elasticAgentRepository) GetServerStats(ctx context.Context, serverID int, startTime time.Time, endTime time.Time) (*ServerPushStats, error) {
	query := fmt.Sprintf(`
		{
			"size": 0,
			"query": {
				"bool": {
					"filter": [
						{
							"term": {
								"server_id": %d
							}
						},
						{
							"range": {
								"@timestamp": {
									"gte": "%s",
									"lte": "%s"
								}
							}
						}
					]
				}
			},
			"aggs": {
				"per_time": {
					"date_histogram": {
						"field": "@timestamp",
						"fixed_interval": "10s",
						"min_doc_count": 1
					},
					"aggs": {
						"cpu_usage_sum": {
							"sum": {
								"field": "cpu.usage"
							}
						},
						"cpu_throttling_avg": {
							"avg": {
								"field": "cpu.throttling"
							}
						},
						"cpu_pressure_avg": {
							"avg": {
								"field": "cpu.pressure"
							}
						},
						"memory_usage_avg": {
							"avg": {
								"field": "memory.usage"
							}
						},
						"memory_ws_sum": {
							"sum": {
								"field": "memory.working_set"
							}
						},
						"memory_rss_sum": {
							"sum": {
								"field": "memory.rss"
							}
						},
						"memory_pressure_avg": {
							"avg": {
								"field": "memory.pressure"
							}
						},
						"read_bps_sum": {
							"sum": {
								"field": "io.read_bps"
							}
						},
						"write_bps_sum": {
							"sum": {
								"field": "io.write_bps"
							}
						},
						"io_pressure_avg": {
							"avg": {
								"field": "io.pressure"
							}
						},
						"pids_sum": {
							"sum": {
								"field": "pids"
							}
						}
					}
				},
				"cpu_usage_overall": {
					"avg_bucket": {
						"buckets_path": "per_time>cpu_usage_sum"
					}
				},
				"cpu_throttling_overall": {
					"avg_bucket": {
						"buckets_path": "per_time>cpu_throttling_avg"
					}
				},
				"cpu_pressure_overall": {
					"avg_bucket": {
						"buckets_path": "per_time>cpu_pressure_avg"
					}
				},
				"memory_usage_overall": {
					"avg_bucket": {
						"buckets_path": "per_time>memory_usage_avg"
					}
				},
				"memory_ws_overall": {
					"avg_bucket": {
						"buckets_path": "per_time>memory_ws_sum"
					}
				},
				"memory_rss_overall": {
					"avg_bucket": {
						"buckets_path": "per_time>memory_rss_sum"
					}
				},
				"memory_pressure_overall": {
					"avg_bucket": {
						"buckets_path": "per_time>memory_pressure_avg"
					}
				},
				"read_bps_overall": {
					"avg_bucket": {
						"buckets_path": "per_time>read_bps_sum"
					}
				},
				"write_bps_overall": {
					"avg_bucket": {
						"buckets_path": "per_time>write_bps_sum"
					}
				},
				"io_pressure_overall": {
					"avg_bucket": {
						"buckets_path": "per_time>io_pressure_avg"
					}
				},
				"pids_overall": {
					"avg_bucket": {
						"buckets_path": "per_time>pids_sum"
					}
				},
				"oom_events_max": {
					"max": {
						"field": "oom.events"
					}
				},
				"oom_events_min": {
					"min": {
						"field": "oom.events"
					}
				},
				"oom_kills_max": {
					"max": {
						"field": "oom.kills"
					}
				},
				"oom_kills_min": {
					"min": {
						"field": "oom.kills"
					}
				}
			}
		}`,
		serverID,
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
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

	aggs := raw["aggregations"].(map[string]interface{})

	getVal := func(name string) float64 {
		if v, ok := aggs[name].(map[string]interface{})["value"].(float64); ok {
			return v
		}
		return 0
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
	}

	return stats, nil
}
