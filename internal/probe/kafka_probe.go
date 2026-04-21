package probe

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/andrewhowdencom/grpc-healthd/internal/metrics"
)

// KafkaProbe checks health of a Kafka broker by establishing a TCP connection
// and performing a minimal Kafka API versions request handshake.
type KafkaProbe struct {
	address string
	timeout time.Duration
}

// NewKafkaProbe creates a new KafkaProbe targeting the given address.
// If timeout is zero, DefaultTimeout is used.
func NewKafkaProbe(address string, timeout time.Duration) *KafkaProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &KafkaProbe{
		address: address,
		timeout: timeout,
	}
}

// Check connects to the Kafka broker and verifies it accepts connections.
// It sends a minimal ApiVersions request (Kafka protocol) and expects a valid response.
func (p *KafkaProbe) Check(ctx context.Context) Result {
	start := time.Now()

	deadline := start.Add(p.timeout)
	conn, err := net.DialTimeout("tcp", p.address, time.Until(deadline))
	if err != nil {
		dur := time.Since(start)
		metrics.RecordProbe("kafka", p.address, StatusUnhealthy.String(), dur)
		return Result{
			Status:   StatusUnhealthy,
			Duration: dur,
			Err:      fmt.Errorf("kafka connect: %w", err),
		}
	}
	defer conn.Close()

	// ApiVersions request v0: length(4) + apiKey(2) + apiVersion(2) + correlationId(4) + clientId(2)
	request := []byte{
		0x00, 0x00, 0x00, 0x0a, // message length = 10
		0x00, 0x12, // ApiKey: ApiVersions = 18
		0x00, 0x00, // ApiVersion: 0
		0x00, 0x00, 0x00, 0x01, // CorrelationId: 1
		0xff, 0xff, // ClientId: null (length -1)
	}

	_ = conn.SetDeadline(deadline)
	if _, err := conn.Write(request); err != nil {
		dur := time.Since(start)
		metrics.RecordProbe("kafka", p.address, StatusUnhealthy.String(), dur)
		return Result{
			Status:   StatusUnhealthy,
			Duration: dur,
			Err:      fmt.Errorf("kafka write: %w", err),
		}
	}

	// Read at least the 4-byte response length header
	header := make([]byte, 4)
	if _, err := conn.Read(header); err != nil {
		dur := time.Since(start)
		metrics.RecordProbe("kafka", p.address, StatusUnhealthy.String(), dur)
		return Result{
			Status:   StatusUnhealthy,
			Duration: dur,
			Err:      fmt.Errorf("kafka read: %w", err),
		}
	}

	dur := time.Since(start)
	metrics.RecordProbe("kafka", p.address, StatusHealthy.String(), dur)
	return Result{
		Status:   StatusHealthy,
		Duration: dur,
	}
}
