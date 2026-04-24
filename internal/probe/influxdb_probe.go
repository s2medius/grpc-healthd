package probe

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yourusername/grpc-healthd/internal/metrics"
)

// InfluxDBProbe checks health by hitting the InfluxDB /ping endpoint.
type InfluxDBProbe struct {
	address string
	timeout time.Duration
	client  *http.Client
}

// NewInfluxDBProbe creates a new InfluxDBProbe for the given address.
// address should be in the form "host:port" (e.g. "localhost:8086").
func NewInfluxDBProbe(address string, timeout time.Duration) *InfluxDBProbe {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &InfluxDBProbe{
		address: address,
		timeout: timeout,
		client:  &http.Client{Timeout: timeout},
	}
}

// Probe performs the health check against the InfluxDB /ping endpoint.
func (p *InfluxDBProbe) Probe(ctx context.Context) Result {
	start := time.Now()

	url := fmt.Sprintf("http://%s/ping", p.address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("influxdb", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("failed to create request: %v", err)}
	}

	resp, err := p.client.Do(req)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("influxdb", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("request failed: %v", err)}
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, resp.Body)

	duration := time.Since(start).Seconds()
	// InfluxDB /ping returns 204 No Content when healthy.
	if resp.StatusCode == http.StatusNoContent {
		metrics.RecordProbe("influxdb", p.address, true, duration)
		return Result{Healthy: true, Message: "influxdb ping ok"}
	}

	metrics.RecordProbe("influxdb", p.address, false, duration)
	return Result{Healthy: false, Message: fmt.Sprintf("unexpected status code: %d", resp.StatusCode)}
}
