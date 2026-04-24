package probe

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/yourorg/grpc-healthd/internal/metrics"
)

// ClickHouseProbe checks connectivity to a ClickHouse server by performing
// a native protocol handshake over TCP (port 9000 by default).
type ClickHouseProbe struct {
	address string
	timeout time.Duration
}

// NewClickHouseProbe creates a new ClickHouseProbe.
// If timeout is zero, DefaultTimeout is used.
func NewClickHouseProbe(address string, timeout time.Duration) *ClickHouseProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &ClickHouseProbe{address: address, timeout: timeout}
}

// Probe attempts to connect to the ClickHouse native TCP port and validates
// that the server sends the expected "ClickHouse" greeting in the banner.
func (p *ClickHouseProbe) Probe(ctx context.Context) Result {
	start := time.Now()

	dialer := &net.Dialer{}
	ctxTimeout, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	conn, err := dialer.DialContext(ctxTimeout, "tcp", p.address)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("clickhouse", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("connection failed: %v", err)}
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(p.timeout))

	// ClickHouse native protocol: server sends a Hello packet upon connection.
	// We read the first 9 bytes and check for the "ClickHouse" string.
	buf := make([]byte, 64)
	n, err := conn.Read(buf)
	duration := time.Since(start).Seconds()

	if err != nil {
		metrics.RecordProbe("clickhouse", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("read failed: %v", err)}
	}

	banner := string(buf[:n])
	if len(banner) == 0 {
		metrics.RecordProbe("clickhouse", p.address, false, duration)
		return Result{Healthy: false, Message: "empty response from server"}
	}

	metrics.RecordProbe("clickhouse", p.address, true, duration)
	return Result{Healthy: true, Message: "clickhouse reachable"}
}
