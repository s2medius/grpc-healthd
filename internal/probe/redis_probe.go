package probe

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/yourusername/grpc-healthd/internal/config"
	"github.com/yourusername/grpc-healthd/internal/metrics"
)

// RedisProbe checks health by sending a PING command to a Redis server.
type RedisProbe struct {
	address string
	timeout time.Duration
}

// NewRedisProbe creates a new RedisProbe from the given config.
func NewRedisProbe(cfg config.ProbeConfig) *RedisProbe {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return &RedisProbe{
		address: cfg.Address,
		timeout: timeout,
	}
}

// Run executes the Redis PING probe and returns a Status.
func (p *RedisProbe) Run(ctx context.Context) Status {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", p.address, p.timeout)
	if err != nil {
		metrics.RecordProbe("redis", p.address, false, time.Since(start).Seconds())
		return Status{Healthy: false, Message: fmt.Sprintf("dial failed: %v", err)}
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(p.timeout))

	// Send Redis inline PING command
	if _, err := fmt.Fprintf(conn, "PING\r\n"); err != nil {
		metrics.RecordProbe("redis", p.address, false, time.Since(start).Seconds())
		return Status{Healthy: false, Message: fmt.Sprintf("write failed: %v", err)}
	}

	buf := make([]byte, 7) // "+PONG\r\n"
	n, err := conn.Read(buf)
	if err != nil {
		metrics.RecordProbe("redis", p.address, false, time.Since(start).Seconds())
		return Status{Healthy: false, Message: fmt.Sprintf("read failed: %v", err)}
	}

	response := string(buf[:n])
	if len(response) < 5 || response[:5] != "+PONG" {
		metrics.RecordProbe("redis", p.address, false, time.Since(start).Seconds())
		return Status{Healthy: false, Message: fmt.Sprintf("unexpected response: %q", response)}
	}

	metrics.RecordProbe("redis", p.address, true, time.Since(start).Seconds())
	return Status{Healthy: true, Message: "PONG received"}
}
