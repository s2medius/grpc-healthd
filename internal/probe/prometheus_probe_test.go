package probe

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestPrometheusProbe_Healthy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("# HELP go_goroutines Number of goroutines\ngo_goroutines 42\n"))
	}))
	defer ts.Close()

	p := NewPrometheusProbe(ts.URL, "go_goroutines", 2*time.Second)
	res := p.Probe(context.Background())

	if !res.Healthy {
		t.Fatalf("expected healthy, got unhealthy: %v", res.Err)
	}
}

func TestPrometheusProbe_Unhealthy_MetricMissing(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("# HELP some_other_metric\nsome_other_metric 1\n"))
	}))
	defer ts.Close()

	p := NewPrometheusProbe(ts.URL, "go_goroutines", 2*time.Second)
	res := p.Probe(context.Background())

	if res.Healthy {
		t.Fatal("expected unhealthy when metric is missing")
	}
}

func TestPrometheusProbe_Unhealthy_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	p := NewPrometheusProbe(ts.URL, "", 2*time.Second)
	res := p.Probe(context.Background())

	if res.Healthy {
		t.Fatal("expected unhealthy on 500 response")
	}
}

func TestPrometheusProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewPrometheusProbe("http://127.0.0.1:19999/metrics", "up", 300*time.Millisecond)
	res := p.Probe(context.Background())

	if res.Healthy {
		t.Fatal("expected unhealthy when connection refused")
	}
}

func TestNewPrometheusProbe_DefaultTimeout(t *testing.T) {
	p := NewPrometheusProbe("http://localhost:9090/metrics", "up", 0)
	if p.timeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %v", p.timeout)
	}
}

func TestPrometheusProbe_NoMetricFilter(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("anything\n"))
	}))
	defer ts.Close()

	p := NewPrometheusProbe(ts.URL, "", 2*time.Second)
	res := p.Probe(context.Background())

	if !res.Healthy {
		t.Fatalf("expected healthy with no metric filter, got: %v", res.Err)
	}
}
