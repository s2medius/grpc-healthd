package probe

import (
	"testing"

	"github.com/your-org/grpc-healthd/internal/config"
)

func TestFromConfig_OpenSearch(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "opensearch-default",
		Type:    "opensearch",
		Address: "localhost:9200",
	}

	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	op, ok := p.(*OpenSearchProbe)
	if !ok {
		t.Fatalf("expected *OpenSearchProbe, got %T", p)
	}
	if op.address != cfg.Address {
		t.Errorf("expected address %q, got %q", cfg.Address, op.address)
	}
	if op.timeout != DefaultTimeout {
		t.Errorf("expected default timeout, got %v", op.timeout)
	}
}

func TestFromConfig_OpenSearch_WithTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "opensearch-custom",
		Type:    "opensearch",
		Address: "localhost:9200",
		Timeout: "3s",
	}

	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	op, ok := p.(*OpenSearchProbe)
	if !ok {
		t.Fatalf("expected *OpenSearchProbe, got %T", p)
	}
	if op.timeout.String() != "3s" {
		t.Errorf("expected timeout 3s, got %v", op.timeout)
	}
}

func TestFromConfig_OpenSearch_InvalidTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "opensearch-bad",
		Type:    "opensearch",
		Address: "localhost:9200",
		Timeout: "notaduration",
	}

	_, err := FromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for invalid timeout, got nil")
	}
}
