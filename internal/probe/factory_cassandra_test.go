package probe

import (
	"testing"
	"time"

	"github.com/patrickdappollonio/grpc-healthd/internal/config"
)

func TestFromConfig_Cassandra(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "cassandra-test",
		Type:    "cassandra",
		Address: "localhost:9042",
	}
	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cp, ok := p.(*CassandraProbe)
	if !ok {
		t.Fatalf("expected *CassandraProbe, got %T", p)
	}
	if cp.address != "localhost:9042" {
		t.Errorf("expected address localhost:9042, got %s", cp.address)
	}
	if cp.timeout != DefaultTimeout {
		t.Errorf("expected default timeout, got %v", cp.timeout)
	}
}

func TestFromConfig_Cassandra_WithTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "cassandra-timeout",
		Type:    "cassandra",
		Address: "localhost:9042",
		Timeout: 5 * time.Second,
	}
	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cp, ok := p.(*CassandraProbe)
	if !ok {
		t.Fatalf("expected *CassandraProbe, got %T", p)
	}
	if cp.timeout != 5*time.Second {
		t.Errorf("expected 5s timeout, got %v", cp.timeout)
	}
}
