package probe_test

import (
	"context"
	"net"
	"testing"
	"time"

	"grpc-healthd/internal/probe"
)

func TestICMPProbe_Healthy(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	p := probe.NewICMPProbe(ln.Addr().String(), time.Second)
	result := p.Probe(context.Background())
	if !result.Healthy {
		t.Fatalf("expected healthy, got: %s", result.Message)
	}
}

func TestICMPProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	// Use a port that is not listening
	p := probe.NewICMPProbe("127.0.0.1:19999", time.Second)
	result := p.Probe(context.Background())
	if result.Healthy {
		t.Fatal("expected unhealthy for refused connection")
	}
}

func TestNewICMPProbe_DefaultTimeout(t *testing.T) {
	p := probe.NewICMPProbe("127.0.0.1:80", 0)
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestICMPProbe_NoPortFallback(t *testing.T) {
	// Providing host without port should not panic; connection will fail gracefully
	p := probe.NewICMPProbe("192.0.2.1", 200*time.Millisecond) // TEST-NET, unreachable
	result := p.Probe(context.Background())
	if result.Healthy {
		t.Fatal("expected unhealthy for unreachable host")
	}
	if result.Duration == 0 {
		t.Fatal("expected non-zero duration")
	}
}

func TestICMPProbe_DurationRecorded(t *testing.T) {
	p := probe.NewICMPProbe("127.0.0.1:19998", 200*time.Millisecond)
	result := p.Probe(context.Background())
	if result.Duration == 0 {
		t.Error("expected non-zero duration")
	}
}
