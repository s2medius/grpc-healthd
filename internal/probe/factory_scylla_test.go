package probe_test

import (
	"testing"

	"github.com/your-org/grpc-healthd/internal/config"
	"github.com/your-org/grpc-healthd/internal/probe"
)

func TestFromConfig_Scylla(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "scylla-check",
		Type:    "scylla",
		Address: "localhost:9042",
	}
	p, err := probe.FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestFromConfig_Scylla_WithTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "scylla-check",
		Type:    "scylla",
		Address: "localhost:9042",
		Timeout: "3s",
	}
	p, err := probe.FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestFromConfig_Scylla_InvalidTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "scylla-check",
		Type:    "scylla",
		Address: "localhost:9042",
		Timeout: "not-a-duration",
	}
	_, err := probe.FromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for invalid timeout, got nil")
	}
}
