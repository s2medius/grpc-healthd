package probe

import (
	"context"
	"fmt"
	"net"
	"time"

	"grpc-healthd/internal/metrics"
)

// MSSQLProbe checks Microsoft SQL Server availability by performing a TCP
// handshake and validating the TDS prelogin banner response.
type MSSQLProbe struct {
	address string
	timeout time.Duration
}

// NewMSSQLProbe creates a new MSSQLProbe. If timeout is zero the default
// probe timeout is used.
func NewMSSQLProbe(address string, timeout time.Duration) *MSSQLProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &MSSQLProbe{address: address, timeout: timeout}
}

// Probe dials the MSSQL server, sends a minimal TDS prelogin packet and
// checks that the server responds with a valid TDS packet header (first byte
// 0x04 = tabular result / prelogin response).
func (p *MSSQLProbe) Probe(ctx context.Context) Result {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", p.address, p.timeout)
	if err != nil {
		metrics.RecordProbe("mssql", p.address, false, time.Since(start).Seconds())
		return Result{Healthy: false, Message: fmt.Sprintf("dial failed: %v", err)}
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(p.timeout))

	// Minimal TDS 7.x prelogin packet (type=0x12, status=0x01, length=0x002F)
	prelogin := []byte{
		0x12, 0x01, 0x00, 0x2F, 0x00, 0x00, 0x00, 0x00, // TDS header
		0x00, 0x00, 0x1A, 0x00, 0x06, // VERSION token
		0x01, 0x00, 0x20, 0x00, 0x01, // ENCRYPTION token
		0x02, 0x00, 0x21, 0x00, 0x01, // INSTOPT token
		0x03, 0x00, 0x22, 0x00, 0x04, // THREADID token
		0xFF,                         // terminator
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // VERSION data
		0x00,                               // ENCRYPTION data
		0x00,                               // INSTOPT data
		0x00, 0x00, 0x00, 0x00,             // THREADID data
	}

	if _, err := conn.Write(prelogin); err != nil {
		metrics.RecordProbe("mssql", p.address, false, time.Since(start).Seconds())
		return Result{Healthy: false, Message: fmt.Sprintf("write failed: %v", err)}
	}

	header := make([]byte, 1)
	if _, err := conn.Read(header); err != nil {
		metrics.RecordProbe("mssql", p.address, false, time.Since(start).Seconds())
		return Result{Healthy: false, Message: fmt.Sprintf("read failed: %v", err)}
	}

	// TDS prelogin response type is 0x04
	if header[0] != 0x04 {
		metrics.RecordProbe("mssql", p.address, false, time.Since(start).Seconds())
		return Result{Healthy: false, Message: fmt.Sprintf("unexpected TDS response type: 0x%02X", header[0])}
	}

	duration := time.Since(start).Seconds()
	metrics.RecordProbe("mssql", p.address, true, duration)
	return Result{Healthy: true, Message: "mssql reachable"}
}
