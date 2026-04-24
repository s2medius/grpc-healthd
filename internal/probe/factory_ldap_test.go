package probe

import (
	"testing"
	"time"

	"grpc-healthd/internal/config"
)

func TestFromConfig_LDAP(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "ldap-check",
		Type:    "ldap",
		Address: "localhost:389",
	}
	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
	lp, ok := p.(*LDAPProbe)
	if !ok {
		t.Fatalf("expected *LDAPProbe, got %T", p)
	}
	if lp.address != "localhost:389" {
		t.Errorf("unexpected address: %s", lp.address)
	}
	if lp.timeout != defaultTimeout {
		t.Errorf("expected defaultTimeout, got %v", lp.timeout)
	}
}

func TestFromConfig_LDAP_WithTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "ldap-timed",
		Type:    "ldap",
		Address: "ldap.example.com:389",
		Timeout: 3 * time.Second,
	}
	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lp, ok := p.(*LDAPProbe)
	if !ok {
		t.Fatalf("expected *LDAPProbe, got %T", p)
	}
	if lp.timeout != 3*time.Second {
		t.Errorf("expected 3s timeout, got %v", lp.timeout)
	}
}
