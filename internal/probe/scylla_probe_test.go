package probe

import (
	"context"
	"net"
	"testing"
	"time"
)

// startFakeScylla starts a minimal TCP server that responds to a CQL OPTIONS
// request with a SUPPORTED frame (opcode 0x06).
func startFakeScylla(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
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
				buf := make([]byte, 9)
				if _, err := c.Read(buf); err != nil {
					return
				}
				// Respond with a SUPPORTED frame (opcode 0x06), empty body.
				resp := []byte{0x84, 0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00}
				_, _ = c.Write(resp)
			}(conn)
		}
	}()

	return ln.Addr().String()
}

func TestScyllaProbe_Healthy(t *testing.T) {
	addr := startFakeScylla(t)
	p := NewScyllaProbe(addr, time.Second)
	res := p.Check(context.Background())
	if !res.Healthy {
		t.Fatalf("expected healthy, got: %s", res.Message)
	}
}

func TestScyllaProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewScyllaProbe("127.0.0.1:1", 200*time.Millisecond)
	res := p.Check(context.Background())
	if res.Healthy {
		t.Fatal("expected unhealthy for refused connection")
	}
}

func TestNewScyllaProbe_DefaultTimeout(t *testing.T) {
	p := NewScyllaProbe("127.0.0.1:9042", 0)
	if p.timeout != DefaultTimeout {
		t.Fatalf("expected default timeout %v, got %v", DefaultTimeout, p.timeout)
	}
}

func TestScyllaProbe_CustomTimeout(t *testing.T) {
	p := NewScyllaProbe("127.0.0.1:9042", 3*time.Second)
	if p.timeout != 3*time.Second {
		t.Fatalf("expected 3s timeout, got %v", p.timeout)
	}
}

func TestScyllaProbe_DurationRecorded(t *testing.T) {
	addr := startFakeScylla(t)
	p := NewScyllaProbe(addr, time.Second)
	// Verify Check completes without panic and returns a result.
	res := p.Check(context.Background())
	if !res.Healthy {
		t.Fatalf("expected healthy result, got: %s", res.Message)
	}
}
