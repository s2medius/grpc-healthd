package probe

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/yourorg/grpc-healthd/internal/metrics"
)

// mongoDBHandshake is the minimal MongoDB isMaster wire protocol message.
var mongoDBHandshake = []byte{
	0x3a, 0x00, 0x00, 0x00, // messageLength
	0x01, 0x00, 0x00, 0x00, // requestID
	0x00, 0x00, 0x00, 0x00, // responseTo
	0xd4, 0x07, 0x00, 0x00, // opCode OP_QUERY
	0x00, 0x00, 0x00, 0x00, // flags
	0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x24, 0x63, 0x6d, 0x64, 0x00, // fullCollectionName: admin.$cmd
	0x00, 0x00, 0x00, 0x00, // numberToSkip
	0x01, 0x00, 0x00, 0x00, // numberToReturn
	// BSON doc: {isMaster: 1}
	0x13, 0x00, 0x00, 0x00,
	0x10, 0x69, 0x73, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x00,
	0x01, 0x00, 0x00, 0x00,
	0x00,
}

// MongoDBProbe checks connectivity to a MongoDB instance.
type MongoDBProbe struct {
	address string
	timeout time.Duration
}

// NewMongoDBProbe creates a new MongoDBProbe.
func NewMongoDBProbe(address string, timeout time.Duration) *MongoDBProbe {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	return &MongoDBProbe{address: address, timeout: timeout}
}

// Probe connects to MongoDB, sends a minimal isMaster handshake, and reads the response header.
func (p *MongoDBProbe) Probe(ctx context.Context) Result {
	start := time.Now()

	dialer := &net.Dialer{Timeout: p.timeout}
	conn, err := dialer.DialContext(ctx, "tcp", p.address)
	if err != nil {
		d := time.Since(start)
		metrics.RecordProbe("mongodb", p.address, false, d)
		return Result{Status: StatusUnhealthy, Duration: d, Err: fmt.Errorf("mongodb connect: %w", err)}
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(p.timeout))

	if _, err := conn.Write(mongoDBHandshake); err != nil {
		d := time.Since(start)
		metrics.RecordProbe("mongodb", p.address, false, d)
		return Result{Status: StatusUnhealthy, Duration: d, Err: fmt.Errorf("mongodb write: %w", err)}
	}

	// Read at least the 4-byte message length from the response header.
	header := make([]byte, 4)
	if _, err := conn.Read(header); err != nil {
		d := time.Since(start)
		metrics.RecordProbe("mongodb", p.address, false, d)
		return Result{Status: StatusUnhealthy, Duration: d, Err: fmt.Errorf("mongodb read: %w", err)}
	}

	d := time.Since(start)
	metrics.RecordProbe("mongodb", p.address, true, d)
	return Result{Status: StatusHealthy, Duration: d}
}
