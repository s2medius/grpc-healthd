package probe_test

import (
	"testing"
	"time"

	"github.com/your-org/grpc-healthd/internal/config"
	"github.com/your-org/grpc-healthd/internal/probe"
)

func TestFromConfig_Valkey(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "valkey-check",
		Type:    "valkey",
		Address: "localhost:6379",
	}

	p, err := probe.FromConfig(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p == nil {
		t.Fatal("expected probe, got nil")
	}
}

func TestFromConfig_Valkey_WithTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "valkey-check",
		Type:    "valkey",
		Address: "localhost:6379",
		Timeout: 3 * time.Second,
	}

	p, err := probe.FromConfig(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if p == nil {
		t.Fatal("expected probe, got nil")
	}
}
