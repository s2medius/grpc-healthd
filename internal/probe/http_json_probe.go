package probe

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/your-org/grpc-healthd/internal/metrics"
)

// HTTPJSONProbe checks an HTTP endpoint and validates a JSON field in the response body.
type HTTPJSONProbe struct {
	address   string
	jsonKey   string
	jsonValue string
	timeout   time.Duration
}

// NewHTTPJSONProbe creates an HTTPJSONProbe. It expects the response body to be a JSON
// object containing jsonKey equal to jsonValue.
func NewHTTPJSONProbe(address, jsonKey, jsonValue string, timeout time.Duration) *HTTPJSONProbe {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &HTTPJSONProbe{
		address:   address,
		jsonKey:   jsonKey,
		jsonValue: jsonValue,
		timeout:   timeout,
	}
}

// Probe performs the HTTP request and validates the JSON response.
func (p *HTTPJSONProbe) Probe(ctx context.Context) Result {
	start := time.Now()

	client := &http.Client{Timeout: p.timeout}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.address, nil)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe(p.address, "http_json", false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("build request: %v", err)}
	}

	resp, err := client.Do(req)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe(p.address, "http_json", false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("request failed: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe(p.address, "http_json", false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("unexpected status: %d", resp.StatusCode)}
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe(p.address, "http_json", false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("read body: %v", err)}
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe(p.address, "http_json", false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("invalid JSON: %v", err)}
	}

	val, ok := data[p.jsonKey]
	duration := time.Since(start).Seconds()
	if !ok {
		metrics.RecordProbe(p.address, "http_json", false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("key %q not found in response", p.jsonKey)}
	}

	if fmt.Sprintf("%v", val) != p.jsonValue {
		metrics.RecordProbe(p.address, "http_json", false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("key %q: got %v, want %s", p.jsonKey, val, p.jsonValue)}
	}

	metrics.RecordProbe(p.address, "http_json", true, duration)
	return Result{Healthy: true, Message: "ok"}
}
