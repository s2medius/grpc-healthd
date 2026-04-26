package probe

import (
	"testing"
	"time"

	"github.com/your-org/grpc-healthd/internal/config"
)

func TestFromConfig_Dgraph(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "dgraph-test",
		Type:    "dgraph",
		Address: "localhost:8080",
	}
	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
	dp, ok := p.(*DgraphProbe)
	if !ok {
		t.Fatalf("expected *DgraphProbe, got %T", p)
	}
	if dp.address != "localhost:8080" {
		t.Errorf("expected address localhost:8080, got %s", dp.address)
	}
}

func TestFromConfig_Dgraph_WithTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "dgraph-timeout",
		Type:    "dgraph",
		Address: "localhost:8080",
		Timeout: "3s",
	}
	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	dp, ok := p.(*DgraphProbe)
	if !ok {
		t.Fatalf("expected *DgraphProbe, got %T", p)
	}
	if dp.timeout != 3*time.Second {
		t.Errorf("expected 3s timeout, got %v", dp.timeout)
	}
}
