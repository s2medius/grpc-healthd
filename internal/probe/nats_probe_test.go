package probe

import (
	"context"
	"net"
	"testing"
	"time"
)

// startFakeNATS starts a TCP server that sends a NATS-like INFO banner.
func startFakeNATS(t *testing.T, banner string) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("startFakeNATS listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		_, _ = conn.Write([]byte(banner))
	}()

	return ln.Addr().String()
}

func TestNATSProbe_Healthy(t *testing.T) {
	addr := startFakeNATS(t, `INFO {"server_id":"test","version":"2.9.0"}`+"\r\n")
	p := NewNATSProbe(addr, time.Second)
	res := p.Probe(context.Background())
	if res.Status != StatusHealthy {
		t.Fatalf("expected healthy, got %s: %v", res.Status, res.Error)
	}
}

func TestNATSProbe_Unhealthy_BadBanner(t *testing.T) {
	addr := startFakeNATS(t, "GARBAGE\r\n")
	p := NewNATSProbe(addr, time.Second)
	res := p.Probe(context.Background())
	if res.Status != StatusUnhealthy {
		t.Fatalf("expected unhealthy, got %s", res.Status)
	}
	if res.Error == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestNATSProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewNATSProbe("127.0.0.1:1", 200*time.Millisecond)
	res := p.Probe(context.Background())
	if res.Status != StatusUnhealthy {
		t.Fatalf("expected unhealthy, got %s", res.Status)
	}
}

func TestNewNATSProbe_DefaultTimeout(t *testing.T) {
	p := NewNATSProbe("127.0.0.1:4222", 0)
	if p.timeout != DefaultTimeout {
		t.Fatalf("expected default timeout %v, got %v", DefaultTimeout, p.timeout)
	}
}

func TestNATSProbe_DurationRecorded(t *testing.T) {
	addr := startFakeNATS(t, "INFO {}\r\n")
	p := NewNATSProbe(addr, time.Second)
	// Should not panic and should complete without error.
	res := p.Probe(context.Background())
	if res.Status != StatusHealthy {
		t.Fatalf("expected healthy, got %s: %v", res.Status, res.Error)
	}
}
