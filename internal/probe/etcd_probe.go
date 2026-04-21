package probe

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// EtcdProbe checks etcd health via its HTTP /health endpoint.
type EtcdProbe struct {
	address string
	timeout time.Duration
	client  *http.Client
}

// NewEtcdProbe creates a new EtcdProbe. If timeout is zero, a 5-second default is used.
func NewEtcdProbe(address string, timeout time.Duration) *EtcdProbe {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return &EtcdProbe{
		address: address,
		timeout: timeout,
		client:  &http.Client{Timeout: timeout},
	}
}

// Probe performs the etcd health check.
func (e *EtcdProbe) Probe() Result {
	start := time.Now()

	url := fmt.Sprintf("http://%s/health", e.address)
	resp, err := e.client.Get(url)
	duration := time.Since(start)

	if err != nil {
		return Result{
			Healthy:  false,
			Message:  fmt.Sprintf("etcd request failed: %v", err),
			Duration: duration,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{
			Healthy:  false,
			Message:  fmt.Sprintf("etcd returned status %d", resp.StatusCode),
			Duration: duration,
		}
	}

	var body struct {
		Health string `json:"health"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return Result{
			Healthy:  false,
			Message:  fmt.Sprintf("failed to decode etcd response: %v", err),
			Duration: duration,
		}
	}

	if body.Health != "true" {
		return Result{
			Healthy:  false,
			Message:  fmt.Sprintf("etcd reports unhealthy: %s", body.Health),
			Duration: duration,
		}
	}

	return Result{
		Healthy:  true,
		Message:  "etcd is healthy",
		Duration: duration,
	}
}
