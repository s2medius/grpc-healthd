package probe

import (
	"context"
	"net"
	"testing"
	"time"
)

// startFakeClickHouse starts a TCP listener that immediately writes one byte
// (simulating the ClickHouse server hello) and returns the address.
func startFakeClickHouse(t *testing.T, sendByte bool) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
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
				if sendByte {
					_, _ = c.Write([]byte{0x05}) // fake server hello byte
				}
			}(conn)
		}
	}()

	return ln.Addr().String()
}

func TestClickHouseProbe_Healthy(t *testing.T) {
	addr := startFakeClickHouse(t, true)
	p := NewClickHouseProbe(addr, time.Second)
	res := p.Execute(context.Background())
	if res.Status != StatusHealthy {
		t.Fatalf("expected healthy, got %s: %v", res.Status, res.Error)
	}
}

func TestClickHouseProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewClickHouseProbe("127.0.0.1:19999", 200*time.Millisecond)
	res := p.Execute(context.Background())
	if res.Status != StatusUnhealthy {
		t.Fatalf("expected unhealthy, got %s", res.Status)
	}
}

func TestNewClickHouseProbe_DefaultTimeout(t *testing.T) {
	p := NewClickHouseProbe("localhost:9000", 0)
	if p.timeout != DefaultTimeout {
		t.Fatalf("expected default timeout %v, got %v", DefaultTimeout, p.timeout)
	}
}

func TestClickHouseProbe_CustomTimeout(t *testing.T) {
	p := NewClickHouseProbe("localhost:9000", 3*time.Second)
	if p.timeout != 3*time.Second {
		t.Fatalf("expected 3s, got %v", p.timeout)
	}
}

func TestClickHouseProbe_Unhealthy_NoResponse(t *testing.T) {
	addr := startFakeClickHouse(t, false) // connects but sends nothing
	p := NewClickHouseProbe(addr, 200*time.Millisecond)
	res := p.Execute(context.Background())
	if res.Status != StatusUnhealthy {
		t.Fatalf("expected unhealthy when no server hello, got %s", res.Status)
	}
}

func TestClickHouseProbe_DurationRecorded(t *testing.T) {
	addr := startFakeClickHouse(t, true)
	p := NewClickHouseProbe(addr, time.Second)
	res := p.Execute(context.Background())
	if res.Duration <= 0 {
		t.Fatalf("expected positive duration, got %v", res.Duration)
	}
}
