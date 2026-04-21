package probe

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/yourorg/grpc-healthd/internal/metrics"
)

// MemcachedProbe checks health by sending a "version\r\n" command
// and verifying the response starts with "VERSION".
type MemcachedProbe struct {
	address string
	timeout time.Duration
}

// NewMemcachedProbe creates a new MemcachedProbe.
// If timeout is zero, DefaultTimeout is used.
func NewMemcachedProbe(address string, timeout time.Duration) *MemcachedProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &MemcachedProbe{address: address, timeout: timeout}
}

// Check dials the Memcached server, sends the VERSION command,
// and verifies the response to determine health.
func (p *MemcachedProbe) Check() Result {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", p.address, p.timeout)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("memcached", p.address, false, duration)
		return Result{Status: StatusUnhealthy, Error: err}
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(p.timeout))

	if _, err = fmt.Fprintf(conn, "version\r\n"); err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("memcached", p.address, false, duration)
		return Result{Status: StatusUnhealthy, Error: err}
	}

	buf := make([]byte, 64)
	n, err := conn.Read(buf)
	if err != nil && err != io.EOF {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("memcached", p.address, false, duration)
		return Result{Status: StatusUnhealthy, Error: err}
	}

	duration := time.Since(start).Seconds()
	response := string(buf[:n])

	if !strings.HasPrefix(response, "VERSION") {
		metrics.RecordProbe("memcached", p.address, false, duration)
		return Result{
			Status: StatusUnhealthy,
			Error:  fmt.Errorf("unexpected memcached response: %q", response),
		}
	}

	metrics.RecordProbe("memcached", p.address, true, duration)
	return Result{Status: StatusHealthy}
}
