package probe

import (
	"context"
	"fmt"
	"net"
	"time"

	"grpc-healthd/internal/metrics"
)

// ICMPProbe checks reachability by attempting a TCP dial to port 7 (echo)
// or falling back to a basic dial to confirm host resolution and TCP stack.
type ICMPProbe struct {
	address string
	timeout time.Duration
}

// NewICMPProbe creates a new ICMPProbe for the given host address.
// If timeout is zero, DefaultTimeout is used.
func NewICMPProbe(address string, timeout time.Duration) *ICMPProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &ICMPProbe{address: address, timeout: timeout}
}

// Probe attempts to dial the host on port 80 to verify network reachability.
// This is a TCP-based ping alternative that does not require raw sockets.
func (p *ICMPProbe) Probe(ctx context.Context) Result {
	start := time.Now()

	target := p.address
	// If no port specified, append :80 for reachability check
	if _, _, err := net.SplitHostPort(target); err != nil {
		target = fmt.Sprintf("%s:80", target)
	}

	dialCtx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	conn, err := (&net.Dialer{}).DialContext(dialCtx, "tcp", target)
	duration := time.Since(start)

	if err != nil {
		metrics.RecordProbe(p.address, "icmp", false, duration)
		return Result{Healthy: false, Message: err.Error(), Duration: duration}
	}
	conn.Close()
	metrics.RecordProbe(p.address, "icmp", true, duration)
	return Result{Healthy: true, Message: "reachable", Duration: duration}
}
