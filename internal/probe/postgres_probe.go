package probe

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/yourusername/grpc-healthd/internal/config"
	"github.com/yourusername/grpc-healthd/internal/metrics"
)

// PostgresProbe checks health by performing a TCP handshake and sending
// a minimal PostgreSQL startup message, verifying the server responds.
type PostgresProbe struct {
	address string
	timeout time.Duration
}

// NewPostgresProbe creates a PostgresProbe from a ProbeConfig.
func NewPostgresProbe(cfg config.ProbeConfig) *PostgresProbe {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return &PostgresProbe{
		address: cfg.Address,
		timeout: timeout,
	}
}

// Execute dials the PostgreSQL address and validates a response is received.
func (p *PostgresProbe) Execute(ctx context.Context) Result {
	start := time.Now()

	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", p.address)
	duration := time.Since(start)

	if err != nil {
		metrics.RecordProbe(p.address, "postgres", false, duration)
		return Result{
			Status:   StatusUnhealthy,
			Message:  fmt.Sprintf("connection failed: %v", err),
			Duration: duration,
		}
	}
	defer conn.Close()

	// Send a minimal PostgreSQL SSLRequest (8 bytes) to provoke a response.
	// The server should reply with 'N' (no SSL) or 'S' (SSL supported).
	sslRequest := []byte{0x00, 0x00, 0x00, 0x08, 0x04, 0xd2, 0x16, 0x2f}
	conn.SetDeadline(time.Now().Add(p.timeout))
	if _, err := conn.Write(sslRequest); err != nil {
		metrics.RecordProbe(p.address, "postgres", false, duration)
		return Result{
			Status:   StatusUnhealthy,
			Message:  fmt.Sprintf("write failed: %v", err),
			Duration: duration,
		}
	}

	buf := make([]byte, 1)
	if _, err := conn.Read(buf); err != nil {
		metrics.RecordProbe(p.address, "postgres", false, duration)
		return Result{
			Status:   StatusUnhealthy,
			Message:  fmt.Sprintf("read failed: %v", err),
			Duration: duration,
		}
	}

	metrics.RecordProbe(p.address, "postgres", true, duration)
	return Result{
		Status:   StatusHealthy,
		Message:  fmt.Sprintf("postgres responded: %q", string(buf)),
		Duration: duration,
	}
}
