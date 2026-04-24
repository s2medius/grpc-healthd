package probe

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func startFakeConsul(t *testing.T, healthy bool) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !healthy {
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"Config": map[string]string{"Datacenter": "dc1"}})
	}))
}

func TestConsulProbe_Healthy(t *testing.T) {
	srv := startFakeConsul(t, true)
	defer srv.Close()

	p := NewConsulProbe(srv.Listener.Addr().String(), time.Second)
	res := p.Check(context.Background())
	if res.Status != StatusHealthy {
		t.Fatalf("expected healthy, got %s: %v", res.Status, res.Err)
	}
}

func TestConsulProbe_Unhealthy_ServiceUnavailable(t *testing.T) {
	srv := startFakeConsul(t, false)
	defer srv.Close()

	p := NewConsulProbe(srv.Listener.Addr().String(), time.Second)
	res := p.Check(context.Background())
	if res.Status != StatusUnhealthy {
		t.Fatalf("expected unhealthy, got %s", res.Status)
	}
}

func TestConsulProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewConsulProbe("127.0.0.1:19999", 200*time.Millisecond)
	res := p.Check(context.Background())
	if res.Status != StatusUnhealthy {
		t.Fatalf("expected unhealthy, got %s", res.Status)
	}
	if res.Err == nil {
		t.Fatal("expected non-nil error")
	}
}

func TestNewConsulProbe_DefaultTimeout(t *testing.T) {
	p := NewConsulProbe("localhost:8500", 0)
	if p.timeout != DefaultTimeout {
		t.Fatalf("expected default timeout %v, got %v", DefaultTimeout, p.timeout)
	}
}

func TestConsulProbe_DurationRecorded(t *testing.T) {
	srv := startFakeConsul(t, true)
	defer srv.Close()

	p := NewConsulProbe(srv.Listener.Addr().String(), time.Second)
	res := p.Check(context.Background())
	if res.Duration <= 0 {
		t.Fatal("expected positive duration")
	}
}
