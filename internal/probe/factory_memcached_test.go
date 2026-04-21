package probe

import (
	"testing"

	"github.com/yourorg/grpc-healthd/internal/config"
)

func TestFromConfig_Memcached(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "memcached-check",
		Type:    "memcached",
		Address: "localhost:11211",
	}
	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
	mp, ok := p.(*MemcachedProbe)
	if !ok {
		t.Fatalf("expected *MemcachedProbe, got %T", p)
	}
	if mp.address != "localhost:11211" {
		t.Errorf("expected address %q, got %q", "localhost:11211", mp.address)
	}
	if mp.timeout != DefaultTimeout {
		t.Errorf("expected default timeout %v, got %v", DefaultTimeout, mp.timeout)
	}
}

func TestFromConfig_Memcached_WithTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "memcached-check",
		Type:    "memcached",
		Address: "localhost:11211",
		Timeout: "2s",
	}
	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	mp, ok := p.(*MemcachedProbe)
	if !ok {
		t.Fatalf("expected *MemcachedProbe, got %T", p)
	}
	if mp.timeout.Seconds() != 2 {
		t.Errorf("expected 2s timeout, got %v", mp.timeout)
	}
}
