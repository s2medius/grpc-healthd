package probe

import (
	"testing"
	"time"

	"github.com/your-org/grpc-healthd/internal/config"
)

func TestFromConfig_HTTPJson(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:      "api-health",
		Type:      "http_json",
		Address:   "http://localhost:8080/health",
		JSONKey:   "status",
		JSONValue: "ok",
	}

	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil probe")
	}

	probe, ok := p.(*HTTPJSONProbe)
	if !ok {
		t.Fatalf("expected *HTTPJSONProbe, got %T", p)
	}
	if probe.address != cfg.Address {
		t.Errorf("address: got %s, want %s", probe.address, cfg.Address)
	}
	if probe.jsonKey != cfg.JSONKey {
		t.Errorf("jsonKey: got %s, want %s", probe.jsonKey, cfg.JSONKey)
	}
	if probe.jsonValue != cfg.JSONValue {
		t.Errorf("jsonValue: got %s, want %s", probe.jsonValue, cfg.JSONValue)
	}
}

func TestFromConfig_HTTPJson_WithTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:      "api-health",
		Type:      "http_json",
		Address:   "http://localhost:8080/health",
		JSONKey:   "status",
		JSONValue: "ok",
		Timeout:   10 * time.Second,
	}

	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	probe, ok := p.(*HTTPJSONProbe)
	if !ok {
		t.Fatalf("expected *HTTPJSONProbe, got %T", p)
	}
	if probe.timeout != 10*time.Second {
		t.Errorf("timeout: got %v, want 10s", probe.timeout)
	}
}
