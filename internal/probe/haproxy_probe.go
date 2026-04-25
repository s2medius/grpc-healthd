package probe

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/your-org/grpc-healthd/internal/metrics"
)

// HAProxyProbe checks HAProxy health via the stats socket or HTTP health endpoint.
// It connects to the address and sends a simple "show info" command over TCP,
// expecting a valid HAProxy response banner.
type HAProxyProbe struct {
	address string
	timeout time.Duration
}

// NewHAProxyProbe creates a new HAProxyProbe. If timeout is zero, defaultTimeout is used.
func NewHAProxyProbe(address string, timeout time.Duration) *HAProxyProbe {
	if timeout == 0 {
		timeout = defaultTimeout
	}
	return &HAProxyProbe{address: address, timeout: timeout}
}

// Probe connects to the HAProxy stats endpoint and validates the response.
func (p *HAProxyProbe) Probe() Result {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", p.address, p.timeout)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("haproxy", p.address, false, duration)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("haproxy: connection failed: %w", err)}
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(p.timeout))

	// Send show info command
	_, err = fmt.Fprintf(conn, "show info\n")
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("haproxy", p.address, false, duration)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("haproxy: write failed: %w", err)}
	}

	scanner := bufio.NewScanner(conn)
	if scanner.Scan() {
		line := scanner.Text()
		duration := time.Since(start).Seconds()
		if strings.HasPrefix(line, "Name:") || strings.Contains(line, "HAProxy") {
			metrics.RecordProbe("haproxy", p.address, true, duration)
			return Result{Status: StatusHealthy}
		}
		metrics.RecordProbe("haproxy", p.address, false, duration)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("haproxy: unexpected response: %q", line)}
	}

	duration := time.Since(start).Seconds()
	metrics.RecordProbe("haproxy", p.address, false, duration)
	return Result{Status: StatusUnhealthy, Error: fmt.Errorf("haproxy: no response received")}
}
