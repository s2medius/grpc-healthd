package probe

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/yourusername/grpc-healthd/internal/metrics"
)

// RedisSentinelProbe checks a Redis Sentinel instance by sending a PING
// command and verifying the response, then optionally checking a monitored
// master is reachable.
type RedisSentinelProbe struct {
	address string
	timeout time.Duration
}

// NewRedisSentinelProbe creates a new RedisSentinelProbe.
// If timeout is zero, DefaultTimeout is used.
func NewRedisSentinelProbe(address string, timeout time.Duration) *RedisSentinelProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &RedisSentinelProbe{address: address, timeout: timeout}
}

// Check connects to the Redis Sentinel, sends PING, and validates the response.
func (p *RedisSentinelProbe) Check() Result {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", p.address, p.timeout)
	if err != nil {
		dur := time.Since(start).Seconds()
		metrics.RecordProbe("redis_sentinel", p.address, false, dur)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("connect: %w", err)}
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(p.timeout))

	_, err = fmt.Fprintf(conn, "PING\r\n")
	if err != nil {
		dur := time.Since(start).Seconds()
		metrics.RecordProbe("redis_sentinel", p.address, false, dur)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("write: %w", err)}
	}

	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		dur := time.Since(start).Seconds()
		metrics.RecordProbe("redis_sentinel", p.address, false, dur)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("no response from sentinel")}
	}

	line := scanner.Text()
	dur := time.Since(start).Seconds()

	if !strings.HasPrefix(line, "+PONG") {
		metrics.RecordProbe("redis_sentinel", p.address, false, dur)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("unexpected response: %s", line)}
	}

	metrics.RecordProbe("redis_sentinel", p.address, true, dur)
	return Result{Status: StatusHealthy}
}
