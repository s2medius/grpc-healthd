package probe

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/grpc-healthd/internal/metrics"
)

// FTPProbe checks health by connecting to an FTP server and verifying the banner.
type FTPProbe struct {
	address string
	timeout time.Duration
}

// NewFTPProbe creates a new FTPProbe. If timeout is zero, defaultTimeout is used.
func NewFTPProbe(address string, timeout time.Duration) *FTPProbe {
	if timeout == 0 {
		timeout = defaultTimeout
	}
	return &FTPProbe{address: address, timeout: timeout}
}

// Check connects to the FTP server and validates the 220 service-ready banner.
func (p *FTPProbe) Check() Result {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", p.address, p.timeout)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("ftp", p.address, false, duration)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("ftp connect: %w", err)}
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(p.timeout))

	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("ftp", p.address, false, duration)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("ftp read banner: %w", err)}
	}

	duration := time.Since(start).Seconds()
	banner := string(buf[:n])

	if !strings.HasPrefix(banner, "220") {
		metrics.RecordProbe("ftp", p.address, false, duration)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("unexpected ftp banner: %q", banner)}
	}

	metrics.RecordProbe("ftp", p.address, true, duration)
	return Result{Status: StatusHealthy}
}
