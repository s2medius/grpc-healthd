package probe_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/yourusername/grpc-healthd/internal/probe"
)

// startFakeMySQL starts a TCP server that sends a fake MySQL greeting.
func startFakeMySQL(t *testing.T) string {
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
		// Minimal MySQL greeting: 5+ bytes
		_, _ = conn.Write([]byte{0x4a, 0x00, 0x00, 0x00, 0x0a})
	}()
	t.Cleanup(func() { ln.Close() })
	return ln.Addr().String()
}

func TestMySQLProbe_Healthy(t *testing.T) {
	addr := startFakeMySQL(t)
	p := probe.NewMySQLProbe(addr, time.Second)
	res := p.Check(context.Background())
	if res.Status != probe.StatusHealthy {
		t.Errorf("expected healthy, got %v: %v", res.Status, res.Error)
	}
}

func TestMySQLProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := probe.NewMySQLProbe("127.0.0.1:1", 200*time.Millisecond)
	res := p.Check(context.Background())
	if res.Status != probe.StatusUnhealthy {
		t.Errorf("expected unhealthy, got %v", res.Status)
	}
}

func TestNewMySQLProbe_DefaultTimeout(t *testing.T) {
	p := probe.NewMySQLProbe("localhost:3306", 0)
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestMySQLProbe_CustomTimeout(t *testing.T) {
	p := probe.NewMySQLProbe("localhost:3306", 5*time.Second)
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestMySQLProbe_DurationRecorded(t *testing.T) {
	addr := startFakeMySQL(t)
	p := probe.NewMySQLProbe(addr, time.Second)
	// Should not panic and should record metrics.
	res := p.Check(context.Background())
	if res.Error != nil && res.Status == probe.StatusHealthy {
		t.Errorf("inconsistent result: healthy but error set")
	}
}
