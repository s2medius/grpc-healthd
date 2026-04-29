package probe

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type openSearchClusterHealth struct {
	Status string `json:"status"`
}

// OpenSearchProbe checks the health of an OpenSearch cluster via its
// /_cluster/health endpoint. Green and yellow statuses are considered healthy.
type OpenSearchProbe struct {
	address string
	timeout time.Duration
}

// NewOpenSearchProbe creates an OpenSearchProbe targeting the given address.
// If timeout is zero, DefaultTimeout is used.
func NewOpenSearchProbe(address string, timeout time.Duration) *OpenSearchProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &OpenSearchProbe{address: address, timeout: timeout}
}

func (p *OpenSearchProbe) Probe(ctx context.Context) Result {
	start := time.Now()

	url := fmt.Sprintf("http://%s/_cluster/health", p.address)
	client := &http.Client{Timeout: p.timeout}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Result{Status: StatusUnhealthy, Duration: time.Since(start), Error: err}
	}

	resp, err := client.Do(req)
	if err != nil {
		return Result{Status: StatusUnhealthy, Duration: time.Since(start), Error: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		return Result{Status: StatusUnhealthy, Duration: time.Since(start), Error: err}
	}

	var health openSearchClusterHealth
	if err = json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return Result{Status: StatusUnhealthy, Duration: time.Since(start), Error: err}
	}

	if health.Status == "red" {
		err = fmt.Errorf("opensearch cluster status is red")
		return Result{Status: StatusUnhealthy, Duration: time.Since(start), Error: err}
	}

	return Result{Status: StatusHealthy, Duration: time.Since(start)}
}
