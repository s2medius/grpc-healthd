package probe

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/yourusername/grpc-healthd/internal/metrics"
)

// ValkeyProbe checks a Valkey (Redis-compatible) instance by sending PING
// and verifying the +PONG response. Valkey is a Redis fork and uses the
// same RESP protocol.
type ValkeyProbe struct {
	address string
	timeout time.Duration
}

// NewValkeyProbe creates a new ValkeyProbe.
// If timeout is zero, DefaultTimeout is used.
func NewValkeyProbe(address string, timeout time.Duration) *ValkeyProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &ValkeyProbe{address: address, timeout: timeout}
}

// Check dials the Valkey server, sends PING, and validates the response.
func (p *ValkeyProbe) Check() Result {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", p.address, p.timeout)
	if err != nil {
		dur := time.Since(start).Seconds()
		metrics.RecordProbe("valkey", p.address, false, dur)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("connect: %w", err)}
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(p.timeout))

	_, err = fmt.Fprintf(conn, "*1\r\n$4\r\nPING\r\n")
	if err != nil {
		dur := time.Since(start).Seconds()
		metrics.RecordProbe("valkey", p.address, false, dur)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("write: %w", err)}
	}

	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		dur := time.Since(start).Seconds()
		metrics.RecordProbe("valkey", p.address, false, dur)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("no response from valkey")}
	}

	line := scanner.Text()
	dur := time.Since(start).Seconds()

	if !strings.HasPrefix(line, "+PONG") {
		metrics.RecordProbe("valkey", p.address, false, dur)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("unexpected response: %s", line)}
	}

	metrics.RecordProbe("valkey", p.address, true, dur)
	return Result{Status: StatusHealthy}
}
