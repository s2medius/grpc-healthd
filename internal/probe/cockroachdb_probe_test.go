package probe

import (
	"net"
	"testing"
	"time"
)

// startFakeCockroachDB spins up a minimal TCP server that mimics the
// CockroachDB SSL-negotiation response ('N' = no SSL).
func startFakeCockroachDB(t *testing.T, response byte) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		buf := make([]byte, 8)
		_, _ = conn.Read(buf)
		_, _ = conn.Write([]byte{response})
	}()
	t.Cleanup(func() { ln.Close() })
	return ln.Addr().String()
}

func TestCockroachDBProbe_Healthy_NoSSL(t *testing.T) {
	addr := startFakeCockroachDB(t, 'N')
	p := NewCockroachDBProbe(addr, time.Second)
	res := p.Check()
	if res.Status != StatusHealthy {
		t.Fatalf("expected healthy, got %s: %s", res.Status, res.Message)
	}
}

func TestCockroachDBProbe_Healthy_SSL(t *testing.T) {
	addr := startFakeCockroachDB(t, 'S')
	p := NewCockroachDBProbe(addr, time.Second)
	res := p.Check()
	if res.Status != StatusHealthy {
		t.Fatalf("expected healthy, got %s: %s", res.Status, res.Message)
	}
}

func TestCockroachDBProbe_Unhealthy_BadResponse(t *testing.T) {
	addr := startFakeCockroachDB(t, 0xFF)
	p := NewCockroachDBProbe(addr, time.Second)
	res := p.Check()
	if res.Status != StatusUnhealthy {
		t.Fatalf("expected unhealthy, got %s", res.Status)
	}
}

func TestCockroachDBProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewCockroachDBProbe("127.0.0.1:19999", 200*time.Millisecond)
	res := p.Check()
	if res.Status != StatusUnhealthy {
		t.Fatalf("expected unhealthy, got %s", res.Status)
	}
}

func TestNewCockroachDBProbe_DefaultTimeout(t *testing.T) {
	p := NewCockroachDBProbe("localhost:26257", 0).(*cockroachDBProbe)
	if p.timeout != DefaultTimeout {
		t.Fatalf("expected default timeout %v, got %v", DefaultTimeout, p.timeout)
	}
}

func TestCockroachDBProbe_CustomTimeout(t *testing.T) {
	p := NewCockroachDBProbe("localhost:26257", 3*time.Second).(*cockroachDBProbe)
	if p.timeout != 3*time.Second {
		t.Fatalf("expected 3s timeout, got %v", p.timeout)
	}
}
