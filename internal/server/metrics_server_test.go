package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"
)

func freePort(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not find free port: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	_ = l.Close()
	return fmt.Sprintf("127.0.0.1:%d", port)
}

func TestNewMetricsServer_NotNil(t *testing.T) {
	s := NewMetricsServer(":9090")
	if s == nil {
		t.Fatal("expected non-nil MetricsServer")
	}
}

func TestMetricsServer_Addr(t *testing.T) {
	s := NewMetricsServer(":9091")
	if s.Addr() != ":9091" {
		t.Errorf("expected :9091, got %s", s.Addr())
	}
}

func TestMetricsServer_HealthzEndpoint(t *testing.T) {
	addr := freePort(t)
	s := NewMetricsServer(addr)

	go func() { _ = s.ListenAndServe() }()
	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get(fmt.Sprintf("http://%s/healthz", addr))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = s.Shutdown(ctx)
}

func TestMetricsServer_MetricsEndpoint(t *testing.T) {
	addr := freePort(t)
	s := NewMetricsServer(addr)

	go func() { _ = s.ListenAndServe() }()
	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get(fmt.Sprintf("http://%s/metrics", addr))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_ = s.Shutdown(ctx)
}
