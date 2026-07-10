package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/HenryNg101/server-service/internal/model"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/go-openapi/testify/v2/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupElasticsearch(t *testing.T) *elasticsearch.Client {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "docker.elastic.co/elasticsearch/elasticsearch:8.13.4",
		ExposedPorts: []string{"9200/tcp"},
		Env: map[string]string{
			"discovery.type":         "single-node",
			"xpack.security.enabled": "false",
			"ES_JAVA_OPTS":           "-Xms512m -Xmx512m",
		},
		WaitingFor: wait.ForHTTP("/").
			WithPort("9200/tcp").
			WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		container.Terminate(ctx)
	})

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "9200")

	addr := fmt.Sprintf("http://%s:%s", host, port.Port())

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{addr},
	})
	require.NoError(t, err)

	// wait until ready
	for i := 0; i < 10; i++ {
		_, err := es.Info()
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	require.NoError(t, err)

	return es
}

func createIndex(t *testing.T, es *elasticsearch.Client) {
	mapping := `{
	  "mappings": {
	    "properties": {
	      "@timestamp": { "type": "date" },
	      "server_id": { "type": "long" },
	      "status": { "type": "boolean" }
	    }
	  }
	}`

	res, err := es.Indices.Create(
		"server-status",
		es.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	require.NoError(t, err)
	defer res.Body.Close()
}

func TestElasticRepository_BulkInsertStatus(t *testing.T) {
	es := setupElasticsearch(t)
	createIndex(t, es)

	repo := NewServerESRepository(es)

	now := time.Now()

	servers := []*model.Server{
		{ID: 1, Status: true, LastUpdated: now},
		{ID: 2, Status: false, LastUpdated: now},
	}

	err := repo.BulkInsertStatus(context.Background(), servers)
	require.NoError(t, err)

	_, err = es.Indices.Refresh(
		es.Indices.Refresh.WithIndex("server-status"),
	)
	require.NoError(t, err)

	// verify inserted
	res, err := es.Count(es.Count.WithIndex("server-status"))
	require.NoError(t, err)
	defer res.Body.Close()

	var body map[string]interface{}
	json.NewDecoder(res.Body).Decode(&body)

	require.Equal(t, float64(2), body["count"])
}

func TestElasticRepository_GetDailyUptime(t *testing.T) {
	es := setupElasticsearch(t)
	createIndex(t, es)

	repo := NewServerESRepository(es)

	now := time.Now()

	// Insert test data
	data := []*model.Server{
		{ID: 1, Status: true, LastUpdated: now.Add(-1 * time.Hour)},
		{ID: 1, Status: false, LastUpdated: now.Add(-30 * time.Minute)},
		{ID: 2, Status: true, LastUpdated: now.Add(-1 * time.Hour)},
	}

	err := repo.BulkInsertStatus(context.Background(), data)
	require.NoError(t, err)

	// ES needs refresh to make docs searchable
	_, err = es.Indices.Refresh(es.Indices.Refresh.WithIndex("server-status"))
	require.NoError(t, err)

	result, err := repo.GetDailyUptime(
		context.Background(),
		now.Add(-2*time.Hour),
		now,
		10,
	)
	require.NoError(t, err)

	require.Contains(t, result, uint(1))
	require.Contains(t, result, uint(2))

	// server 1: (true + false) / 2 = 0.5
	require.InDelta(t, 0.5, result[1], 0.01)

	// server 2: only true = 1.0
	require.InDelta(t, 1.0, result[2], 0.01)
}

func TestElasticRepository_GetDailyUptime_ESFailure(t *testing.T) {
	es, _ := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9999"}, // invalid
	})

	repo := NewServerESRepository(es)

	_, err := repo.GetDailyUptime(context.Background(), time.Now(), time.Now(), 10)

	require.Error(t, err)
}

func TestElasticRepository_BulkInsertStatus_Error(t *testing.T) {
	es, _ := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"http://localhost:9999"},
	})

	repo := NewServerESRepository(es)

	err := repo.BulkInsertStatus(context.Background(), []*model.Server{
		{ID: 1, Status: true, LastUpdated: time.Now()},
	})

	require.Error(t, err)
}
