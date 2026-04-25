package probe

import (
	"bufio"
	"net"
	"testing"
	"time"
)

func startFakeRedisSentinel(t *testing.T, response string) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
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
		_, _ = conn.Write([]byte(response + "\r\n"))
	}()
	return ln.Addr().String()
}

func TestRedisSentinelProbe_Healthy(t *testing.T) {
	addr := startFakeRedisSentinel(t, "+PONG")
	p := NewRedisSentinelProbe(addr, time.Second)
	res := p.Check()
	if res.Status != StatusHealthy {
		t.Errorf("expected healthy, got %s: %v", res.Status, res.Error)
	}
}

func TestRedisSentinelProbe_Unhealthy_BadResponse(t *testing.T) {
	addr := startFakeRedisSentinel(t, "-ERR unknown command")
	p := NewRedisSentinelProbe(addr, time.Second)
	res := p.Check()
	if res.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", res.Status)
	}
}

func TestRedisSentinelProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewRedisSentinelProbe("127.0.0.1:1", time.Second)
	res := p.Check()
	if res.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", res.Status)
	}
}

func TestNewRedisSentinelProbe_DefaultTimeout(t *testing.T) {
	p := NewRedisSentinelProbe("127.0.0.1:26379", 0)
	if p.timeout != DefaultTimeout {
		t.Errorf("expected default timeout %v, got %v", DefaultTimeout, p.timeout)
	}
}

func TestRedisSentinelProbe_CustomTimeout(t *testing.T) {
	p := NewRedisSentinelProbe("127.0.0.1:26379", 3*time.Second)
	if p.timeout != 3*time.Second {
		t.Errorf("expected 3s, got %v", p.timeout)
	}
}

func TestRedisSentinelProbe_DurationRecorded(t *testing.T) {
	addr := startFakeRedisSentinel(t, "+PONG")
	p := NewRedisSentinelProbe(addr, time.Second)
	res := p.Check()
	if res.Status != StatusHealthy {
		t.Errorf("expected healthy, got %s", res.Status)
	}
}
