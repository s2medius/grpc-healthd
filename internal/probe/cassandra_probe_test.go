package probe

import (
	"context"
	"net"
	"testing"
	"time"
)

// startFakeCassandra starts a minimal TCP server that responds with a
// valid CQL SUPPORTED frame header so the probe considers it healthy.
func startFakeCassandra(t *testing.T) string {
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
				buf := make([]byte, 16)
				_, _ = c.Read(buf)
				// SUPPORTED frame: version=0x83, flags=0, stream=0, opcode=6, length=0
				_, _ = c.Write([]byte{0x83, 0x00, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00})
			}(conn)
		}
	}()
	return ln.Addr().String()
}

func TestCassandraProbe_Healthy(t *testing.T) {
	addr := startFakeCassandra(t)
	p := NewCassandraProbe(addr, time.Second)
	res := p.Check(context.Background())
	if res.Status != StatusHealthy {
		t.Errorf("expected healthy, got %v: %v", res.Status, res.Error)
	}
}

func TestCassandraProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewCassandraProbe("127.0.0.1:19999", 200*time.Millisecond)
	res := p.Check(context.Background())
	if res.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %v", res.Status)
	}
}

func TestNewCassandraProbe_DefaultTimeout(t *testing.T) {
	p := NewCassandraProbe("localhost:9042", 0)
	if p.timeout != DefaultTimeout {
		t.Errorf("expected default timeout %v, got %v", DefaultTimeout, p.timeout)
	}
}

func TestCassandraProbe_CustomTimeout(t *testing.T) {
	p := NewCassandraProbe("localhost:9042", 3*time.Second)
	if p.timeout != 3*time.Second {
		t.Errorf("expected 3s, got %v", p.timeout)
	}
}

func TestCassandraProbe_DurationRecorded(t *testing.T) {
	addr := startFakeCassandra(t)
	p := NewCassandraProbe(addr, time.Second)
	// Should not panic; metrics recording is exercised implicitly.
	res := p.Check(context.Background())
	if res.Status != StatusHealthy {
		t.Errorf("expected healthy, got %v", res.Status)
	}
}
