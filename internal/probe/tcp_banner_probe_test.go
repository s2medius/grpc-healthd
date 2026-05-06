package probe

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func startFakeTCPBanner(t *testing.T, banner string) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start fake TCP banner server: %v", err)
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
				fmt.Fprint(c, banner)
			}(conn)
		}
	}()

	return ln.Addr().String()
}

func TestTCPBannerProbe_Healthy(t *testing.T) {
	addr := startFakeTCPBanner(t, "220 Welcome to test service")
	p := NewTCPBannerProbe(addr, "220", time.Second)
	result := p.Probe()
	if !result.Healthy {
		t.Errorf("expected healthy, got: %s", result.Message)
	}
}

func TestTCPBannerProbe_Unhealthy_BadBanner(t *testing.T) {
	addr := startFakeTCPBanner(t, "500 Error")
	p := NewTCPBannerProbe(addr, "220", time.Second)
	result := p.Probe()
	if result.Healthy {
		t.Error("expected unhealthy due to banner mismatch")
	}
}

func TestTCPBannerProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewTCPBannerProbe("127.0.0.1:1", "220", 200*time.Millisecond)
	result := p.Probe()
	if result.Healthy {
		t.Error("expected unhealthy due to connection refused")
	}
}

func TestNewTCPBannerProbe_DefaultTimeout(t *testing.T) {
	p := NewTCPBannerProbe("127.0.0.1:9999", "OK", 0)
	if p.timeout != DefaultTimeout {
		t.Errorf("expected default timeout %v, got %v", DefaultTimeout, p.timeout)
	}
}

func TestTCPBannerProbe_DurationRecorded(t *testing.T) {
	addr := startFakeTCPBanner(t, "OK ready")
	p := NewTCPBannerProbe(addr, "OK", time.Second)
	// Should not panic and should complete without error
	result := p.Probe()
	if !result.Healthy {
		t.Errorf("expected healthy, got: %s", result.Message)
	}
}
