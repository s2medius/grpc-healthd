package probe

import (
	"context"
	"net"
	"testing"
	"time"
)

func TestStatusString(t *testing.T) {
	cases := []struct {
		status Status
		want   string
	}{
		{StatusHealthy, "healthy"},
		{StatusUnhealthy, "unhealthy"},
		{StatusUnknown, "unknown"},
	}
	for _, c := range cases {
		if got := c.status.String(); got != c.want {
			t.Errorf("Status(%d).String() = %q, want %q", c.status, got, c.want)
		}
	}
}

func TestTCPProbe_Healthy(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	p := NewTCPProbe("test-tcp", ln.Addr().String(), time.Second)
	if p.Name() != "test-tcp" {
		t.Errorf("Name() = %q, want %q", p.Name(), "test-tcp")
	}

	result := p.Check(context.Background())
	if result.Status != StatusHealthy {
		t.Errorf("expected healthy, got %s: %v", result.Status, result.Err)
	}
	if result.Latency <= 0 {
		t.Error("expected positive latency")
	}
}

func TestTCPProbe_Unhealthy(t *testing.T) {
	// Use a port that is not listening.
	p := NewTCPProbe("test-tcp-fail", "127.0.0.1:1", 200*time.Millisecond)

	result := p.Check(context.Background())
	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", result.Status)
	}
	if result.Err == nil {
		t.Error("expected non-nil error")
	}
}

func TestNewTCPProbe_DefaultTimeout(t *testing.T) {
	p := NewTCPProbe("default", "localhost:80", 0)
	if p.timeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %s", p.timeout)
	}
}
