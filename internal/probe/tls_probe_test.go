package probe_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/your-org/grpc-healthd/internal/probe"
)

func selfSignedTLSServer(t *testing.T) (addr string, shutdown func()) {
	t.Helper()
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "localhost"},
		DNSNames:     []string{"localhost"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
	}
	certDER, _ := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	tlsCert := tls.Certificate{Certificate: [][]byte{certDER}, PrivateKey: key}
	ln, err := tls.Listen("tcp", "127.0.0.1:0", &tls.Config{Certificates: []tls.Certificate{tlsCert}})
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return ln.Addr().String(), func() { ln.Close() }
}

func TestTLSProbe_Healthy(t *testing.T) {
	addr, shutdown := selfSignedTLSServer(t)
	defer shutdown()
	p := probe.NewTLSProbe(addr, 3*time.Second, true)
	res := p.Probe(context.Background())
	if res.Status != probe.StatusHealthy {
		t.Fatalf("expected healthy, got %v: %v", res.Status, res.Error)
	}
}

func TestTLSProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := probe.NewTLSProbe("127.0.0.1:19999", time.Second, true)
	res := p.Probe(context.Background())
	if res.Status != probe.StatusUnhealthy {
		t.Fatal("expected unhealthy")
	}
}

func TestNewTLSProbe_DefaultTimeout(t *testing.T) {
	p := probe.NewTLSProbe("localhost:443", 0, false)
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestTLSProbe_DurationRecorded(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	ln.Close()
	p := probe.NewTLSProbe(ln.Addr().String(), time.Second, true)
	res := p.Probe(context.Background())
	if res.Duration <= 0 {
		t.Fatal("expected positive duration")
	}
}
