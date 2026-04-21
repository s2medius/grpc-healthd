package probe

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"grpc-healthd/internal/metrics"
)

// NATSProbe checks connectivity to a NATS server by establishing a TCP
// connection and verifying the server sends the expected INFO banner.
type NATSProbe struct {
	address string
	timeout time.Duration
}

// NewNATSProbe creates a NATSProbe for the given address.
// If timeout is zero, DefaultTimeout is used.
func NewNATSProbe(address string, timeout time.Duration) *NATSProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &NATSProbe{address: address, timeout: timeout}
}

// Probe connects to the NATS server and checks for the INFO banner.
func (p *NATSProbe) Probe(ctx context.Context) Result {
	start := time.Now()

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(p.timeout)
	}

	conn, err := net.DialTimeout("tcp", p.address, time.Until(deadline))
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("nats", p.address, false, duration)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("nats connect: %w", err)}
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(deadline)

	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("nats", p.address, false, duration)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("nats read banner: %w", err)}
	}

	duration := time.Since(start).Seconds()
	banner := string(buf[:n])
	if !strings.HasPrefix(banner, "INFO ") {
		metrics.RecordProbe("nats", p.address, false, duration)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("nats unexpected banner: %q", banner)}
	}

	metrics.RecordProbe("nats", p.address, true, duration)
	return Result{Status: StatusHealthy}
}
