package probe

import (
	"testing"
	"time"

	"github.com/your-org/grpc-healthd/internal/config"
)

func TestFromConfig_HAProxy(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "haproxy-stats",
		Type:    "haproxy",
		Address: "127.0.0.1:1936",
	}
	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
	hp, ok := p.(*HAProxyProbe)
	if !ok {
		t.Fatalf("expected *HAProxyProbe, got %T", p)
	}
	if hp.address != "127.0.0.1:1936" {
		t.Errorf("expected address 127.0.0.1:1936, got %s", hp.address)
	}
	if hp.timeout != defaultTimeout {
		t.Errorf("expected default timeout, got %v", hp.timeout)
	}
}

func TestFromConfig_HAProxy_WithTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "haproxy-custom",
		Type:    "haproxy",
		Address: "127.0.0.1:1936",
		Timeout: 3 * time.Second,
	}
	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	hp, ok := p.(*HAProxyProbe)
	if !ok {
		t.Fatalf("expected *HAProxyProbe, got %T", p)
	}
	if hp.timeout != 3*time.Second {
		t.Errorf("expected 3s timeout, got %v", hp.timeout)
	}
}
