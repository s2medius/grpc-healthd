package probe_test

import (
	"bufio"
	"context"
	"net"
	"testing"
	"time"

	"github.com/yourusername/grpc-healthd/internal/config"
	"github.com/yourusername/grpc-healthd/internal/probe"
)

func startFakeRedis(t *testing.T, respond string) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start fake redis: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		scanner := bufio.NewScanner(conn)
		scanner.Scan()
		conn.Write([]byte(respond))
	}()
	return ln.Addr().String()
}

func TestRedisProbe_Healthy(t *testing.T) {
	addr := startFakeRedis(t, "+PONG\r\n")
	p := probe.NewRedisProbe(config.ProbeConfig{Address: addr})
	status := p.Run(context.Background())
	if !status.Healthy {
		t.Errorf("expected healthy, got: %s", status.Message)
	}
}

func TestRedisProbe_Unhealthy_BadResponse(t *testing.T) {
	addr := startFakeRedis(t, "-ERR unknown command\r\n")
	p := probe.NewRedisProbe(config.ProbeConfig{Address: addr})
	status := p.Run(context.Background())
	if status.Healthy {
		t.Error("expected unhealthy for bad response")
	}
}

func TestRedisProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := probe.NewRedisProbe(config.ProbeConfig{Address: "127.0.0.1:19999"})
	status := p.Run(context.Background())
	if status.Healthy {
		t.Error("expected unhealthy for refused connection")
	}
}

func TestNewRedisProbe_DefaultTimeout(t *testing.T) {
	p := probe.NewRedisProbe(config.ProbeConfig{Address: "localhost:6379"})
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestRedisProbe_CustomTimeout(t *testing.T) {
	p := probe.NewRedisProbe(config.ProbeConfig{
		Address: "localhost:6379",
		Timeout: 2 * time.Second,
	})
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}
