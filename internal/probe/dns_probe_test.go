package probe

import (
	"testing"
	"time"
)

func TestDNSProbe_Healthy(t *testing.T) {
	p := NewDNSProbe("localhost", 5*time.Second)
	result := p.Probe()
	if result.Status != StatusHealthy {
		t.Errorf("expected healthy, got %s: %s", result.Status, result.Message)
	}
	if result.Duration <= 0 {
		t.Error("expected positive duration")
	}
}

func TestDNSProbe_Unhealthy_InvalidHost(t *testing.T) {
	p := NewDNSProbe("this.host.does.not.exist.invalid", 3*time.Second)
	result := p.Probe()
	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", result.Status)
	}
	if result.Message == "" {
		t.Error("expected non-empty message on failure")
	}
}

func TestNewDNSProbe_DefaultTimeout(t *testing.T) {
	p := NewDNSProbe("localhost", 0)
	if p.timeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %s", p.timeout)
	}
}

func TestDNSProbe_CustomTimeout(t *testing.T) {
	p := NewDNSProbe("localhost", 2*time.Second)
	if p.timeout != 2*time.Second {
		t.Errorf("expected 2s timeout, got %s", p.timeout)
	}
}

func TestDNSProbe_DurationRecorded(t *testing.T) {
	p := NewDNSProbe("localhost", 5*time.Second)
	result := p.Probe()
	if result.Duration == 0 {
		t.Error("expected non-zero duration")
	}
}
