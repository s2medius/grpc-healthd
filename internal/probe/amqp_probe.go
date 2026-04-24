package probe

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/yourorg/grpc-healthd/internal/config"
	"github.com/yourorg/grpc-healthd/internal/metrics"
)

// AMQPProbe checks health of an AMQP broker (e.g. RabbitMQ) by verifying
// the server sends a valid AMQP 0-9-1 protocol header upon connection.
type AMQPProbe struct {
	address string
	timeout time.Duration
}

// NewAMQPProbe creates a new AMQPProbe from the provided config.
func NewAMQPProbe(cfg config.ProbeConfig) *AMQPProbe {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return &AMQPProbe{
		address: cfg.Address,
		timeout: timeout,
	}
}

// Probe connects to the AMQP broker and validates the protocol header.
func (p *AMQPProbe) Probe(ctx context.Context) Result {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", p.address, p.timeout)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("amqp", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("connection failed: %v", err)}
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(p.timeout))

	// AMQP 0-9-1 servers send "AMQP" as the first 4 bytes of their greeting.
	buf := make([]byte, 4)
	if _, err := conn.Read(buf); err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("amqp", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("read failed: %v", err)}
	}

	duration := time.Since(start).Seconds()
	if string(buf) != "AMQP" {
		metrics.RecordProbe("amqp", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("unexpected banner: %q", string(buf))}
	}

	metrics.RecordProbe("amqp", p.address, true, duration)
	return Result{Healthy: true, Message: "AMQP broker reachable"}
}
