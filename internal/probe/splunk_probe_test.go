package probe_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/your-org/grpc-healthd/internal/probe"
)

func TestSplunkProbe_Healthy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"health":"green"}`))
	}))
	defer ts.Close()

	p := probe.NewSplunkProbe(ts.URL, 2*time.Second)
	res := p.Execute(context.Background())

	if res.Status != probe.StatusHealthy {
		t.Errorf("expected healthy, got %s: %v", res.Status, res.Err)
	}
	if res.Duration <= 0 {
		t.Error("expected positive duration")
	}
}

func TestSplunkProbe_Unhealthy_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	p := probe.NewSplunkProbe(ts.URL, 2*time.Second)
	res := p.Execute(context.Background())

	if res.Status != probe.StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", res.Status)
	}
	if res.Err == nil {
		t.Error("expected non-nil error")
	}
}

func TestSplunkProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := probe.NewSplunkProbe("http://127.0.0.1:19999", 500*time.Millisecond)
	res := p.Execute(context.Background())

	if res.Status != probe.StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", res.Status)
	}
}

func TestNewSplunkProbe_DefaultTimeout(t *testing.T) {
	p := probe.NewSplunkProbe("http://splunk:8089", 0)
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestSplunkProbe_DurationRecorded(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	p := probe.NewSplunkProbe(ts.URL, 2*time.Second)
	res := p.Execute(context.Background())

	if res.Duration < 10*time.Millisecond {
		t.Errorf("expected duration >= 10ms, got %s", res.Duration)
	}
}
