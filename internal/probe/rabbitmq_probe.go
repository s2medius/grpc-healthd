package probe

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/yourorg/grpc-healthd/internal/metrics"
)

// RabbitMQProbe checks RabbitMQ availability by connecting and verifying the AMQP banner.
type RabbitMQProbe struct {
	address string
	timeout time.Duration
}

// NewRabbitMQProbe creates a new RabbitMQProbe.
// If timeout is zero, DefaultTimeout is used.
func NewRabbitMQProbe(address string, timeout time.Duration) *RabbitMQProbe {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	return &RabbitMQProbe{address: address, timeout: timeout}
}

// Check dials the RabbitMQ server and verifies the AMQP protocol banner.
func (p *RabbitMQProbe) Check(ctx context.Context) Result {
	start := time.Now()

	dialer := &net.Dialer{Timeout: p.timeout}
	conn, err := dialer.DialContext(ctx, "tcp", p.address)
	if err != nil {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("rabbitmq", p.address, StatusUnhealthy.String(), duration)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("rabbitmq: dial failed: %w", err)}
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(p.timeout))

	reader := bufio.NewReader(conn)
	banner, err := reader.ReadString('\n')
	if err != nil && banner == "" {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("rabbitmq", p.address, StatusUnhealthy.String(), duration)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("rabbitmq: failed to read banner: %w", err)}
	}

	// AMQP 0-9-1 servers respond with "AMQP" in the initial banner
	if !strings.Contains(banner, "AMQP") {
		duration := time.Since(start).Seconds()
		metrics.RecordProbe("rabbitmq", p.address, StatusUnhealthy.String(), duration)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("rabbitmq: unexpected banner: %q", banner)}
	}

	duration := time.Since(start).Seconds()
	metrics.RecordProbe("rabbitmq", p.address, StatusHealthy.String(), duration)
	return Result{Status: StatusHealthy}
}
