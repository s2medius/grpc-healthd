package probe

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/your-org/grpc-healthd/internal/metrics"
)

// SplunkProbe checks Splunk health via its REST API /services/server/health/splunkd/details.
type SplunkProbe struct {
	address string
	timeout time.Duration
	client  *http.Client
}

// NewSplunkProbe creates a new SplunkProbe targeting the given address (e.g. "https://splunk:8089").
// If timeout is zero, DefaultTimeout is used.
func NewSplunkProbe(address string, timeout time.Duration) *SplunkProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &SplunkProbe{
		address: address,
		timeout: timeout,
		client:  &http.Client{Timeout: timeout},
	}
}

// Execute performs the Splunk health check.
func (p *SplunkProbe) Execute(ctx context.Context) Result {
	start := time.Now()

	url := fmt.Sprintf("%s/services/server/health/splunkd?output_mode=json", p.address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		dur := time.Since(start)
		metrics.RecordProbe("splunk", StatusUnhealthy.String(), dur)
		return Result{Status: StatusUnhealthy, Duration: dur, Err: err}
	}

	resp, err := p.client.Do(req)
	dur := time.Since(start)
	if err != nil {
		metrics.RecordProbe("splunk", StatusUnhealthy.String(), dur)
		return Result{Status: StatusUnhealthy, Duration: dur, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		metrics.RecordProbe("splunk", StatusUnhealthy.String(), dur)
		return Result{Status: StatusUnhealthy, Duration: dur, Err: err}
	}

	metrics.RecordProbe("splunk", StatusHealthy.String(), dur)
	return Result{Status: StatusHealthy, Duration: dur}
}
