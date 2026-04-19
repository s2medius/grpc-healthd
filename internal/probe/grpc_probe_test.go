package probe_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/yourorg/grpc-healthd/internal/probe"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func startHealthServer(t *testing.T, serving bool) string {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	healthSrv := health.NewServer()
	status := grpc_health_v1.HealthCheckResponse_SERVING
	if !serving {
		status = grpc_health_v1.HealthCheckResponse_NOT_SERVING
	}
	healthSrv.SetServingStatus("", status)
	grpc_health_v1.RegisterHealthServer(srv, healthSrv)
	go srv.Serve(lis) //nolint:errcheck
	t.Cleanup(srv.Stop)
	return lis.Addr().String()
}

func TestGRPCProbe_Healthy(t *testing.T) {
	addr := startHealthServer(t, true)
	p := probe.NewGRPCProbe(addr, "", 3*time.Second)
	result := p.Probe(context.Background())
	if result.Status != probe.StatusHealthy {
		t.Errorf("expected healthy, got %s: %s", result.Status, result.Message)
	}
}

func TestGRPCProbe_NotServing(t *testing.T) {
	addr := startHealthServer(t, false)
	p := probe.NewGRPCProbe(addr, "", 3*time.Second)
	result := p.Probe(context.Background())
	if result.Status != probe.StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", result.Status)
	}
}

func TestGRPCProbe_ConnectionRefused(t *testing.T) {
	p := probe.NewGRPCProbe("127.0.0.1:1", "", 500*time.Millisecond)
	result := p.Probe(context.Background())
	if result.Status != probe.StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", result.Status)
	}
}

func TestNewGRPCProbe_DefaultTimeout(t *testing.T) {
	p := probe.NewGRPCProbe("localhost:50051", "svc", 0)
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestGRPCProbe_DurationRecorded(t *testing.T) {
	addr := startHealthServer(t, true)
	p := probe.NewGRPCProbe(addr, "", 3*time.Second)
	result := p.Probe(context.Background())
	if result.Duration <= 0 {
		t.Errorf("expected positive duration, got %v", result.Duration)
	}
}
