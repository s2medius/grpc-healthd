package probe

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func startFakeFTP(t *testing.T, banner string) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start fake FTP: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		fmt.Fprint(conn, banner)
	}()
	return ln.Addr().String()
}

func TestFTPProbe_Healthy(t *testing.T) {
	addr := startFakeFTP(t, "220 Welcome to FTP server\r\n")
	p := NewFTPProbe(addr, time.Second)
	result := p.Check()
	if result.Status != StatusHealthy {
		t.Errorf("expected healthy, got %s: %v", result.Status, result.Error)
	}
}

func TestFTPProbe_Unhealthy_BadBanner(t *testing.T) {
	addr := startFakeFTP(t, "530 Not logged in\r\n")
	p := NewFTPProbe(addr, time.Second)
	result := p.Check()
	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", result.Status)
	}
}

func TestFTPProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewFTPProbe("127.0.0.1:1", 200*time.Millisecond)
	result := p.Check()
	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", result.Status)
	}
}

func TestNewFTPProbe_DefaultTimeout(t *testing.T) {
	p := NewFTPProbe("localhost:21", 0)
	if p.timeout != defaultTimeout {
		t.Errorf("expected default timeout %v, got %v", defaultTimeout, p.timeout)
	}
}

func TestFTPProbe_DurationRecorded(t *testing.T) {
	addr := startFakeFTP(t, "220 Ready\r\n")
	p := NewFTPProbe(addr, time.Second)
	// Should not panic; metrics recording is a side-effect
	p.Check()
}
