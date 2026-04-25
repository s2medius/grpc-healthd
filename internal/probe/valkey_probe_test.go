package probe

import (
	"bufio"
	"net"
	"testing"
	"time"
)

func startFakeValkey(t *testing.T, response string) string {
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
		// consume the PING command lines
		for scanner.Scan() {
			line := scanner.Text()
			if line == "PING" {
				break
			}
		}
		_, _ = conn.Write([]byte(response + "\r\n"))
	}()
	return ln.Addr().String()
}

func TestValkeyProbe_Healthy(t *testing.T) {
	addr := startFakeValkey(t, "+PONG")
	p := NewValkeyProbe(addr, time.Second)
	res := p.Check()
	if res.Status != StatusHealthy {
		t.Errorf("expected healthy, got %s: %v", res.Status, res.Error)
	}
}

func TestValkeyProbe_Unhealthy_BadResponse(t *testing.T) {
	addr := startFakeValkey(t, "-ERR unknown")
	p := NewValkeyProbe(addr, time.Second)
	res := p.Check()
	if res.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", res.Status)
	}
}

func TestValkeyProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewValkeyProbe("127.0.0.1:1", time.Second)
	res := p.Check()
	if res.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", res.Status)
	}
}

func TestNewValkeyProbe_DefaultTimeout(t *testing.T) {
	p := NewValkeyProbe("127.0.0.1:6379", 0)
	if p.timeout != DefaultTimeout {
		t.Errorf("expected default timeout %v, got %v", DefaultTimeout, p.timeout)
	}
}

func TestValkeyProbe_DurationRecorded(t *testing.T) {
	addr := startFakeValkey(t, "+PONG")
	p := NewValkeyProbe(addr, time.Second)
	res := p.Check()
	if res.Status != StatusHealthy {
		t.Errorf("expected healthy, got %s", res.Status)
	}
}
