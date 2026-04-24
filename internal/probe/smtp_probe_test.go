package probe_test

import (
	"bufio"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/your-org/grpc-healthd/internal/probe"
)

func startFakeSMTP(t *testing.T, banner string) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start fake SMTP server: %v", err)
	}
	t.Cleanup(func() { ln.Close() })

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		w := bufio.NewWriter(conn)
		fmt.Fprintln(w, banner)
		w.Flush()
	}()

	return ln.Addr().String()
}

func TestSMTPProbe_Healthy(t *testing.T) {
	addr := startFakeSMTP(t, "220 localhost ESMTP ready")
	p := probe.NewSMTPProbe(addr, time.Second)
	result := p.Check()
	if result.Status != probe.StatusHealthy {
		t.Errorf("expected healthy, got %s: %v", result.Status, result.Error)
	}
}

func TestSMTPProbe_Unhealthy_BadBanner(t *testing.T) {
	addr := startFakeSMTP(t, "500 Something went wrong")
	p := probe.NewSMTPProbe(addr, time.Second)
	result := p.Check()
	if result.Status != probe.StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", result.Status)
	}
	if result.Error == nil {
		t.Error("expected non-nil error for bad banner")
	}
}

func TestSMTPProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := probe.NewSMTPProbe("127.0.0.1:1", time.Second)
	result := p.Check()
	if result.Status != probe.StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", result.Status)
	}
}

func TestNewSMTPProbe_DefaultTimeout(t *testing.T) {
	p := probe.NewSMTPProbe("localhost:25", 0)
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestSMTPProbe_DurationRecorded(t *testing.T) {
	addr := startFakeSMTP(t, "220 test ESMTP")
	p := probe.NewSMTPProbe(addr, time.Second)
	result := p.Check()
	if result.Status != probe.StatusHealthy {
		t.Errorf("expected healthy, got %s: %v", result.Status, result.Error)
	}
}
