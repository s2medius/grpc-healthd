package probe

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/grpc-healthd/internal/metrics"
)

// TCPBannerProbe connects to a TCP address and checks that the server
// banner (first line of response) contains an expected substring.
type TCPBannerProbe struct {
	address string
	banner  string
	timeout time.Duration
}

// NewTCPBannerProbe creates a new TCPBannerProbe.
// If timeout is zero, DefaultTimeout is used.
func NewTCPBannerProbe(address, banner string, timeout time.Duration) *TCPBannerProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &TCPBannerProbe{
		address: address,
		banner:  banner,
		timeout: timeout,
	}
}

// Probe connects to the TCP address, reads the first 256 bytes, and
// verifies the expected banner substring is present.
func (p *TCPBannerProbe) Probe() Result {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", p.address, p.timeout)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("tcp_banner", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("connection failed: %v", err)}
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(p.timeout))

	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	duration := time.Since(start).Seconds()

	if err != nil && n == 0 {
		metrics.RecordProbe("tcp_banner", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("read banner failed: %v", err)}
	}

	recv := string(buf[:n])
	if !strings.Contains(recv, p.banner) {
		metrics.RecordProbe("tcp_banner", p.address, false, duration)
		return Result{
			Healthy: false,
			Message: fmt.Sprintf("banner mismatch: expected %q in %q", p.banner, recv),
		}
	}

	metrics.RecordProbe("tcp_banner", p.address, true, duration)
	return Result{Healthy: true, Message: "banner matched"}
}
