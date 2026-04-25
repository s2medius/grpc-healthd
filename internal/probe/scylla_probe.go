package probe

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/your-org/grpc-healthd/internal/metrics"
)

// ScyllaProbe checks health of a ScyllaDB node by connecting to its native
// transport port and verifying the CQL OPTIONS handshake response.
type ScyllaProbe struct {
	address string
	timeout time.Duration
}

// NewScyllaProbe creates a new ScyllaProbe. If timeout is zero the default
// probe timeout is used.
func NewScyllaProbe(address string, timeout time.Duration) *ScyllaProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &ScyllaProbe{address: address, timeout: timeout}
}

// Check dials the ScyllaDB native transport port, sends a CQL OPTIONS request
// and expects a SUPPORTED response (opcode 0x06).
func (p *ScyllaProbe) Check(ctx context.Context) Result {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", p.address, p.timeout)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("scylla", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("dial error: %v", err)}
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(p.timeout))

	// CQL OPTIONS request: version=0x04, flags=0x00, stream=0x0000,
	// opcode=0x05 (OPTIONS), length=0x00000000
	options := []byte{0x04, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00}
	if _, err := conn.Write(options); err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("scylla", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("write error: %v", err)}
	}

	// Read the 9-byte CQL frame header.
	header := make([]byte, 9)
	if _, err := conn.Read(header); err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("scylla", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("read error: %v", err)}
	}

	// opcode is at byte index 4; 0x06 == SUPPORTED
	if header[4] != 0x06 {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("scylla", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("unexpected opcode: 0x%02x", header[4])}
	}

	duration := time.Since(start).Seconds()
	metrics.RecordProbe("scylla", p.address, true, duration)
	return Result{Healthy: true, Message: "scylla ok"}
}
