package probe

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/yourusername/grpc-healthd/internal/metrics"
)

// MySQLProbe checks MySQL availability by performing a TCP handshake
// and validating the server greeting packet.
type MySQLProbe struct {
	address string
	timeout time.Duration
}

// NewMySQLProbe creates a new MySQLProbe for the given address.
// If timeout is zero, DefaultTimeout is used.
func NewMySQLProbe(address string, timeout time.Duration) *MySQLProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &MySQLProbe{address: address, timeout: timeout}
}

// Check connects to the MySQL server and validates the server greeting.
func (p *MySQLProbe) Check(ctx context.Context) Result {
	start := time.Now()

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", p.address)
	if err != nil {
		dur := time.Since(start)
		metrics.RecordProbe("mysql", p.address, false, dur)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("mysql connect: %w", err)}
	}
	defer conn.Close()

	// MySQL sends a greeting packet; read at least 5 bytes to confirm.
	_ = conn.SetReadDeadline(time.Now().Add(p.timeout))
	buf := make([]byte, 5)
	n, err := conn.Read(buf)
	dur := time.Since(start)

	if err != nil || n < 5 {
		metrics.RecordProbe("mysql", p.address, false, dur)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("mysql greeting: %w", err)}
	}

	metrics.RecordProbe("mysql", p.address, true, dur)
	return Result{Status: StatusHealthy}
}
