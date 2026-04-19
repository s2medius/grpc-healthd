package health_test

import (
	"context"
	"testing"
	"time"

	"github.com/yourorg/grpc-healthd/internal/health"
	"github.com/yourorg/grpc-healthd/internal/probe"
)

type mockProbe struct {
	status probe.Status
}

func (m *mockProbe) Check(_ context.Context) probe.Result {
	return probe.Result{Status: m.status, Duration: time.Millisecond}
}

func TestNewChecker(t *testing.T) {
	c := health.NewChecker(time.Second)
	if c == nil {
		t.Fatal("expected non-nil checker")
	}
}

func TestRegisterAndGetStatus(t *testing.T) {
	c := health.NewChecker(time.Second)
	c.Register("svc", &mockProbe{status: probe.StatusHealthy})

	s, ok := c.GetStatus("svc")
	if !ok {
		t.Fatal("expected status to exist")
	}
	if s.Name != "svc" {
		t.Errorf("expected name=svc, got %s", s.Name)
	}
}

func TestGetStatus_NotFound(t *testing.T) {
	c := health.NewChecker(time.Second)
	_, ok := c.GetStatus("missing")
	if ok {
		t.Error("expected status not found")
	}
}

func TestChecker_RunUpdatesStatus(t *testing.T) {
	c := health.NewChecker(50 * time.Millisecond)
	c.Register("svc", &mockProbe{status: probe.StatusHealthy})

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	go c.Run(ctx)
	<-ctx.Done()

	s, ok := c.GetStatus("svc")
	if !ok {
		t.Fatal("expected status to exist after run")
	}
	if s.Status != probe.StatusHealthy {
		t.Errorf("expected healthy, got %s", s.Status)
	}
	if s.LastChecked.IsZero() {
		t.Error("expected LastChecked to be set")
	}
}

func TestChecker_UnhealthyProbe(t *testing.T) {
	c := health.NewChecker(50 * time.Millisecond)
	c.Register("bad", &mockProbe{status: probe.StatusUnhealthy})

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	go c.Run(ctx)
	<-ctx.Done()

	s, _ := c.GetStatus("bad")
	if s.Status != probe.StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", s.Status)
	}
}
