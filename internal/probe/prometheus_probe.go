package probe

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// PrometheusProbe checks a Prometheus-compatible /metrics or custom endpoint
// and verifies that a specific metric name is present in the response.
type PrometheusProbe struct {
	address    string
	metricName string
	timeout    time.Duration
	client     *http.Client
}

// NewPrometheusProbe creates a new PrometheusProbe.
// address should be a full URL (e.g. http://localhost:9090/metrics).
// metricName is the metric label to look for in the response body.
func NewPrometheusProbe(address, metricName string, timeout time.Duration) *PrometheusProbe {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &PrometheusProbe{
		address:    address,
		metricName: metricName,
		timeout:    timeout,
		client:     &http.Client{Timeout: timeout},
	}
}

// Probe performs the health check by fetching the metrics endpoint and
// verifying the expected metric name appears in the output.
func (p *PrometheusProbe) Probe(ctx context.Context) Result {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.address, nil)
	if err != nil {
		return Result{Healthy: false, Duration: time.Since(start), Err: err}
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return Result{Healthy: false, Duration: time.Since(start), Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{
			Healthy:  false,
			Duration: time.Since(start),
			Err:      fmt.Errorf("unexpected status code: %d", resp.StatusCode),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Healthy: false, Duration: time.Since(start), Err: err}
	}

	if p.metricName != "" && !strings.Contains(string(body), p.metricName) {
		return Result{
			Healthy:  false,
			Duration: time.Since(start),
			Err:      fmt.Errorf("metric %q not found in response", p.metricName),
		}
	}

	return Result{Healthy: true, Duration: time.Since(start)}
}
