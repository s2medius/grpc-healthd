package probe

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/your-org/grpc-healthd/internal/metrics"
)

// ClickHouseProbe checks ClickHouse HTTP interface availability.
type ClickHouseProbe struct {
	address string
	timeout time.Duration
	client  *http.Client
}

// NewClickHouseProbe creates a new ClickHouseProbe targeting address (host:port).
// If timeout is zero, DefaultTimeout is used.
func NewClickHouseProbe(address string, timeout time.Duration) *ClickHouseProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &ClickHouseProbe{
		address: address,
		timeout: timeout,
		client:  &http.Client{Timeout: timeout},
	}
}

// Probe performs the health check against the ClickHouse HTTP ping endpoint.
func (p *ClickHouseProbe) Probe(ctx context.Context) Result {
	start := time.Now()

	host, port, err := net.SplitHostPort(p.address)
	if err != nil {
		dur := time.Since(start).Seconds()
		metrics.RecordProbe("clickhouse", p.address, false, dur)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("invalid address %q: %w", p.address, err)}
	}

	url := fmt.Sprintf("http://%s:%s/ping", host, port)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		dur := time.Since(start).Seconds()
		metrics.RecordProbe("clickhouse", p.address, false, dur)
		return Result{Status: StatusUnhealthy, Error: err}
	}

	resp, err := p.client.Do(req)
	dur := time.Since(start).Seconds()
	if err != nil {
		metrics.RecordProbe("clickhouse", p.address, false, dur)
		return Result{Status: StatusUnhealthy, Error: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		metrics.RecordProbe("clickhouse", p.address, false, dur)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("unexpected status %d", resp.StatusCode)}
	}

	metrics.RecordProbe("clickhouse", p.address, true, dur)
	return Result{Status: StatusHealthy}
}
