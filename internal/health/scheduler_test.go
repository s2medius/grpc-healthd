package health

import (
	"context"
	"testing"
	"time"

	"github.com/yourorg/grpc-healthd/internal/config"
	"github.com/yourorg/grpc-healthd/internal/probe"
)

func TestNewScheduler_NotNil(t *testing.T) {
	checker := NewChecker()
	s := NewScheduler(checker, nil)
	if s == nil {
		t.Fatal("expected non-nil scheduler")
	}
}

func TestScheduler_UpdatesStatus(t *testing.T) {
	// Use a real TCP server started by net.Listen so the probe succeeds.
	ln := startTCPListener(t)
	defer ln.Close()

	addr := ln.Addr().String()
	probes := []config.ProbeConfig{
		{
			Name:     "tcp-test",
			Type:     "tcp",
			Address:  addr,
			Interval: 50 * time.Millisecond,
			Timeout:  200 * time.Millisecond,
		},
	}

	checker := NewChecker()
	s := NewScheduler(checker, probes)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.Start(ctx)
	time.Sleep(200 * time.Millisecond)

	result, err := checker.GetStatus("tcp-test")
	if err != nil {
		t.Fatalf("expected status for tcp-test: %v", err)
	}
	if result.Status != probe.StatusHealthy {
		t.Errorf("expected healthy, got %s", result.Status)
	}
}

func TestScheduler_Stop(t *testing.T) {
	checker := NewChecker()
	s := NewScheduler(checker, nil)
	ctx := context.Background()
	s.Start(ctx)
	// Should not panic or block.
	s.Stop()
}

func TestScheduler_InvalidProbeSkipped(t *testing.T) {
	probes := []config.ProbeConfig{
		{Name: "bad", Type: "unknown", Address: "x"},
	}
	checker := NewChecker()
	s := NewScheduler(checker, probes)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	s.Start(ctx)
	time.Sleep(50 * time.Millisecond)
	_, err := checker.GetStatus("bad")
	if err == nil {
		t.Error("expected error for unregistered probe")
	}
}

// TestScheduler_MultipleProbes verifies that the scheduler correctly tracks
// status for multiple probes running concurrently.
func TestScheduler_MultipleProbes(t *testing.T) {
	ln1 := startTCPListener(t)
	defer ln1.Close()
	ln2 := startTCPListener(t)
	defer ln2.Close()

	probes := []config.ProbeConfig{
		{
			Name:     "svc-a",
			Type:     "tcp",
			Address:  ln1.Addr().String(),
			Interval: 50 * time.Millisecond,
			Timeout:  200 * time.Millisecond,
		},
		{
			Name:     "svc-b",
			Type:     "tcp",
			Address:  ln2.Addr().String(),
			Interval: 50 * time.Millisecond,
			Timeout:  200 * time.Millisecond,
		},
	}

	checker := NewChecker()
	s := NewScheduler(checker, probes)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.Start(ctx)
	time.Sleep(200 * time.Millisecond)

	for _, name := range []string{"svc-a", "svc-b"} {
		result, err := checker.GetStatus(name)
		if err != nil {
			t.Fatalf("expected status for %s: %v", name, err)
		}
		if result.Status != probe.StatusHealthy {
			t.Errorf("probe %s: expected healthy, got %s", name, result.Status)
		}
	}
}
