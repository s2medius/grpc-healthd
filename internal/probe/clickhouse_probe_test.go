package probe

import (
	"context"
	"net"
	"testing"
	"time"
)

// startFakeClickHouse starts a TCP server that sends a fake ClickHouse banner.
func startFakeClickHouse(t *testing.T, banner string) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start fake clickhouse: %v", err)
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

func TestClickHouseProbe_Healthy(t *testing.T) {
	addr := startFakeClickHouse(t, "ClickHouse server version 23.8")
	p := NewClickHouseProbe(addr, time.Second)
	res := p.Probe(context.Background())
	if !res.Healthy {
		t.Errorf("expected healthy, got: %s", res.Message)
	}
}

func TestClickHouseProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewClickHouseProbe("127.0.0.1:19999", 300*time.Millisecond)
	res := p.Probe(context.Background())
	if res.Healthy {
		t.Error("expected unhealthy for refused connection")
	}
}

func TestNewClickHouseProbe_DefaultTimeout(t *testing.T) {
	p := NewClickHouseProbe("localhost:9000", 0)
	if p.timeout != DefaultTimeout {
		t.Errorf("expected default timeout %v, got %v", DefaultTimeout, p.timeout)
	}
}

func TestClickHouseProbe_CustomTimeout(t *testing.T) {
	p := NewClickHouseProbe("localhost:9000", 5*time.Second)
	if p.timeout != 5*time.Second {
		t.Errorf("expected 5s timeout, got %v", p.timeout)
	}
}

func TestClickHouseProbe_DurationRecorded(t *testing.T) {
	addr := startFakeClickHouse(t, "ClickHouse server version 23.8")
	p := NewClickHouseProbe(addr, time.Second)
	start := time.Now()
	res := p.Probe(context.Background())
	elapsed := time.Since(start)
	if !res.Healthy {
		t.Errorf("expected healthy: %s", res.Message)
	}
	if elapsed > 2*time.Second {
		t.Errorf("probe took too long: %v", elapsed)
	}
}
