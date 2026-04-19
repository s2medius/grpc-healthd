package probe

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// GRPCProbe checks health via gRPC Health Checking Protocol.
type GRPCProbe struct {
	address string
	service string
	timeout time.Duration
}

// NewGRPCProbe creates a new GRPCProbe. Timeout defaults to 5s if zero.
func NewGRPCProbe(address, service string, timeout time.Duration) *GRPCProbe {
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	return &GRPCProbe{address: address, service: service, timeout: timeout}
}

// Probe performs the gRPC health check and returns a Result.
func (p *GRPCProbe) Probe(ctx context.Context) Result {
	start := time.Now()

	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, p.address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return Result{Status: StatusUnhealthy, Duration: time.Since(start),
			Message: fmt.Sprintf("dial failed: %v", err)}
	}
	defer conn.Close()

	client := grpc_health_v1.NewHealthClient(conn)
	resp, err := client.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: p.service})
	duration := time.Since(start)
	if err != nil {
		return Result{Status: StatusUnhealthy, Duration: duration,
			Message: fmt.Sprintf("health check rpc failed: %v", err)}
	}
	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		return Result{Status: StatusUnhealthy, Duration: duration,
			Message: fmt.Sprintf("service status: %s", resp.Status)}
	}
	return Result{Status: StatusHealthy, Duration: duration, Message: "SERVING"}
}
