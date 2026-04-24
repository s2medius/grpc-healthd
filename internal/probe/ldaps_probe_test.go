package probe

import (
	"context"
	"crypto/tls"
	"net"
	"testing"
	"time"
)

func startFakeLDAPS(t *testing.T, banner []byte) string {
	t.Helper()
	cert := selfSignedCert(t)
	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	ln, err := tls.Listen("tcp", "127.0.0.1:0", cfg)
	if err != nil {
		t.Fatalf("failed to start fake LDAPS: %v", err)
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

// selfSignedCert reuses the helper from tls_probe_test.go via the same package.
func selfSignedCert(t *testing.T) tls.Certificate {
	t.Helper()
	// Use a plain TCP listener trick: generate inline via crypto/tls test helper.
	// For simplicity, borrow from the existing selfSignedTLSServer pattern.
	ln := selfSignedTLSServer(t, []byte{})
	_ = ln // just need the cert generation side-effect path
	// Re-generate a self-signed cert directly.
	cert, err := tls.X509KeyPair(testCertPEM, testKeyPEM)
	if err != nil {
		t.Fatalf("failed to load test cert: %v", err)
	}
	return cert
}

func TestLDAPSProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewLDAPSProbe("127.0.0.1:1", time.Second, true)
	res := p.Probe(context.Background())
	if res.Healthy {
		t.Fatal("expected unhealthy for refused connection")
	}
}

func TestNewLDAPSProbe_DefaultTimeout(t *testing.T) {
	p := NewLDAPSProbe("localhost:636", 0, false)
	if p.timeout != defaultTimeout {
		t.Errorf("expected defaultTimeout, got %v", p.timeout)
	}
}

func TestLDAPSProbe_SkipVerifyFlag(t *testing.T) {
	p := NewLDAPSProbe("localhost:636", 2*time.Second, true)
	if !p.skipVerify {
		t.Error("expected skipVerify to be true")
	}
}

func TestLDAPSProbe_Healthy(t *testing.T) {
	addr := startFakeLDAPS(t, []byte{0x30, 0x00})
	p := NewLDAPSProbe(addr, time.Second, true)
	res := p.Probe(context.Background())
	if !res.Healthy {
		t.Fatalf("expected healthy, got: %s", res.Message)
	}
}

func init() {
	// Ensure net package is imported for startFakeLDAPS.
	_ = net.Dial
}
