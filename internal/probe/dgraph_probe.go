package probe

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/your-org/grpc-healthd/internal/metrics"
)

// DgraphProbe checks health of a Dgraph instance via its /health HTTP endpoint.
type DgraphProbe struct {
	address string
	timeout time.Duration
	client  *http.Client
}

// NewDgraphProbe creates a new DgraphProbe. address should be host:port.
func NewDgraphProbe(address string, timeout time.Duration) *DgraphProbe {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &DgraphProbe{
		address: address,
		timeout: timeout,
		client:  &http.Client{Timeout: timeout},
	}
}

func (p *DgraphProbe) Probe(ctx context.Context) Result {
	start := time.Now()
	url := fmt.Sprintf("http://%s/health", p.address)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		dur := time.Since(start)
		metrics.RecordProbe("dgraph", p.address, false, dur)
		return Result{Healthy: false, Duration: dur, Err: err}
	}

	resp, err := p.client.Do(req)
	dur := time.Since(start)
	if err != nil {
		metrics.RecordProbe("dgraph", p.address, false, dur)
		return Result{Healthy: false, Duration: dur, Err: err}
	}
	defer resp.Body.Close()

	healthy := resp.StatusCode == http.StatusOK
	metrics.RecordProbe("dgraph", p.address, healthy, dur)
	if !healthy {
		return Result{
			Healthy:  false,
			Duration: dur,
			Err:      fmt.Errorf("dgraph /health returned status %d", resp.StatusCode),
		}
	}
	return Result{Healthy: true, Duration: dur}
}
