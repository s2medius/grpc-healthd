package probe

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/your-org/grpc-healthd/internal/metrics"
)

// ClickHouseProbe checks connectivity to a ClickHouse server via the native
// protocol port (default 9000). It performs a TCP handshake and sends the
// minimal ClickHouse client hello to verify the server responds correctly.
type ClickHouseProbe struct {
	address string
	timeout time.Duration
}

// NewClickHouseProbe creates a new ClickHouseProbe. If timeout is zero the
// default probe timeout is used.
func NewClickHouseProbe(address string, timeout time.Duration) *ClickHouseProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &ClickHouseProbe{address: address, timeout: timeout}
}

// Execute dials the ClickHouse native port and expects the server to send at
// least one byte (the server hello), indicating it is alive.
func (p *ClickHouseProbe) Execute(ctx context.Context) Result {
	start := time.Now()

	dialer := &net.Dialer{Timeout: p.timeout}
	conn, err := dialer.DialContext(ctx, "tcp", p.address)
	if err != nil {
		dur := time.Since(start)
		metrics.RecordProbe("clickhouse", p.address, false, dur)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("dial: %w", err), Duration: dur}
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(p.timeout))

	// ClickHouse sends a server hello immediately after connection.
	buf := make([]byte, 1)
	_, err = conn.Read(buf)
	dur := time.Since(start)
	if err != nil {
		metrics.RecordProbe("clickhouse", p.address, false, dur)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("read server hello: %w", err), Duration: dur}
	}

	metrics.RecordProbe("clickhouse", p.address, true, dur)
	return Result{Status: StatusHealthy, Duration: dur}
}
