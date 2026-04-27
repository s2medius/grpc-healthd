package probe

import (
	"testing"
	"time"

	"github.com/your-org/grpc-healthd/internal/config"
)

func TestFromConfig_Prometheus(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "prom-test",
		Type:    "prometheus",
		Address: "http://localhost:9090/metrics",
	}

	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil prober")
	}
}

func TestFromConfig_Prometheus_WithTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "prom-timeout",
		Type:    "prometheus",
		Address: "http://localhost:9090/metrics",
		Timeout: "3s",
	}

	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	pp, ok := p.(*PrometheusProbe)
	if !ok {
		t.Fatal("expected *PrometheusProbe")
	}
	if pp.timeout != 3*time.Second {
		t.Errorf("expected 3s timeout, got %v", pp.timeout)
	}
}

func TestFromConfig_Prometheus_WithMetricName(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "prom-metric",
		Type:    "prometheus",
		Address: "http://localhost:9090/metrics",
		Options: map[string]string{"metric_name": "up"},
	}

	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	pp, ok := p.(*PrometheusProbe)
	if !ok {
		t.Fatal("expected *PrometheusProbe")
	}
	if pp.metricName != "up" {
		t.Errorf("expected metric_name 'up', got %q", pp.metricName)
	}
}

func TestFromConfig_Prometheus_InvalidTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "prom-bad",
		Type:    "prometheus",
		Address: "http://localhost:9090/metrics",
		Timeout: "not-a-duration",
	}

	_, err := FromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for invalid timeout")
	}
}
