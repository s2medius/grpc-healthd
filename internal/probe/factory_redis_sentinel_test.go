package probe

import (
	"testing"
	"time"

	"github.com/yourusername/grpc-healthd/internal/config"
)

func TestFromConfig_RedisSentinel(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "sentinel-check",
		Type:    "redis_sentinel",
		Address: "127.0.0.1:26379",
	}
	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
	sentinel, ok := p.(*RedisSentinelProbe)
	if !ok {
		t.Fatalf("expected *RedisSentinelProbe, got %T", p)
	}
	if sentinel.timeout != DefaultTimeout {
		t.Errorf("expected default timeout, got %v", sentinel.timeout)
	}
}

func TestFromConfig_RedisSentinel_WithTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "sentinel-check",
		Type:    "redis_sentinel",
		Address: "127.0.0.1:26379",
		Timeout: 4 * time.Second,
	}
	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sentinel, ok := p.(*RedisSentinelProbe)
	if !ok {
		t.Fatalf("expected *RedisSentinelProbe, got %T", p)
	}
	if sentinel.timeout != 4*time.Second {
		t.Errorf("expected 4s timeout, got %v", sentinel.timeout)
	}
}

func TestFromConfig_RedisSentinel_AddressPreserved(t *testing.T) {
	address := "127.0.0.1:26379"
	cfg := config.ProbeConfig{
		Name:    "sentinel-check",
		Type:    "redis_sentinel",
		Address: address,
	}
	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sentinel, ok := p.(*RedisSentinelProbe)
	if !ok {
		t.Fatalf("expected *RedisSentinelProbe, got %T", p)
	}
	if sentinel.address != address {
		t.Errorf("expected address %q, got %q", address, sentinel.address)
	}
}
