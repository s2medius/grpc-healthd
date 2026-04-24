package probe

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/your-org/grpc-healthd/internal/metrics"
)

// SMTPProbe checks health of an SMTP server by connecting and
// verifying the 220 greeting banner.
type SMTPProbe struct {
	address string
	timeout time.Duration
}

// NewSMTPProbe creates a new SMTPProbe targeting the given address.
// If timeout is zero, DefaultTimeout is used.
func NewSMTPProbe(address string, timeout time.Duration) *SMTPProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &SMTPProbe{address: address, timeout: timeout}
}

// Check connects to the SMTP server and validates the 220 greeting.
func (p *SMTPProbe) Check() Result {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", p.address, p.timeout)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("smtp", p.address, StatusUnhealthy.String(), duration)
		return Result{Status: StatusUnhealthy, Error: err}
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(p.timeout))

	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("smtp", p.address, StatusUnhealthy.String(), duration)
		return Result{Status: StatusUnhealthy, Error: err}
	}

	duration := time.Since(start).Seconds()

	if !strings.HasPrefix(strings.TrimSpace(line), "220") {
		err = fmt.Errorf("unexpected SMTP banner: %q", strings.TrimSpace(line))
		metrics.RecordProbe("smtp", p.address, StatusUnhealthy.String(), duration)
		return Result{Status: StatusUnhealthy, Error: err}
	}

	metrics.RecordProbe("smtp", p.address, StatusHealthy.String(), duration)
	return Result{Status: StatusHealthy}
}
