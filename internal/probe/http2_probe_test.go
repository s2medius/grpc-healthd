package probe_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourorg/grpc-healthd/internal/probe"
)

func TestHTTP2Probe_Healthy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	p := probe.NewHTTP2Probe(ts.URL, 2*time.Second)
	res := p.Probe(context.Background())

	if res.Status != probe.StatusHealthy {
		t.Errorf("expected healthy, got %s (err: %v)", res.Status, res.Err)
	}
}

func TestHTTP2Probe_Unhealthy_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	p := probe.NewHTTP2Probe(ts.URL, 2*time.Second)
	res := p.Probe(context.Background())

	if res.Status != probe.StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", res.Status)
	}
}

func TestHTTP2Probe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := probe.NewHTTP2Probe("http://127.0.0.1:19999", 500*time.Millisecond)
	res := p.Probe(context.Background())

	if res.Status != probe.StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", res.Status)
	}
	if res.Err == nil {
		t.Error("expected non-nil error")
	}
}

func TestNewHTTP2Probe_DefaultTimeout(t *testing.T) {
	p := probe.NewHTTP2Probe("http://localhost", 0)
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestHTTP2Probe_DurationRecorded(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	p := probe.NewHTTP2Probe(ts.URL, 2*time.Second)
	res := p.Probe(context.Background())

	if res.Duration < 10*time.Millisecond {
		t.Errorf("expected duration >= 10ms, got %s", res.Duration)
	}
}
