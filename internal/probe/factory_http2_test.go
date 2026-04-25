package probe_test

import (
	"testing"

	"github.com/yourorg/grpc-healthd/internal/config"
	"github.com/yourorg/grpc-healthd/internal/probe"
)

func TestFromConfig_HTTP2(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "my-http2-service",
		Type:    "http2",
		Address: "http://localhost:8080/health",
	}

	p, err := probe.FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestFromConfig_HTTP2_WithTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "my-http2-service",
		Type:    "http2",
		Address: "http://localhost:8080/health",
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

func TestFromConfig_HTTP2_InvalidTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "my-http2-service",
		Type:    "http2",
		Address: "http://localhost:8080/health",
		Timeout: "not-a-duration",
	}

	_, err := probe.FromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for invalid timeout, got nil")
	}
}
