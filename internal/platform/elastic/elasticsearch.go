package elastic

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"

	_ "embed"

	"github.com/HenryNg101/server-management-system/internal/config"
	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v9"
)

func NewElasticsearchSession(config *config.ElasticSearchConfig) *elasticsearch.Client {
	es, err := elasticsearch.New(
		elasticsearch.WithAddresses(fmt.Sprintf("%s:%s", config.Host, config.Port)),
		elasticsearch.WithTransportOptions(
			// Intentionally disable any attempt of using HTTP/2.0 for simplicity
			// Got network hanging issue before, and I checked that default config use this option as "true"
			elastictransport.WithTransport(&http.Transport{
				ForceAttemptHTTP2: false,
			}),
		),
		elasticsearch.WithBasicAuth(config.User, config.Password),
	)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// Test connection
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	return es
}

func InitElasticsearch(es *elasticsearch.Client) {
	ctx := context.Background()

	createIndexTemplateIfNotExist(es, ctx)
	createDataStreamIfNotExist(es, ctx)
}

//go:embed server-status-template-v1.json
var serverStatusTemplate []byte

func createIndexTemplateIfNotExist(es *elasticsearch.Client, ctx context.Context) {
	// Check if the template exists already
	templateName := "server-status-template"
	res, err := es.Indices.ExistsIndexTemplate(templateName)
	if err != nil {
		log.Fatalf("check template error: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		log.Println("Index template for data stream is already exists")
		return
	}

	//
	// Create index template
	res, err = es.Indices.PutIndexTemplate(
		templateName,
		bytes.NewReader(serverStatusTemplate),
	)
	if err != nil || res.IsError() {
		log.Fatalf("create template failed: %v", err)
	}

	log.Println("Index template for data stream created")
}

func createDataStreamIfNotExist(es *elasticsearch.Client, ctx context.Context) {
	name := "server-status"

	// Check existence
	res, err := es.Indices.GetDataStream(
		es.Indices.GetDataStream.WithName(name),
	)
	if err == nil && res.StatusCode == 200 {
		log.Println("Data stream already exists")
		res.Body.Close()
		return
	}

	// Create
	res, err = es.Indices.CreateDataStream(name)
	if err != nil || res.IsError() {
		log.Fatalf("Create data stream failed: %v", err)
	}

	log.Println("Data stream created")
}
