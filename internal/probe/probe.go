package probe

import (
	"context"
	"fmt"
	"net"
	"time"
)

// Status represents the health state of a service.
type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
	StatusUnknown   Status = "unknown"
)

func (s Status) String() string { return string(s) }

// Result holds the outcome of a single probe execution.
type Result struct {
	Status   Status
	Duration time.Duration
	Err      error
}

// Probe is the interface implemented by all probe types.
type Probe interface {
	Check(ctx context.Context) Result
}

// TCPProbe checks health by opening a TCP connection.
type TCPProbe struct {
	Address string
	Timeout time.Duration
}

const defaultTCPTimeout = 5 * time.Second

// NewTCPProbe creates a TCPProbe, applying a default timeout if zero.
func NewTCPProbe(address string, timeout time.Duration) *TCPProbe {
	if timeout == 0 {
		timeout = defaultTCPTimeout
	}
	return &TCPProbe{Address: address, Timeout: timeout}
}

// Check attempts a TCP dial and returns the result.
func (p *TCPProbe) Check(ctx context.Context) Result {
	start := time.Now()
	dialCtx, cancel := context.WithTimeout(ctx, p.Timeout)
	defer cancel()

	conn, err := (&net.Dialer{}).DialContext(dialCtx, "tcp", p.Address)
	duration := time.Since(start)
	if err != nil {
		return Result{
			Status:   StatusUnhealthy,
			Duration: duration,
			Err:      fmt.Errorf("tcp probe failed: %w", err),
		}
	}
	conn.Close()
	return Result{Status: StatusHealthy, Duration: duration}
}
