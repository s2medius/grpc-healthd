package probe

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/your-org/grpc-healthd/internal/metrics"
)

// ElasticsearchProbe checks the health of an Elasticsearch cluster by querying
// the /_cluster/health endpoint and verifying the cluster status is not "red".
type ElasticsearchProbe struct {
	address string
	timeout time.Duration
	client  *http.Client
}

// elasticsearchClusterHealth represents the relevant fields from the
// Elasticsearch /_cluster/health API response.
type elasticsearchClusterHealth struct {
	Status string `json:"status"`
}

// NewElasticsearchProbe creates a new ElasticsearchProbe targeting the given
// address (e.g. "http://localhost:9200"). If timeout is zero, DefaultTimeout
// is used.
func NewElasticsearchProbe(address string, timeout time.Duration) *ElasticsearchProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &ElasticsearchProbe{
		address: address,
		timeout: timeout,
		client:  &http.Client{Timeout: timeout},
	}
}

// Probe performs the Elasticsearch cluster health check. It returns a healthy
// Status when the cluster status is "green" or "yellow", and unhealthy when
// the status is "red" or the endpoint is unreachable.
func (e *ElasticsearchProbe) Probe(ctx context.Context) Status {
	start := time.Now()

	url := fmt.Sprintf("%s/_cluster/health", e.address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("elasticsearch", e.address, false, duration)
		return Status{
			Healthy: false,
			Message: fmt.Sprintf("failed to build request: %v", err),
		}
	}

	resp, err := e.client.Do(req)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("elasticsearch", e.address, false, duration)
		return Status{
			Healthy: false,
			Message: fmt.Sprintf("connection failed: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("elasticsearch", e.address, false, duration)
		return Status{
			Healthy: false,
			Message: fmt.Sprintf("unexpected HTTP status: %d", resp.StatusCode),
		}
	}

	var health elasticsearchClusterHealth
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("elasticsearch", e.address, false, duration)
		return Status{
			Healthy: false,
			Message: fmt.Sprintf("failed to decode response: %v", err),
		}
	}

	duration := time.Since(start).Seconds()

	if health.Status == "red" {
		metrics.RecordProbe("elasticsearch", e.address, false, duration)
		return Status{
			Healthy: false,
			Message: fmt.Sprintf("cluster status is red"),
		}
	}

	metrics.RecordProbe("elasticsearch", e.address, true, duration)
	return Status{
		Healthy: true,
		Message: fmt.Sprintf("cluster status: %s", health.Status),
	}
}
