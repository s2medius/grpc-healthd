package probe

import (
	"context"
	"testing"
	"time"
)

func TestExecProbe_Healthy(t *testing.T) {
	p := NewExecProbe("true", nil, 0)
	result := p.Probe(context.Background())
	if result.Status != StatusHealthy {
		t.Errorf("expected healthy, got %s: %s", result.Status, result.Message)
	}
}

func TestExecProbe_Unhealthy_NonZeroExit(t *testing.T) {
	p := NewExecProbe("false", nil, 0)
	result := p.Probe(context.Background())
	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", result.Status)
	}
	if result.Message == "" {
		t.Error("expected non-empty message on failure")
	}
}

func TestExecProbe_Unhealthy_CommandNotFound(t *testing.T) {
	p := NewExecProbe("no-such-command-xyz", nil, 0)
	result := p.Probe(context.Background())
	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", result.Status)
	}
}

func TestNewExecProbe_DefaultTimeout(t *testing.T) {
	p := NewExecProbe("true", nil, 0)
	if p.timeout != DefaultTimeout {
		t.Errorf("expected default timeout %v, got %v", DefaultTimeout, p.timeout)
	}
}

func TestExecProbe_CustomArgs(t *testing.T) {
	p := NewExecProbe("echo", []string{"hello"}, time.Second)
	result := p.Probe(context.Background())
	if result.Status != StatusHealthy {
		t.Errorf("expected healthy, got %s: %s", result.Status, result.Message)
	}
}

func TestExecProbe_DurationRecorded(t *testing.T) {
	p := NewExecProbe("true", nil, 0)
	result := p.Probe(context.Background())
	if result.Duration <= 0 {
		t.Error("expected positive duration")
	}
}
