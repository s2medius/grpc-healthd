package probe_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/patrickward/grpc-healthd/internal/probe"
)

func TestNeo4jProbe_Healthy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	address := strings.TrimPrefix(ts.URL, "http://")
	p := probe.NewNeo4jProbe(address, time.Second)
	res := p.Check(context.Background())

	if !res.Healthy {
		t.Fatalf("expected healthy, got error: %v", res.Error)
	}
}

func TestNeo4jProbe_Unhealthy_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	address := strings.TrimPrefix(ts.URL, "http://")
	p := probe.NewNeo4jProbe(address, time.Second)
	res := p.Check(context.Background())

	if res.Healthy {
		t.Fatal("expected unhealthy for 503 response")
	}
	if res.Error == nil {
		t.Fatal("expected non-nil error")
	}
}

func TestNeo4jProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := probe.NewNeo4jProbe("127.0.0.1:19999", time.Second)
	res := p.Check(context.Background())

	if res.Healthy {
		t.Fatal("expected unhealthy when connection is refused")
	}
}

func TestNewNeo4jProbe_DefaultTimeout(t *testing.T) {
	p := probe.NewNeo4jProbe("localhost:7474", 0)
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestNeo4jProbe_DurationRecorded(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	address := strings.TrimPrefix(ts.URL, "http://")
	p := probe.NewNeo4jProbe(address, time.Second)
	res := p.Check(context.Background())

	if !res.Healthy {
		t.Fatalf("expected healthy: %v", res.Error)
	}
}
