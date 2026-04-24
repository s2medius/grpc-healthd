package probe

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/patrickdappollonio/grpc-healthd/internal/metrics"
)

// CassandraProbe checks connectivity to a Cassandra node by establishing
// a TCP connection and validating the native transport protocol banner.
type CassandraProbe struct {
	address string
	timeout time.Duration
}

// NewCassandraProbe creates a new CassandraProbe targeting the given address.
// If timeout is zero, DefaultTimeout is used.
func NewCassandraProbe(address string, timeout time.Duration) *CassandraProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &CassandraProbe{address: address, timeout: timeout}
}

// Check dials the Cassandra native transport port and reads the initial
// OPTIONS/READY handshake frame to confirm the node is accepting connections.
func (p *CassandraProbe) Check(ctx context.Context) Result {
	start := time.Now()

	dialer := &net.Dialer{Timeout: p.timeout}
	conn, err := dialer.DialContext(ctx, "tcp", p.address)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("cassandra", p.address, false, duration)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("cassandra connect: %w", err)}
	}
	defer conn.Close()

	// Send a CQL OPTIONS request (version=3, flags=0, stream=0, opcode=5, length=0)
	options := []byte{0x03, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00}
	_ = conn.SetDeadline(time.Now().Add(p.timeout))
	if _, err := conn.Write(options); err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("cassandra", p.address, false, duration)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("cassandra write: %w", err)}
	}

	// Read at least 9 bytes (frame header) to confirm a response
	header := make([]byte, 9)
	if _, err := conn.Read(header); err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("cassandra", p.address, false, duration)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("cassandra read: %w", err)}
	}

	duration := time.Since(start).Seconds()
	metrics.RecordProbe("cassandra", p.address, true, duration)
	return Result{Status: StatusHealthy}
}
