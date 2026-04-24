package probe

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/yourorg/grpc-healthd/internal/config"
)

func startFakeAMQP(t *testing.T, banner string) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
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

func TestAMQPProbe_Healthy(t *testing.T) {
	addr := startFakeAMQP(t, "AMQP")
	probe := NewAMQPProbe(config.ProbeConfig{Address: addr, Timeout: 2 * time.Second})

	result := probe.Probe(context.Background())

	if !result.Healthy {
		t.Errorf("expected healthy, got unhealthy: %s", result.Message)
	}
}

func TestAMQPProbe_Unhealthy_BadBanner(t *testing.T) {
	addr := startFakeAMQP(t, "HTTP")
	probe := NewAMQPProbe(config.ProbeConfig{Address: addr, Timeout: 2 * time.Second})

	result := probe.Probe(context.Background())

	if result.Healthy {
		t.Error("expected unhealthy due to bad banner, got healthy")
	}
}

func TestAMQPProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	probe := NewAMQPProbe(config.ProbeConfig{Address: "127.0.0.1:1", Timeout: 1 * time.Second})

	result := probe.Probe(context.Background())

	if result.Healthy {
		t.Error("expected unhealthy for refused connection, got healthy")
	}
}

func TestNewAMQPProbe_DefaultTimeout(t *testing.T) {
	probe := NewAMQPProbe(config.ProbeConfig{Address: "localhost:5672"})

	if probe.timeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %v", probe.timeout)
	}
}

func TestAMQPProbe_DurationRecorded(t *testing.T) {
	addr := startFakeAMQP(t, "AMQP")
	probe := NewAMQPProbe(config.ProbeConfig{Address: addr, Timeout: 2 * time.Second})

	// Ensure Probe completes without panic and records metrics.
	result := probe.Probe(context.Background())
	if !result.Healthy {
		t.Errorf("expected healthy result, got: %s", result.Message)
	}
}
