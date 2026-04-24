package probe

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"grpc-healthd/internal/metrics"
)

// LDAPSProbe checks connectivity to an LDAPS (LDAP over TLS) server.
type LDAPSProbe struct {
	address    string
	timeout    time.Duration
	skipVerify bool
}

// NewLDAPSProbe creates a new LDAPSProbe.
// If timeout is zero, defaultTimeout is used.
func NewLDAPSProbe(address string, timeout time.Duration, skipVerify bool) *LDAPSProbe {
	if timeout == 0 {
		timeout = defaultTimeout
	}
	return &LDAPSProbe{address: address, timeout: timeout, skipVerify: skipVerify}
}

// Probe dials the LDAPS server over TLS and reads the first byte of the response.
func (p *LDAPSProbe) Probe(ctx context.Context) Result {
	start := time.Now()

	dialer := &tls.Dialer{
		Config: &tls.Config{
			InsecureSkipVerify: p.skipVerify, //nolint:gosec
		},
	}

	conn, err := dialer.DialContext(ctx, "tcp", p.address)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("ldaps", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("tls connection failed: %v", err)}
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(p.timeout))

	buf := make([]byte, 1)
	_, err = conn.Read(buf)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("ldaps", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("read failed: %v", err)}
	}

	duration := time.Since(start).Seconds()

	if buf[0] != 0x30 {
		metrics.RecordProbe("ldaps", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("unexpected banner byte: 0x%02x", buf[0])}
	}

	metrics.RecordProbe("ldaps", p.address, true, duration)
	return Result{Healthy: true, Message: "ldaps ok"}
}
