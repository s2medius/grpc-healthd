package probe

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/your-org/grpc-healthd/internal/config"
)

func TestPrometheusProbe_EndToEnd_ViaFactory(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("# HELP up Is the service up\nup 1\n"))
	}))
	defer ts.Close()

	cfg := config.ProbeConfig{
		Name:    "e2e-prometheus",
		Type:    "prometheus",
		Address: ts.URL,
		Timeout: "2s",
		Options: map[string]string{"metric_name": "up"},
	}

	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("factory error: %v", err)
	}

	res := p.Probe(context.Background())
	if !res.Healthy {
		t.Fatalf("expected healthy result, got: %v", res.Err)
	}
	if res.Duration <= 0 {
		t.Error("expected positive duration")
	}
}

func TestPrometheusProbe_EndToEnd_MetricAbsent(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("# HELP other_metric\nother_metric 0\n"))
	}))
	defer ts.Close()

	cfg := config.ProbeConfig{
		Name:    "e2e-prometheus-absent",
		Type:    "prometheus",
		Address: ts.URL,
		Options: map[string]string{"metric_name": "up"},
	}

	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("factory error: %v", err)
	}

	res := p.Probe(context.Background())
	if res.Healthy {
		t.Fatal("expected unhealthy when required metric is absent")
	}
}
