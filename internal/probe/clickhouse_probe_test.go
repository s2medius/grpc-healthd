package probe_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/your-org/grpc-healthd/internal/probe"
)

func startFakeClickHouse(t *testing.T, statusCode int) string {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ping" {
			w.WriteHeader(statusCode)
			if statusCode == http.StatusOK {
				_, _ = w.Write([]byte("Ok.\n"))
			}
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(ts.Close)
	return ts.Listener.Addr().String()
}

func TestClickHouseProbe_Healthy(t *testing.T) {
	addr := startFakeClickHouse(t, http.StatusOK)
	p := probe.NewClickHouseProbe(addr, 2*time.Second)
	res := p.Probe(context.Background())
	if res.Status != probe.StatusHealthy {
		t.Fatalf("expected healthy, got %s: %v", res.Status, res.Error)
	}
}

func TestClickHouseProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	ln.Close()

	p := probe.NewClickHouseProbe(addr, 500*time.Millisecond)
	res := p.Probe(context.Background())
	if res.Status != probe.StatusUnhealthy {
		t.Fatalf("expected unhealthy, got %s", res.Status)
	}
}

func TestNewClickHouseProbe_DefaultTimeout(t *testing.T) {
	p := probe.NewClickHouseProbe("localhost:8123", 0)
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestClickHouseProbe_CustomTimeout(t *testing.T) {
	p := probe.NewClickHouseProbe("localhost:8123", 3*time.Second)
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestClickHouseProbe_Unhealthy_ServiceUnavailable(t *testing.T) {
	addr := startFakeClickHouse(t, http.StatusServiceUnavailable)
	p := probe.NewClickHouseProbe(addr, 2*time.Second)
	res := p.Probe(context.Background())
	if res.Status != probe.StatusUnhealthy {
		t.Fatalf("expected unhealthy, got %s", res.Status)
	}
	if res.Error == nil {
		t.Fatal("expected non-nil error")
	}
	expected := fmt.Sprintf("unexpected status %d", http.StatusServiceUnavailable)
	if res.Error.Error() != expected {
		t.Fatalf("expected error %q, got %q", expected, res.Error.Error())
	}
}
