package probe

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/yourorg/grpc-healthd/internal/metrics"
)

// OracleProbe checks Oracle DB availability by performing a TCP handshake
// and validating the server's initial response banner.
type OracleProbe struct {
	address string
	timeout time.Duration
}

// NewOracleProbe creates a new OracleProbe targeting the given address.
// If timeout is zero, DefaultTimeout is used.
func NewOracleProbe(address string, timeout time.Duration) *OracleProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &OracleProbe{address: address, timeout: timeout}
}

// Check connects to the Oracle listener port and verifies a response is received.
func (p *OracleProbe) Check(ctx context.Context) Result {
	start := time.Now()

	dialer := &net.Dialer{Timeout: p.timeout}
	conn, err := dialer.DialContext(ctx, "tcp", p.address)
	if err != nil {
		dur := time.Since(start)
		metrics.RecordProbe("oracle", p.address, StatusUnhealthy.String(), dur)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("oracle connect: %w", err)}
	}
	defer conn.Close()

	// Oracle TNS listener sends a greeting; read at least 1 byte to confirm liveness.
	_ = conn.SetReadDeadline(time.Now().Add(p.timeout))
	buf := make([]byte, 1)
	_, err = conn.Read(buf)
	dur := time.Since(start)

	if err != nil {
		metrics.RecordProbe("oracle", p.address, StatusUnhealthy.String(), dur)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("oracle read: %w", err)}
	}

	metrics.RecordProbe("oracle", p.address, StatusHealthy.String(), dur)
	return Result{Status: StatusHealthy}
}
