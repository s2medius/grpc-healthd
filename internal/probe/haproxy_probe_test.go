package probe

import (
	"bufio"
	"net"
	"testing"
	"time"
)

func startFakeHAProxy(t *testing.T, banner string) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start fake haproxy: %v", err)
	}
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		scanner := bufio.NewScanner(conn)
		scanner.Scan() // read "show info\n"
		_, _ = conn.Write([]byte(banner + "\n"))
		ln.Close()
	}()
	return ln.Addr().String()
}

func TestHAProxyProbe_Healthy(t *testing.T) {
	addr := startFakeHAProxy(t, "Name: HAProxy")
	p := NewHAProxyProbe(addr, time.Second)
	result := p.Probe()
	if result.Status != StatusHealthy {
		t.Errorf("expected healthy, got %s: %v", result.Status, result.Error)
	}
}

func TestHAProxyProbe_Unhealthy_BadBanner(t *testing.T) {
	addr := startFakeHAProxy(t, "unexpected garbage")
	p := NewHAProxyProbe(addr, time.Second)
	result := p.Probe()
	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", result.Status)
	}
}

func TestHAProxyProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewHAProxyProbe("127.0.0.1:19999", 200*time.Millisecond)
	result := p.Probe()
	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", result.Status)
	}
	if result.Error == nil {
		t.Error("expected non-nil error")
	}
}

func TestNewHAProxyProbe_DefaultTimeout(t *testing.T) {
	p := NewHAProxyProbe("127.0.0.1:1936", 0)
	if p.timeout != defaultTimeout {
		t.Errorf("expected default timeout %v, got %v", defaultTimeout, p.timeout)
	}
}

func TestHAProxyProbe_DurationRecorded(t *testing.T) {
	addr := startFakeHAProxy(t, "Name: HAProxy")
	p := NewHAProxyProbe(addr, time.Second)
	result := p.Probe()
	if result.Status != StatusHealthy {
		t.Errorf("expected healthy, got %s", result.Status)
	}
}
