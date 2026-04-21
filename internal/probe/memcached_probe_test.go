package probe

import (
	"fmt"
	"net"
	"testing"
	"time"
)

// startFakeMemcached starts a fake Memcached server that responds to
// the VERSION command with the provided banner.
func startFakeMemcached(t *testing.T, banner string) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start fake memcached: %v", err)
	}
	t.Cleanup(func() { ln.Close() })

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		buf := make([]byte, 32)
		_, _ = conn.Read(buf)
		fmt.Fprint(conn, banner)
	}()

	return ln.Addr().String()
}

func TestMemcachedProbe_Healthy(t *testing.T) {
	addr := startFakeMemcached(t, "VERSION 1.6.12\r\n")
	p := NewMemcachedProbe(addr, time.Second)
	result := p.Check()
	if result.Status != StatusHealthy {
		t.Errorf("expected healthy, got %s: %v", result.Status, result.Error)
	}
}

func TestMemcachedProbe_Unhealthy_BadResponse(t *testing.T) {
	addr := startFakeMemcached(t, "ERROR\r\n")
	p := NewMemcachedProbe(addr, time.Second)
	result := p.Check()
	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", result.Status)
	}
	if result.Error == nil {
		t.Error("expected non-nil error for bad response")
	}
}

func TestMemcachedProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewMemcachedProbe("127.0.0.1:19999", 200*time.Millisecond)
	result := p.Check()
	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", result.Status)
	}
}

func TestNewMemcachedProbe_DefaultTimeout(t *testing.T) {
	p := NewMemcachedProbe("localhost:11211", 0)
	if p.timeout != DefaultTimeout {
		t.Errorf("expected default timeout %v, got %v", DefaultTimeout, p.timeout)
	}
}

func TestMemcachedProbe_CustomTimeout(t *testing.T) {
	custom := 3 * time.Second
	p := NewMemcachedProbe("localhost:11211", custom)
	if p.timeout != custom {
		t.Errorf("expected %v, got %v", custom, p.timeout)
	}
}

func TestMemcachedProbe_DurationRecorded(t *testing.T) {
	addr := startFakeMemcached(t, "VERSION 1.6.12\r\n")
	p := NewMemcachedProbe(addr, time.Second)
	// Should not panic; metrics recording is a side effect.
	p.Check()
}
