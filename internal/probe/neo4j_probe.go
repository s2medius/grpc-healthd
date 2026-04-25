package probe

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/patrickward/grpc-healthd/internal/metrics"
)

// Neo4jProbe checks Neo4j availability via the HTTP discovery endpoint.
type Neo4jProbe struct {
	address string
	timeout time.Duration
	client  *http.Client
}

// NewNeo4jProbe creates a new Neo4jProbe for the given address.
// address should be in the form "host:port" (e.g. "localhost:7474").
func NewNeo4jProbe(address string, timeout time.Duration) *Neo4jProbe {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &Neo4jProbe{
		address: address,
		timeout: timeout,
		client:  &http.Client{Timeout: timeout},
	}
}

// Check performs an HTTP GET against the Neo4j discovery endpoint and
// returns a healthy Result when the server responds with 200 OK.
func (p *Neo4jProbe) Check(ctx context.Context) Result {
	start := time.Now()

	url := fmt.Sprintf("http://%s/", p.address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		dur := time.Since(start).Seconds()
		metrics.RecordProbe("neo4j", p.address, false, dur)
		return Result{Healthy: false, Error: err}
	}

	resp, err := p.client.Do(req)
	dur := time.Since(start).Seconds()
	if err != nil {
		metrics.RecordProbe("neo4j", p.address, false, dur)
		return Result{Healthy: false, Error: err}
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		metrics.RecordProbe("neo4j", p.address, false, dur)
		return Result{
			Healthy: false,
			Error:   fmt.Errorf("neo4j: unexpected status %d", resp.StatusCode),
		}
	}

	metrics.RecordProbe("neo4j", p.address, true, dur)
	return Result{Healthy: true}
}
