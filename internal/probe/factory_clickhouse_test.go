package probe_test

import (
	"testing"
	"time"

	"github.com/your-org/grpc-healthd/internal/config"
	"github.com/your-org/grpc-healthd/internal/probe"
)

func TestFromConfig_ClickHouse(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "clickhouse-test",
		Type:    "clickhouse",
		Address: "localhost:8123",
	}
	p, err := probe.FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestFromConfig_ClickHouse_WithTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "clickhouse-timeout",
		Type:    "clickhouse",
		Address: "localhost:8123",
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
