package probe

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/grpc-healthd/internal/metrics"
)

// CouchDBProbe checks CouchDB health via its /_up endpoint.
type CouchDBProbe struct {
	address string
	timeout time.Duration
	client  *http.Client
}

// NewCouchDBProbe creates a new CouchDBProbe.
func NewCouchDBProbe(address string, timeout time.Duration) *CouchDBProbe {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return &CouchDBProbe{
		address: address,
		timeout: timeout,
		client:  &http.Client{Timeout: timeout},
	}
}

// Execute runs the CouchDB health check.
func (p *CouchDBProbe) Execute(ctx context.Context) Result {
	start := time.Now()

	url := fmt.Sprintf("http://%s/_up", p.address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("couchdb", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("failed to create request: %v", err)}
	}

	resp, err := p.client.Do(req)
	duration := time.Since(start).Seconds()
	if err != nil {
		metrics.RecordProbe("couchdb", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("connection failed: %v", err)}
	}
	defer resp.Body.Close()

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		metrics.RecordProbe("couchdb", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("failed to decode response: %v", err)}
	}

	if resp.StatusCode != http.StatusOK || body.Status != "ok" {
		metrics.RecordProbe("couchdb", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("unhealthy status: %s (HTTP %d)", body.Status, resp.StatusCode)}
	}

	metrics.RecordProbe("couchdb", p.address, true, duration)
	return Result{Healthy: true, Message: "couchdb is up"}
}
