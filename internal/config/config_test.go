package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/yourorg/grpc-healthd/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "config-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()
	if cfg.Server.GRPCListenAddr != ":50051" {
		t.Errorf("expected :50051, got %s", cfg.Server.GRPCListenAddr)
	}
	if cfg.Server.MetricsListenAddr != ":9090" {
		t.Errorf("expected :9090, got %s", cfg.Server.MetricsListenAddr)
	}
}

func TestLoad_ValidConfig(t *testing.T) {
	yaml := `
server:
  grpc_listen_addr: ":50052"
  metrics_listen_addr: ":9091"
probes:
  - name: my-service
    address: localhost:8080
    service: myapp.MyService
    timeout: 3s
    interval: 10s
`
	path := writeTempConfig(t, yaml)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Server.GRPCListenAddr != ":50052" {
		t.Errorf("unexpected grpc addr: %s", cfg.Server.GRPCListenAddr)
	}
	if len(cfg.Probes) != 1 {
		t.Fatalf("expected 1 probe, got %d", len(cfg.Probes))
	}
	if cfg.Probes[0].Timeout != 3*time.Second {
		t.Errorf("unexpected timeout: %v", cfg.Probes[0].Timeout)
	}
}

func TestLoad_DefaultsApplied(t *testing.T) {
	yaml := `
probes:
  - name: bare
    address: localhost:9000
`
	path := writeTempConfig(t, yaml)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Probes[0].Timeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %v", cfg.Probes[0].Timeout)
	}
	if cfg.Probes[0].Interval != 30*time.Second {
		t.Errorf("expected default interval 30s, got %v", cfg.Probes[0].Interval)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
