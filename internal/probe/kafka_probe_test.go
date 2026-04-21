package probe_test

import (
	"context"
	"encoding/binary"
	"net"
	"testing"
	"time"

	"github.com/andrewhowdencom/grpc-healthd/internal/probe"
)

// startFakeKafka starts a minimal TCP server that mimics a Kafka broker
// by responding to any connection with a 4-byte length-prefixed response.
func startFakeKafka(t *testing.T, healthy bool) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start fake kafka: %v", err)
	}
	t.Cleanup(func() { ln.Close() })

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		if !healthy {
			return // close without responding
		}

		// Drain the incoming request
		buf := make([]byte, 128)
		_, _ = conn.Read(buf)

		// Write a minimal valid response: 4-byte length + 4-byte correlationId + 2-byte error
		body := []byte{0x00, 0x00, 0x00, 0x01, 0x00, 0x00}
		respLen := make([]byte, 4)
		binary.BigEndian.PutUint32(respLen, uint32(len(body)))
		_, _ = conn.Write(respLen)
		_, _ = conn.Write(body)
	}()

	return ln.Addr().String()
}

func TestKafkaProbe_Healthy(t *testing.T) {
	addr := startFakeKafka(t, true)
	p := probe.NewKafkaProbe(addr, time.Second)
	res := p.Check(context.Background())
	if res.Status != probe.StatusHealthy {
		t.Errorf("expected healthy, got %s: %v", res.Status, res.Err)
	}
}

func TestKafkaProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := probe.NewKafkaProbe("127.0.0.1:19999", time.Second)
	res := p.Check(context.Background())
	if res.Status != probe.StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", res.Status)
	}
}

func TestKafkaProbe_Unhealthy_NoResponse(t *testing.T) {
	addr := startFakeKafka(t, false)
	p := probe.NewKafkaProbe(addr, 500*time.Millisecond)
	res := p.Check(context.Background())
	if res.Status != probe.StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", res.Status)
	}
}

func TestNewKafkaProbe_DefaultTimeout(t *testing.T) {
	p := probe.NewKafkaProbe("localhost:9092", 0)
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestKafkaProbe_DurationRecorded(t *testing.T) {
	addr := startFakeKafka(t, true)
	p := probe.NewKafkaProbe(addr, time.Second)
	res := p.Check(context.Background())
	if res.Duration <= 0 {
		t.Errorf("expected positive duration, got %v", res.Duration)
	}
}
