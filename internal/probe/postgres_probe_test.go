package probe_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/yourusername/grpc-healthd/internal/config"
	"github.com/yourusername/grpc-healthd/internal/probe"
)

// startFakePostgres starts a TCP server that responds with 'N' to any data.
func startFakePostgres(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start fake postgres: %v", err)
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
				buf := make([]byte, 8)
				c.Read(buf) //nolint:errcheck
				c.Write([]byte("N")) //nolint:errcheck
			}(conn)
		}
	}()

	return ln.Addr().String()
}

func TestPostgresProbe_Healthy(t *testing.T) {
	addr := startFakePostgres(t)
	p := probe.NewPostgresProbe(config.ProbeConfig{Address: addr})

	result := p.Execute(context.Background())
	if result.Status != probe.StatusHealthy {
		t.Errorf("expected healthy, got %s: %s", result.Status, result.Message)
	}
}

func TestPostgresProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := probe.NewPostgresProbe(config.ProbeConfig{Address: "127.0.0.1:19999"})

	result := p.Execute(context.Background())
	if result.Status != probe.StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", result.Status)
	}
}

func TestNewPostgresProbe_DefaultTimeout(t *testing.T) {
	p := probe.NewPostgresProbe(config.ProbeConfig{Address: "localhost:5432"})
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestPostgresProbe_CustomTimeout(t *testing.T) {
	p := probe.NewPostgresProbe(config.ProbeConfig{
		Address: "127.0.0.1:19999",
		Timeout: 100 * time.Millisecond,
	})

	start := time.Now()
	result := p.Execute(context.Background())
	elapsed := time.Since(start)

	if result.Status != probe.StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", result.Status)
	}
	if elapsed > 2*time.Second {
		t.Errorf("probe took too long: %v", elapsed)
	}
}

func TestPostgresProbe_DurationRecorded(t *testing.T) {
	addr := startFakePostgres(t)
	p := probe.NewPostgresProbe(config.ProbeConfig{Address: addr})

	result := p.Execute(context.Background())
	if result.Duration <= 0 {
		t.Errorf("expected positive duration, got %v", result.Duration)
	}
}
