package server_test

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/grpc-healthd/internal/health"
	"github.com/grpc-healthd/internal/probe"
	"github.com/grpc-healthd/internal/server"
)

func newTestChecker(t *testing.T) *health.Checker {
	t.Helper()
	return health.NewChecker()
}

func TestCheck_ServiceUnknown(t *testing.T) {
	checker := newTestChecker(t)
	hs := server.NewHealthServer(checker)
	resp, err := hs.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: "missing"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status != grpc_health_v1.HealthCheckResponse_SERVICE_UNKNOWN {
		t.Errorf("expected SERVICE_UNKNOWN, got %v", resp.Status)
	}
}

func TestCheck_Serving(t *testing.T) {
	checker := newTestChecker(t)
	p := probe.NewTCPProbe("localhost:1", 100*time.Millisecond)
	checker.Register("svc", p, 50*time.Millisecond)
	time.Sleep(200 * time.Millisecond)

	hs := server.NewHealthServer(checker)
	resp, err := hs.Check(context.Background(), &grpc_health_v1.HealthCheckRequest{Service: "svc"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Status == grpc_health_v1.HealthCheckResponse_SERVICE_UNKNOWN {
		t.Error("expected known status")
	}
}

func TestNewHealthServer_NotNil(t *testing.T) {
	checker := newTestChecker(t)
	hs := server.NewHealthServer(checker)
	if hs == nil {
		t.Error("expected non-nil HealthServer")
	}
}
