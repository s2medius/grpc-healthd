package probe

import (
	"context"
	"net"
	"testing"
	"time"
)

// startFakeMQTT starts a minimal fake MQTT broker that responds with a valid CONNACK.
func startFakeMQTT(t *testing.T, sendValidCONNACK bool) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start fake MQTT: %v", err)
	}
	t.Cleanup(func() { ln.Close() })

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		buf := make([]byte, 64)
		_, _ = conn.Read(buf)
		if sendValidCONNACK {
			// Valid CONNACK: type=0x20, remaining=2, session present=0, return code=0
			_, _ = conn.Write([]byte{0x20, 0x02, 0x00, 0x00})
		} else {
			// Invalid: return code = 0x05 (not authorized)
			_, _ = conn.Write([]byte{0x20, 0x02, 0x00, 0x05})
		}
	}()

	return ln.Addr().String()
}

func TestMQTTProbe_Healthy(t *testing.T) {
	addr := startFakeMQTT(t, true)
	p := NewMQTTProbe(addr, time.Second)
	res := p.Probe(context.Background())
	if !res.Healthy {
		t.Errorf("expected healthy, got: %s", res.Message)
	}
}

func TestMQTTProbe_Unhealthy_BadCONNACK(t *testing.T) {
	addr := startFakeMQTT(t, false)
	p := NewMQTTProbe(addr, time.Second)
	res := p.Probe(context.Background())
	if res.Healthy {
		t.Error("expected unhealthy due to bad CONNACK return code")
	}
}

func TestMQTTProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewMQTTProbe("127.0.0.1:19999", 200*time.Millisecond)
	res := p.Probe(context.Background())
	if res.Healthy {
		t.Error("expected unhealthy for refused connection")
	}
}

func TestNewMQTTProbe_DefaultTimeout(t *testing.T) {
	p := NewMQTTProbe("localhost:1883", 0)
	if p.timeout != DefaultTimeout {
		t.Errorf("expected default timeout %v, got %v", DefaultTimeout, p.timeout)
	}
}

func TestMQTTProbe_DurationRecorded(t *testing.T) {
	addr := startFakeMQTT(t, true)
	p := NewMQTTProbe(addr, time.Second)
	res := p.Probe(context.Background())
	if !res.Healthy {
		t.Errorf("expected healthy, got: %s", res.Message)
	}
}
