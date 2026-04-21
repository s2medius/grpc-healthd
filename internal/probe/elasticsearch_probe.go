package probe

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ElasticsearchProbe checks Elasticsearch cluster health via the REST API.
type ElasticsearchProbe struct {
	address string
	timeout time.Duration
	client  *http.Client
}

type esClusterHealth struct {
	Status string `json:"status"`
}

// NewElasticsearchProbe creates a new ElasticsearchProbe.
// If timeout is zero, a default of 5 seconds is used.
func NewElasticsearchProbe(address string, timeout time.Duration) Probe {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &ElasticsearchProbe{
		address: address,
		timeout: timeout,
		client:  &http.Client{Timeout: timeout},
	}
}

// Check performs an HTTP GET to the Elasticsearch /_cluster/health endpoint
// and returns healthy if the cluster status is "green" or "yellow".
func (e *ElasticsearchProbe) Check() Result {
	start := time.Now()
	url := fmt.Sprintf("http://%s/_cluster/health", e.address)

	resp, err := e.client.Get(url)
	if err != nil {
		return Result{
			Healthy:  false,
			Message:  fmt.Sprintf("elasticsearch connection failed: %v", err),
			Duration: time.Since(start),
		}
	}
	defer resp.Body.Close()

	var health esClusterHealth
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return Result{
			Healthy:  false,
			Message:  fmt.Sprintf("failed to decode cluster health response: %v", err),
			Duration: time.Since(start),
		}
	}

	duration := time.Since(start)

	switch health.Status {
	case "green", "yellow":
		return Result{
			Healthy:  true,
			Message:  fmt.Sprintf("elasticsearch cluster status: %s", health.Status),
			Duration: duration,
		}
	default:
		return Result{
			Healthy:  false,
			Message:  fmt.Sprintf("elasticsearch cluster status unhealthy: %s", health.Status),
			Duration: duration,
		}
	}
}
