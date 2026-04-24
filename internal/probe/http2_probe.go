package probe

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"

	"github.com/yourorg/grpc-healthd/internal/metrics"
)

// HTTP2Probe checks health by performing an HTTP/2 GET request.
type HTTP2Probe struct {
	address string
	timeout time.Duration
}

// NewHTTP2Probe creates a new HTTP2Probe. If timeout is zero, DefaultTimeout is used.
func NewHTTP2Probe(address string, timeout time.Duration) *HTTP2Probe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &HTTP2Probe{address: address, timeout: timeout}
}

// Probe performs an HTTP/2 GET against the configured address and returns a Result.
func (p *HTTP2Probe) Probe(ctx context.Context) Result {
	start := time.Now()

	transport := &http2.Transport{
		AllowHTTP: true,
		DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
			return net.DialTimeout(network, addr, p.timeout)
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   p.timeout,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.address, nil)
	if err != nil {
		dur := time.Since(start)
		metrics.RecordProbe("http2", p.address, false, dur)
		return Result{Status: StatusUnhealthy, Duration: dur, Err: fmt.Errorf("build request: %w", err)}
	}

	resp, err := client.Do(req)
	dur := time.Since(start)
	if err != nil {
		metrics.RecordProbe("http2", p.address, false, dur)
		return Result{Status: StatusUnhealthy, Duration: dur, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		metrics.RecordProbe("http2", p.address, true, dur)
		return Result{Status: StatusHealthy, Duration: dur}
	}

	metrics.RecordProbe("http2", p.address, false, dur)
	return Result{
		Status:   StatusUnhealthy,
		Duration: dur,
		Err:      fmt.Errorf("unexpected status code: %d", resp.StatusCode),
	}
}
