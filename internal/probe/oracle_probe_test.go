package probe

import (
	"context"
	"net"
	"testing"
	"time"
)

// startFakeOracle starts a TCP server that sends a single byte greeting,
// simulating a minimal Oracle TNS listener.
func startFakeOracle(t *testing.T, greeting []byte) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("startFakeOracle: listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				_, _ = c.Write(greeting)
			}(conn)
		}
	}()
	return ln.Addr().String()
}

func TestOracleProbe_Healthy(t *testing.T) {
	addr := startFakeOracle(t, []byte{0x00})
	p := NewOracleProbe(addr, time.Second)
	res := p.Check(context.Background())
	if res.Status != StatusHealthy {
		t.Errorf("expected healthy, got %s: %v", res.Status, res.Error)
	}
}

func TestOracleProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewOracleProbe("127.0.0.1:1", 200*time.Millisecond)
	res := p.Check(context.Background())
	if res.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", res.Status)
	}
}

func TestNewOracleProbe_DefaultTimeout(t *testing.T) {
	p := NewOracleProbe("localhost:1521", 0)
	if p.timeout != DefaultTimeout {
		t.Errorf("expected DefaultTimeout %v, got %v", DefaultTimeout, p.timeout)
	}
}

func TestOracleProbe_CustomTimeout(t *testing.T) {
	p := NewOracleProbe("localhost:1521", 3*time.Second)
	if p.timeout != 3*time.Second {
		t.Errorf("expected 3s, got %v", p.timeout)
	}
}

func TestOracleProbe_DurationRecorded(t *testing.T) {
	addr := startFakeOracle(t, []byte{0xFF})
	p := NewOracleProbe(addr, time.Second)
	start := time.Now()
	res := p.Check(context.Background())
	elapsed := time.Since(start)
	if res.Status != StatusHealthy {
		t.Fatalf("expected healthy: %v", res.Error)
	}
	if elapsed < 0 {
		t.Error("duration should be non-negative")
	}
}
