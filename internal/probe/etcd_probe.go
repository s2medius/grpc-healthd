package probe

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/your-org/grpc-healthd/internal/config"
	"github.com/your-org/grpc-healthd/internal/metrics"
)

// EtcdProbe checks the health of an etcd cluster by querying its
// /health HTTP endpoint on the client port (default 2379).
type EtcdProbe struct {
	address string
	timeout time.Duration
	client  *http.Client
}

// NewEtcdProbe creates an EtcdProbe from the given ProbeConfig.
// The address should be in the form "host:port" where port is
// typically 2379 (etcd client port). If Timeout is zero the
// default probe timeout is used.
func NewEtcdProbe(cfg config.ProbeConfig) *EtcdProbe {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &EtcdProbe{
		address: cfg.Address,
		timeout: timeout,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Probe performs an HTTP GET against the etcd /health endpoint and
// returns a healthy Status when the response body contains
// {"health":"true"}. Any non-200 response or connection error is
// treated as unhealthy.
func (p *EtcdProbe) Probe(ctx context.Context) Status {
	start := time.Now()

	url := fmt.Sprintf("http://%s/health", p.address)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("etcd", p.address, false, duration)
		return Status{
			Healthy: false,
			Message: fmt.Sprintf("etcd: failed to build request: %v", err),
		}
	}

	resp, err := p.client.Do(req)
	duration := time.Since(start).Seconds()

	if err != nil {
		metrics.RecordProbe("etcd", p.address, false, duration)
		return Status{
			Healthy: false,
			Message: fmt.Sprintf("etcd: connection error: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		metrics.RecordProbe("etcd", p.address, false, duration)
		return Status{
			Healthy: false,
			Message: fmt.Sprintf("etcd: unexpected status %d from /health", resp.StatusCode),
		}
	}

	metrics.RecordProbe("etcd", p.address, true, duration)
	return Status{
		Healthy: true,
		Message: fmt.Sprintf("etcd: healthy (%s)", p.address),
	}
}
