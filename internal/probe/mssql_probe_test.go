package probe

import (
	"context"
	"net"
	"testing"
	"time"
)

// startFakeMSSQL starts a TCP server that responds with a valid TDS prelogin
// response header (0x04) followed by enough padding bytes.
func startFakeMSSQL(t *testing.T, healthy bool) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		buf := make([]byte, 128)
		_, _ = conn.Read(buf)
		if healthy {
			// Valid TDS prelogin response: type 0x04
			_, _ = conn.Write([]byte{0x04, 0x01, 0x00, 0x0B, 0x00, 0x00, 0x01, 0x00, 0xFF})
		} else {
			// Invalid TDS type
			_, _ = conn.Write([]byte{0xFF})
		}
	}()

	return ln.Addr().String()
}

func TestMSSQLProbe_Healthy(t *testing.T) {
	addr := startFakeMSSQL(t, true)
	p := NewMSSQLProbe(addr, time.Second)
	res := p.Probe(context.Background())
	if !res.Healthy {
		t.Fatalf("expected healthy, got: %s", res.Message)
	}
}

func TestMSSQLProbe_Unhealthy_BadResponse(t *testing.T) {
	addr := startFakeMSSQL(t, false)
	p := NewMSSQLProbe(addr, time.Second)
	res := p.Probe(context.Background())
	if res.Healthy {
		t.Fatal("expected unhealthy for bad TDS response")
	}
}

func TestMSSQLProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewMSSQLProbe("127.0.0.1:1", 200*time.Millisecond)
	res := p.Probe(context.Background())
	if res.Healthy {
		t.Fatal("expected unhealthy for refused connection")
	}
}

func TestNewMSSQLProbe_DefaultTimeout(t *testing.T) {
	p := NewMSSQLProbe("localhost:1433", 0)
	if p.timeout != DefaultTimeout {
		t.Errorf("expected default timeout %v, got %v", DefaultTimeout, p.timeout)
	}
}

func TestMSSQLProbe_CustomTimeout(t *testing.T) {
	p := NewMSSQLProbe("localhost:1433", 3*time.Second)
	if p.timeout != 3*time.Second {
		t.Errorf("expected 3s timeout, got %v", p.timeout)
	}
}

func TestMSSQLProbe_DurationRecorded(t *testing.T) {
	addr := startFakeMSSQL(t, true)
	p := NewMSSQLProbe(addr, time.Second)
	res := p.Probe(context.Background())
	if !res.Healthy {
		t.Fatalf("expected healthy: %s", res.Message)
	}
}
