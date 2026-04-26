package probe

import (
	"context"
	"net"
	"testing"
	"time"
)

func startFakeTimescaleDB(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start fake TimescaleDB: %v", err)
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
				// Read startup message
				buf := make([]byte, 256)
				_, _ = c.Read(buf)
				// Respond with 'R' (AuthenticationRequest)
				_, _ = c.Write([]byte{'R', 0, 0, 0, 8, 0, 0, 0, 0})
			}(conn)
		}
	}()

	return ln.Addr().String()
}

func TestTimescaleDBProbe_Healthy(t *testing.T) {
	addr := startFakeTimescaleDB(t)
	p := NewTimescaleDBProbe(addr, time.Second)
	res := p.Probe(context.Background())
	if res.Status != StatusHealthy {
		t.Errorf("expected healthy, got %s: %v", res.Status, res.Error)
	}
}

func TestTimescaleDBProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewTimescaleDBProbe("127.0.0.1:19999", 200*time.Millisecond)
	res := p.Probe(context.Background())
	if res.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", res.Status)
	}
}

func TestNewTimescaleDBProbe_DefaultTimeout(t *testing.T) {
	p := NewTimescaleDBProbe("localhost:5432", 0)
	if p.timeout != DefaultTimeout {
		t.Errorf("expected default timeout %v, got %v", DefaultTimeout, p.timeout)
	}
}

func TestTimescaleDBProbe_CustomTimeout(t *testing.T) {
	p := NewTimescaleDBProbe("localhost:5432", 3*time.Second)
	if p.timeout != 3*time.Second {
		t.Errorf("expected 3s timeout, got %v", p.timeout)
	}
}

func TestTimescaleDBProbe_DurationRecorded(t *testing.T) {
	addr := startFakeTimescaleDB(t)
	p := NewTimescaleDBProbe(addr, time.Second)
	res := p.Probe(context.Background())
	if res.Duration <= 0 {
		t.Errorf("expected positive duration, got %v", res.Duration)
	}
}
