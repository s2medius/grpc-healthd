package probe

import (
	"context"
	"net"
	"testing"
	"time"
)

// startFakeRabbitMQ starts a fake TCP server that sends an AMQP-like banner.
func startFakeRabbitMQ(t *testing.T, banner string) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start fake rabbitmq: %v", err)
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
				_, _ = c.Write([]byte(banner + "\n"))
			}(conn)
		}
	}()

	return ln.Addr().String()
}

func TestRabbitMQProbe_Healthy(t *testing.T) {
	addr := startFakeRabbitMQ(t, "AMQP 0-9-1 ready")
	p := NewRabbitMQProbe(addr, time.Second)
	result := p.Check(context.Background())
	if result.Status != StatusHealthy {
		t.Errorf("expected healthy, got %v: %v", result.Status, result.Error)
	}
}

func TestRabbitMQProbe_Unhealthy_BadBanner(t *testing.T) {
	addr := startFakeRabbitMQ(t, "NOT A RABBIT")
	p := NewRabbitMQProbe(addr, time.Second)
	result := p.Check(context.Background())
	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %v", result.Status)
	}
}

func TestRabbitMQProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewRabbitMQProbe("127.0.0.1:1", 200*time.Millisecond)
	result := p.Check(context.Background())
	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %v", result.Status)
	}
}

func TestNewRabbitMQProbe_DefaultTimeout(t *testing.T) {
	p := NewRabbitMQProbe("localhost:5672", 0)
	if p.timeout != DefaultTimeout {
		t.Errorf("expected default timeout %v, got %v", DefaultTimeout, p.timeout)
	}
}

func TestRabbitMQProbe_CustomTimeout(t *testing.T) {
	p := NewRabbitMQProbe("localhost:5672", 3*time.Second)
	if p.timeout != 3*time.Second {
		t.Errorf("expected 3s timeout, got %v", p.timeout)
	}
}

func TestRabbitMQProbe_DurationRecorded(t *testing.T) {
	addr := startFakeRabbitMQ(t, "AMQP 0-9-1")
	p := NewRabbitMQProbe(addr, time.Second)
	// Should not panic; metrics recording is a side effect
	_ = p.Check(context.Background())
}
