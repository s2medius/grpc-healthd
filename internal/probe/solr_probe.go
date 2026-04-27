package probe

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourusername/grpc-healthd/internal/metrics"
)

// SolrProbe checks health of an Apache Solr instance via its admin ping endpoint.
type SolrProbe struct {
	address string
	timeout time.Duration
	client  *http.Client
}

// NewSolrProbe creates a new SolrProbe targeting the given address (host:port).
// Timeout defaults to 5 seconds if not specified.
func NewSolrProbe(address string, timeout time.Duration) *SolrProbe {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return &SolrProbe{
		address: address,
		timeout: timeout,
		client:  &http.Client{Timeout: timeout},
	}
}

// Check performs an HTTP GET against the Solr admin ping endpoint and
// returns Healthy if the status field equals "OK".
func (p *SolrProbe) Check(ctx context.Context) Result {
	start := time.Now()
	url := fmt.Sprintf("http://%s/solr/admin/ping?wt=json", p.address)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("solr", p.address, false, duration)
		return Result{Status: Unhealthy, Error: err}
	}

	resp, err := p.client.Do(req)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("solr", p.address, false, duration)
		return Result{Status: Unhealthy, Error: err}
	}
	defer resp.Body.Close()

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("solr", p.address, false, duration)
		return Result{Status: Unhealthy, Error: fmt.Errorf("decode response: %w", err)}
	}

	duration := time.Since(start).Seconds()
	if body.Status != "OK" {
		metrics.RecordProbe("solr", p.address, false, duration)
		return Result{Status: Unhealthy, Error: fmt.Errorf("unexpected status: %q", body.Status)}
	}

	metrics.RecordProbe("solr", p.address, true, duration)
	return Result{Status: Healthy}
}
