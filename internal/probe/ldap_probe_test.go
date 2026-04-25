package probe

import (
	"context"
	"net"
	"testing"
	"time"
)

func startFakeLDAP(t *testing.T, banner []byte) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start fake LDAP: %v", err)
	}
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		_, _ = conn.Write(banner)
	}()
	t.Cleanup(func() { ln.Close() })
	return ln.Addr().String()
}

func TestLDAPProbe_Healthy(t *testing.T) {
	addr := startFakeLDAP(t, []byte{0x30, 0x00})
	p := NewLDAPProbe(addr, time.Second)
	res := p.Probe(context.Background())
	if !res.Healthy {
		t.Fatalf("expected healthy, got: %s", res.Message)
	}
}

func TestLDAPProbe_Unhealthy_BadBanner(t *testing.T) {
	addr := startFakeLDAP(t, []byte{0xFF})
	p := NewLDAPProbe(addr, time.Second)
	res := p.Probe(context.Background())
	if res.Healthy {
		t.Fatal("expected unhealthy for bad banner byte")
	}
}

func TestLDAPProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewLDAPProbe("127.0.0.1:1", time.Second)
	res := p.Probe(context.Background())
	if res.Healthy {
		t.Fatal("expected unhealthy for refused connection")
	}
}

func TestNewLDAPProbe_DefaultTimeout(t *testing.T) {
	p := NewLDAPProbe("localhost:389", 0)
	if p.timeout != defaultTimeout {
		t.Errorf("expected defaultTimeout, got %v", p.timeout)
	}
}

func TestLDAPProbe_DurationRecorded(t *testing.T) {
	addr := startFakeLDAP(t, []byte{0x30, 0x00})
	p := NewLDAPProbe(addr, time.Second)
	res := p.Probe(context.Background())
	if !res.Healthy {
		t.Fatalf("expected healthy: %s", res.Message)
	}
	if res.Duration <= 0 {
		t.Errorf("expected positive duration, got %v", res.Duration)
	}
}

func TestLDAPProbe_Unhealthy_EmptyBanner(t *testing.T) {
	// A server that closes the connection immediately without sending data
	// should be treated as unhealthy.
	addr := startFakeLDAP(t, []byte{})
	p := NewLDAPProbe(addr, time.Second)
	res := p.Probe(context.Background())
	if res.Healthy {
		t.Fatal("expected unhealthy for empty banner (connection closed immediately)")
	}
}
