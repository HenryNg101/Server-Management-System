package elastic

import (
	"fmt"
	"log"
	"net/http"

	_ "embed"

	"github.com/HenryNg101/metrics-consumer/internal/config"
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
