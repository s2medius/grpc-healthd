package probe

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/yourorg/grpc-healthd/internal/metrics"
)

// ConsulProbe checks health by querying the Consul agent health endpoint.
type ConsulProbe struct {
	address string
	timeout time.Duration
	client  *http.Client
}

// NewConsulProbe creates a new ConsulProbe targeting the given address.
// If timeout is zero, DefaultTimeout is used.
func NewConsulProbe(address string, timeout time.Duration) *ConsulProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &ConsulProbe{
		address: address,
		timeout: timeout,
		client:  &http.Client{Timeout: timeout},
	}
}

// Check performs the Consul health probe.
// It queries /v1/agent/self and expects a 200 response.
func (p *ConsulProbe) Check(ctx context.Context) Result {
	start := time.Now()
	url := fmt.Sprintf("http://%s/v1/agent/self", p.address)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		dur := time.Since(start)
		metrics.RecordProbe("consul", p.address, false, dur)
		return Result{Status: StatusUnhealthy, Duration: dur, Err: err}
	}

	resp, err := p.client.Do(req)
	dur := time.Since(start)
	if err != nil {
		metrics.RecordProbe("consul", p.address, false, dur)
		return Result{Status: StatusUnhealthy, Duration: dur, Err: err}
	}
	defer resp.Body.Close()

	// Decode just enough to confirm the agent responded meaningfully.
	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		metrics.RecordProbe("consul", p.address, false, dur)
		return Result{Status: StatusUnhealthy, Duration: dur, Err: fmt.Errorf("invalid response: %w", err)}
	}

	if resp.StatusCode != http.StatusOK {
		metrics.RecordProbe("consul", p.address, false, dur)
		return Result{Status: StatusUnhealthy, Duration: dur, Err: fmt.Errorf("unexpected status: %d", resp.StatusCode)}
	}

	metrics.RecordProbe("consul", p.address, true, dur)
	return Result{Status: StatusHealthy, Duration: dur}
}
