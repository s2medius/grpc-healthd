package probe

import (
	"context"
	"fmt"
	"net"
	"time"

	"grpc-healthd/internal/metrics"
)

// LDAPProbe checks connectivity to an LDAP server by performing a TCP handshake
// and verifying the server responds with a valid LDAP banner byte (0x30).
type LDAPProbe struct {
	address string
	timeout time.Duration
}

// NewLDAPProbe creates a new LDAPProbe for the given address.
// If timeout is zero, defaultTimeout is used.
func NewLDAPProbe(address string, timeout time.Duration) *LDAPProbe {
	if timeout == 0 {
		timeout = defaultTimeout
	}
	return &LDAPProbe{address: address, timeout: timeout}
}

// Probe connects to the LDAP server and checks for a valid response byte.
func (p *LDAPProbe) Probe(ctx context.Context) Result {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", p.address, p.timeout)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("ldap", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("connection failed: %v", err)}
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(p.timeout))

	buf := make([]byte, 1)
	_, err = conn.Read(buf)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("ldap", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("read failed: %v", err)}
	}

	duration := time.Since(start).Seconds()

	// LDAP responses start with ASN.1 SEQUENCE tag 0x30
	if buf[0] != 0x30 {
		metrics.RecordProbe("ldap", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("unexpected banner byte: 0x%02x", buf[0])}
	}

	metrics.RecordProbe("ldap", p.address, true, duration)
	return Result{Healthy: true, Message: "ldap ok"}
}
