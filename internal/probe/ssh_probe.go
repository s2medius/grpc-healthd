package probe

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/yourorg/grpc-healthd/internal/metrics"
)

// SSHProbe checks health by establishing a TCP connection to an SSH server
// and verifying that the server responds with a valid SSH protocol banner
// (i.e. a line beginning with "SSH-").
type SSHProbe struct {
	address string
	timeout time.Duration
}

// NewSSHProbe creates a new SSHProbe targeting the given address (host:port).
// If timeout is zero, DefaultTimeout is used.
func NewSSHProbe(address string, timeout time.Duration) *SSHProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &SSHProbe{
		address: address,
		timeout: timeout,
	}
}

// Probe connects to the SSH server and validates the banner.
// It returns a Result indicating whether the server is healthy.
func (p *SSHProbe) Probe(ctx context.Context) Result {
	start := time.Now()

	// Respect context deadline alongside the probe timeout.
	deadline := start.Add(p.timeout)
	if d, ok := ctx.Deadline(); ok && d.Before(deadline) {
		deadline = d
	}

	conn, err := net.DialTimeout("tcp", p.address, time.Until(deadline))
	if err != nil {
		duration := time.Since(start)
		metrics.RecordProbe("ssh", p.address, false, duration)
		return Result{
			Status:   StatusUnhealthy,
			Duration: duration,
			Message:  fmt.Sprintf("connection failed: %v", err),
		}
	}
	defer conn.Close()

	// Set a read deadline for receiving the banner.
	if err := conn.SetReadDeadline(deadline); err != nil {
		duration := time.Since(start)
		metrics.RecordProbe("ssh", p.address, false, duration)
		return Result{
			Status:   StatusUnhealthy,
			Duration: duration,
			Message:  fmt.Sprintf("failed to set read deadline: %v", err),
		}
	}

	// Read the SSH banner (up to 255 bytes as per RFC 4253).
	buf := make([]byte, 255)
	n, err := conn.Read(buf)
	if err != nil {
		duration := time.Since(start)
		metrics.RecordProbe("ssh", p.address, false, duration)
		return Result{
			Status:   StatusUnhealthy,
			Duration: duration,
			Message:  fmt.Sprintf("failed to read banner: %v", err),
		}
	}

	banner := string(buf[:n])
	duration := time.Since(start)

	// RFC 4253 §4.2: the server identification string MUST begin with "SSH-".
	if len(banner) < 4 || banner[:4] != "SSH-" {
		metrics.RecordProbe("ssh", p.address, false, duration)
		return Result{
			Status:   StatusUnhealthy,
			Duration: duration,
			Message:  fmt.Sprintf("unexpected banner: %q", banner),
		}
	}

	metrics.RecordProbe("ssh", p.address, true, duration)
	return Result{
		Status:   StatusHealthy,
		Duration: duration,
		Message:  fmt.Sprintf("SSH banner received: %q", banner),
	}
}
