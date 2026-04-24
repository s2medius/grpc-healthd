package probe

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/grpc-healthd/internal/metrics"
)

// MQTTProbe checks health of an MQTT broker by performing a TCP handshake
// and verifying the server sends a valid MQTT CONNACK response to a CONNECT packet.
type MQTTProbe struct {
	address string
	timeout time.Duration
}

// NewMQTTProbe creates a new MQTTProbe for the given address.
// If timeout is zero, DefaultTimeout is used.
func NewMQTTProbe(address string, timeout time.Duration) *MQTTProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &MQTTProbe{address: address, timeout: timeout}
}

// mqttConnectPacket returns a minimal MQTT 3.1.1 CONNECT packet.
func mqttConnectPacket() []byte {
	return []byte{
		0x10,       // CONNECT packet type
		0x0c,       // Remaining length = 12
		0x00, 0x04, // Protocol name length
		'M', 'Q', 'T', 'T', // Protocol name
		0x04,       // Protocol level (3.1.1)
		0x00,       // Connect flags (clean session off)
		0x00, 0x3c, // Keep-alive 60s
		0x00, 0x00, // Client ID length = 0
	}
}

func (p *MQTTProbe) Probe(ctx context.Context) Result {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", p.address, p.timeout)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("mqtt", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("connection failed: %v", err)}
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(p.timeout))

	if _, err := conn.Write(mqttConnectPacket()); err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("mqtt", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("write failed: %v", err)}
	}

	// Read CONNACK: fixed header (2 bytes) + variable header (2 bytes) = 4 bytes
	buf := make([]byte, 4)
	if _, err := conn.Read(buf); err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("mqtt", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("read failed: %v", err)}
	}

	// CONNACK packet type is 0x20, return code 0x00 means success
	if buf[0] != 0x20 || buf[3] != 0x00 {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("mqtt", p.address, false, duration)
		return Result{Healthy: false, Message: fmt.Sprintf("unexpected CONNACK: %x", buf)}
	}

	duration := time.Since(start).Seconds()
	metrics.RecordProbe("mqtt", p.address, true, duration)
	return Result{Healthy: true, Message: "MQTT broker reachable"}
}
