package probe

import (
	"context"
	"fmt"
	"net"
	"time"
)

// Status represents the result of a probe check.
type Status int

const (
	StatusUnknown Status = iota
	StatusHealthy
	StatusUnhealthy
)

func (s Status) String() string {
	switch s {
	case StatusHealthy:
		return "healthy"
	case StatusUnhealthy:
		return "unhealthy"
	default:
		return "unknown"
	}
}

// Result holds the outcome of a single probe execution.
type Result struct {
	Name    string
	Status  Status
	Latency time.Duration
	Err     error
}

// Probe defines the interface for health check probes.
type Probe interface {
	Name() string
	Check(ctx context.Context) Result
}

// TCPProbe checks connectivity to a TCP endpoint.
type TCPProbe struct {
	name    string
	address string
	timeout time.Duration
}

// NewTCPProbe creates a new TCP probe.
func NewTCPProbe(name, address string, timeout time.Duration) *TCPProbe {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return &TCPProbe{name: name, address: address, timeout: timeout}
}

func (p *TCPProbe) Name() string { return p.name }

func (p *TCPProbe) Check(ctx context.Context) Result {
	start := time.Now()
	result := Result{Name: p.name}

	dialer := &net.Dialer{Timeout: p.timeout}
	conn, err := dialer.DialContext(ctx, "tcp", p.address)
	result.Latency = time.Since(start)

	if err != nil {
		result.Status = StatusUnhealthy
		result.Err = fmt.Errorf("tcp dial %s: %w", p.address, err)
		return result
	}
	conn.Close()
	result.Status = StatusHealthy
	return result
}
