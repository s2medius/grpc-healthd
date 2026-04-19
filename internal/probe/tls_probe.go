package probe

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

// TLSProbe checks that a TLS handshake succeeds and the certificate is valid.
type TLSProbe struct {
	address string
	timeout time.Duration
	skipVerify bool
}

// NewTLSProbe creates a new TLSProbe for the given address (host:port).
func NewTLSProbe(address string, timeout time.Duration, skipVerify bool) *TLSProbe {
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return &TLSProbe{address: address, timeout: timeout, skipVerify: skipVerify}
}

// Probe performs the TLS check and returns a Result.
func (p *TLSProbe) Probe(ctx context.Context) Result {
	start := time.Now()

	dialer := &tls.Dialer{
		NetDialer: &net.Dialer{Timeout: p.timeout},
		Config:    &tls.Config{InsecureSkipVerify: p.skipVerify}, //nolint:gosec
	}

	conn, err := dialer.DialContext(ctx, "tcp", p.address)
	duration := time.Since(start)

	if err != nil {
		return Result{Status: StatusUnhealthy, Duration: duration, Error: fmt.Errorf("tls dial: %w", err)}
	}
	conn.Close()
	return Result{Status: StatusHealthy, Duration: duration}
}
