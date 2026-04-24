package probe_test

import (
	"net"
	"testing"
	"time"

	"github.com/yourorg/grpc-healthd/internal/probe"
)

func startFakeSSH(t *testing.T, banner string) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start fake SSH server: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		conn.Write([]byte(banner))
	}()
	return ln.Addr().String()
}

func TestSSHProbe_Healthy(t *testing.T) {
	addr := startFakeSSH(t, "SSH-2.0-OpenSSH_8.9\r\n")
	p := probe.NewSSHProbe(addr, 3*time.Second)
	result := p.Probe()
	if !result.Healthy {
		t.Errorf("expected healthy, got unhealthy: %s", result.Message)
	}
}

func TestSSHProbe_Unhealthy_BadBanner(t *testing.T) {
	addr := startFakeSSH(t, "NOT-SSH-BANNER\r\n")
	p := probe.NewSSHProbe(addr, 3*time.Second)
	result := p.Probe()
	if result.Healthy {
		t.Error("expected unhealthy for bad SSH banner")
	}
}

func TestSSHProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := probe.NewSSHProbe("127.0.0.1:1", 1*time.Second)
	result := p.Probe()
	if result.Healthy {
		t.Error("expected unhealthy for connection refused")
	}
}

func TestNewSSHProbe_DefaultTimeout(t *testing.T) {
	p := probe.NewSSHProbe("localhost:22", 0)
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestSSHProbe_DurationRecorded(t *testing.T) {
	addr := startFakeSSH(t, "SSH-2.0-OpenSSH_8.9\r\n")
	p := probe.NewSSHProbe(addr, 3*time.Second)
	result := p.Probe()
	if result.Duration <= 0 {
		t.Errorf("expected positive duration, got %v", result.Duration)
	}
}
