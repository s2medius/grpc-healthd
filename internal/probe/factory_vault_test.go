package probe

import (
	"testing"

	"github.com/grpc-healthd/internal/config"
)

func TestFromConfig_Vault(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "vault-check",
		Type:    "vault",
		Address: "http://vault.internal:8200",
	}

	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil prober")
	}
	vp, ok := p.(*VaultProbe)
	if !ok {
		t.Fatalf("expected *VaultProbe, got %T", p)
	}
	if vp.address != "http://vault.internal:8200" {
		t.Errorf("unexpected address: %s", vp.address)
	}
}

func TestFromConfig_Vault_WithTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "vault-check",
		Type:    "vault",
		Address: "http://vault.internal:8200",
		Timeout: "3s",
	}

	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	vp := p.(*VaultProbe)
	if vp.timeout.Seconds() != 3 {
		t.Errorf("expected 3s timeout, got %v", vp.timeout)
	}
}

func TestFromConfig_Vault_InvalidTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "vault-check",
		Type:    "vault",
		Address: "http://vault.internal:8200",
		Timeout: "not-a-duration",
	}

	_, err := FromConfig(cfg)
	if err == nil {
		t.Fatal("expected error for invalid timeout")
	}
}
