package probe

import (
	"testing"
	"time"

	"github.com/yourorg/grpc-healthd/internal/config"
)

func TestFromConfig_Consul(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "consul-local",
		Type:    "consul",
		Address: "localhost:8500",
	}

	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
	cp, ok := p.(*ConsulProbe)
	if !ok {
		t.Fatalf("expected *ConsulProbe, got %T", p)
	}
	if cp.address != "localhost:8500" {
		t.Fatalf("expected address localhost:8500, got %s", cp.address)
	}
	if cp.timeout != DefaultTimeout {
		t.Fatalf("expected default timeout, got %v", cp.timeout)
	}
}

func TestFromConfig_Consul_WithTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "consul-custom",
		Type:    "consul",
		Address: "consul.internal:8500",
		Timeout: 3 * time.Second,
	}

	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cp, ok := p.(*ConsulProbe)
	if !ok {
		t.Fatalf("expected *ConsulProbe, got %T", p)
	}
	if cp.timeout != 3*time.Second {
		t.Fatalf("expected 3s timeout, got %v", cp.timeout)
	}
}
