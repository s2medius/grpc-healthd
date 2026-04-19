package probe

import (
	"fmt"
	"net/http"
	"time"
)

// HTTPProbe checks health by performing an HTTP GET request.
type HTTPProbe struct {
	URL     string
	Timeout time.Duration
}

// NewHTTPProbe creates a new HTTPProbe with the given URL and optional timeout.
// If timeout is zero, DefaultTimeout is used.
func NewHTTPProbe(url string, timeout time.Duration) *HTTPProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &HTTPProbe{URL: url, Timeout: timeout}
}

// Check performs an HTTP GET to the configured URL and returns a Result.
func (p *HTTPProbe) Check() Result {
	client := &http.Client{Timeout: p.Timeout}

	start := time.Now()
	resp, err := client.Get(p.URL)
	duration := time.Since(start)

	if err != nil {
		return Result{
			Status:   StatusUnhealthy,
			Duration: duration,
			Message:  fmt.Sprintf("http get failed: %v", err),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return Result{
			Status:   StatusHealthy,
			Duration: duration,
			Message:  fmt.Sprintf("http %d", resp.StatusCode),
		}
	}

	return Result{
		Status:   StatusUnhealthy,
		Duration: duration,
		Message:  fmt.Sprintf("unexpected http status: %d", resp.StatusCode),
	}
}
