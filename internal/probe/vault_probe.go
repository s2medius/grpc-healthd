package probe

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/grpc-healthd/internal/metrics"
)

// VaultProbe checks HashiCorp Vault health via its /v1/sys/health endpoint.
type VaultProbe struct {
	address string
	timeout time.Duration
	client  *http.Client
}

// NewVaultProbe creates a new VaultProbe targeting the given address.
func NewVaultProbe(address string, timeout time.Duration) *VaultProbe {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &VaultProbe{
		address: address,
		timeout: timeout,
		client:  &http.Client{Timeout: timeout},
	}
}

// Probe performs the health check against Vault's sys/health endpoint.
// Vault returns 200 (initialized, unsealed, active) or 429 (standby) as healthy.
func (p *VaultProbe) Probe(ctx context.Context) Result {
	start := time.Now()

	url := fmt.Sprintf("%s/v1/sys/health?standbyok=true", p.address)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe(p.address, "vault", false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("failed to build request: %v", err)}
	}

	resp, err := p.client.Do(req)
	duration := time.Since(start).Seconds()
	if err != nil {
		metrics.RecordProbe(p.address, "vault", false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("request failed: %v", err)}
	}
	defer resp.Body.Close()

	// 200 = active, 429 = standby (healthy with standbyok=true), 473 = performance standby
	healthy := resp.StatusCode == http.StatusOK ||
		resp.StatusCode == http.StatusTooManyRequests ||
		resp.StatusCode == 473

	var body struct {
		Initialized bool   `json:"initialized"`
		Sealed      bool   `json:"sealed"`
		Version     string `json:"version"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&body)

	metrics.RecordProbe(p.address, "vault", healthy, duration)
	if !healthy {
		return Result{Healthy: false, Message: fmt.Sprintf("unexpected status %d", resp.StatusCode)}
	}
	return Result{Healthy: true, Message: fmt.Sprintf("vault %s ok (sealed=%v)", body.Version, body.Sealed)}
}
