package probe

import (
	"testing"
	"time"

	"github.com/your-org/grpc-healthd/internal/config"
)

func TestFromConfig_ClickHouse(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "ch-test",
		Type:    "clickhouse",
		Address: "localhost:9000",
	}
	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
	cp, ok := p.(*ClickHouseProbe)
	if !ok {
		t.Fatalf("expected *ClickHouseProbe, got %T", p)
	}
	if cp.address != "localhost:9000" {
		t.Fatalf("expected address localhost:9000, got %s", cp.address)
	}
	if cp.timeout != DefaultTimeout {
		t.Fatalf("expected default timeout, got %v", cp.timeout)
	}
}

func TestFromConfig_ClickHouse_WithTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "ch-timeout",
		Type:    "clickhouse",
		Address: "localhost:9000",
		Timeout: "4s",
	}
	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cp := p.(*ClickHouseProbe)
	if cp.timeout != 4*time.Second {
		t.Fatalf("expected 4s timeout, got %v", cp.timeout)
	}
}
