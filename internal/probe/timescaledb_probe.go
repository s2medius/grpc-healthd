package probe

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/grpc-healthd/internal/metrics"
)

// TimescaleDBProbe checks connectivity to a TimescaleDB instance by performing
// a TCP handshake followed by a minimal PostgreSQL startup message exchange.
type TimescaleDBProbe struct {
	address string
	timeout time.Duration
}

// NewTimescaleDBProbe creates a new TimescaleDBProbe.
// If timeout is zero the default probe timeout is used.
func NewTimescaleDBProbe(address string, timeout time.Duration) *TimescaleDBProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &TimescaleDBProbe{address: address, timeout: timeout}
}

// Probe dials the TimescaleDB address, sends a PostgreSQL startup message,
// and expects the server to respond with an authentication request ('R'),
// error ('E'), or any valid single-byte response indicating a live server.
func (p *TimescaleDBProbe) Probe(ctx context.Context) Result {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", p.address, p.timeout)
	if err != nil {
		return Result{Status: StatusUnhealthy, Error: err, Duration: time.Since(start)}
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(p.timeout))

	// PostgreSQL startup message: length(4) + protocol(4) + user param
	startup := buildPostgresStartupMessage("healthcheck")
	if _, err := conn.Write(startup); err != nil {
		return Result{Status: StatusUnhealthy, Error: err, Duration: time.Since(start)}
	}

	buf := make([]byte, 1)
	if _, err := conn.Read(buf); err != nil {
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("no response from timescaledb: %w", err), Duration: time.Since(start)}
	}

	dur := time.Since(start)
	metrics.RecordProbe("timescaledb", p.address, true, dur)
	return Result{Status: StatusHealthy, Duration: dur}
}
