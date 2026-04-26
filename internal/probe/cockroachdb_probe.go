package probe

import (
	"fmt"
	"net"
	"time"

	"github.com/your-org/grpc-healthd/internal/metrics"
)

// NewCockroachDBProbe returns a Probe that checks CockroachDB availability
// by initiating a TCP connection and validating the PostgreSQL wire protocol
// startup banner (CockroachDB speaks the PostgreSQL protocol).
func NewCockroachDBProbe(address string, timeout time.Duration) Probe {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	return &cockroachDBProbe{address: address, timeout: timeout}
}

type cockroachDBProbe struct {
	address string
	timeout time.Duration
}

func (p *cockroachDBProbe) Check() Result {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", p.address, p.timeout)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("cockroachdb", p.address, false, duration)
		return Result{
			Status:  StatusUnhealthy,
			Message: fmt.Sprintf("connection failed: %v", err),
		}
	}
	defer conn.Close()

	// Send a minimal PostgreSQL SSLRequest packet (8 bytes).
	// CockroachDB responds with 'N' (no SSL) or 'S' (SSL supported).
	sslRequest := []byte{0x00, 0x00, 0x00, 0x08, 0x04, 0xd2, 0x16, 0x2f}
	_ = conn.SetDeadline(time.Now().Add(p.timeout))
	if _, err := conn.Write(sslRequest); err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("cockroachdb", p.address, false, duration)
		return Result{
			Status:  StatusUnhealthy,
			Message: fmt.Sprintf("write failed: %v", err),
		}
	}

	buf := make([]byte, 1)
	if _, err := conn.Read(buf); err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("cockroachdb", p.address, false, duration)
		return Result{
			Status:  StatusUnhealthy,
			Message: fmt.Sprintf("read failed: %v", err),
		}
	}

	// 'N' = no SSL, 'S' = SSL — both indicate a live CockroachDB node.
	if buf[0] != 'N' && buf[0] != 'S' {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("cockroachdb", p.address, false, duration)
		return Result{
			Status:  StatusUnhealthy,
			Message: fmt.Sprintf("unexpected SSL response byte: 0x%02x", buf[0]),
		}
	}

	duration := time.Since(start).Seconds()
	metrics.RecordProbe("cockroachdb", p.address, true, duration)
	return Result{Status: StatusHealthy, Message: "cockroachdb reachable"}
}
