package probe_test

import (
	"testing"
	"time"

	"github.com/yourorg/grpc-healthd/internal/config"
	"github.com/yourorg/grpc-healthd/internal/probe"
)

func TestFromConfig_SSH(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "ssh-check",
		Type:    "ssh",
		Address: "localhost:22",
	}
	p, err := probe.FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestFromConfig_SSH_WithTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "ssh-check",
		Type:    "ssh",
		Address: "localhost:22",
		Timeout: 5 * time.Second,
	}
	p, err := probe.FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}
