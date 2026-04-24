package probe_test

import (
	"testing"

	"github.com/yourorg/grpc-healthd/internal/config"
	"github.com/yourorg/grpc-healthd/internal/probe"
)

func TestFromConfig_WebSocket(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "ws-service",
		Type:    "websocket",
		Address: "ws://localhost:8080/ws",
	}

	p, err := probe.FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestFromConfig_WebSocket_WithTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "ws-service",
		Type:    "websocket",
		Address: "ws://localhost:8080/health",
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
