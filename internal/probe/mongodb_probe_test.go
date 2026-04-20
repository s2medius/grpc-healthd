package probe

import (
	"context"
	"encoding/binary"
	"net"
	"testing"
	"time"
)

// startFakeMongoDB starts a minimal TCP server that responds with a valid MongoDB
// wire-protocol header so the probe considers it healthy.
func startFakeMongoDB(t *testing.T) string {
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
		// Drain the incoming handshake.
		buf := make([]byte, 256)
		_, _ = conn.Read(buf)
		// Write a minimal 4-byte response length (36 bytes total).
		resp := make([]byte, 4)
		binary.LittleEndian.PutUint32(resp, 36)
		_, _ = conn.Write(resp)
	}()

	return ln.Addr().String()
}

func TestMongoDBProbe_Healthy(t *testing.T) {
	addr := startFakeMongoDB(t)
	p := NewMongoDBProbe(addr, 2*time.Second)
	res := p.Probe(context.Background())
	if res.Status != StatusHealthy {
		t.Errorf("expected healthy, got %v (err: %v)", res.Status, res.Err)
	}
}

func TestMongoDBProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewMongoDBProbe("127.0.0.1:1", 500*time.Millisecond)
	res := p.Probe(context.Background())
	if res.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %v", res.Status)
	}
}

func TestNewMongoDBProbe_DefaultTimeout(t *testing.T) {
	p := NewMongoDBProbe("localhost:27017", 0)
	if p.timeout != DefaultTimeout {
		t.Errorf("expected default timeout %v, got %v", DefaultTimeout, p.timeout)
	}
}

func TestMongoDBProbe_CustomTimeout(t *testing.T) {
	const want = 3 * time.Second
	p := NewMongoDBProbe("localhost:27017", want)
	if p.timeout != want {
		t.Errorf("expected %v, got %v", want, p.timeout)
	}
}

func TestMongoDBProbe_DurationRecorded(t *testing.T) {
	addr := startFakeMongoDB(t)
	p := NewMongoDBProbe(addr, 2*time.Second)
	res := p.Probe(context.Background())
	if res.Duration <= 0 {
		t.Errorf("expected positive duration, got %v", res.Duration)
	}
}
