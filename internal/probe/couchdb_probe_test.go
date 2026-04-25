package probe

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestCouchDBProbe_Healthy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/_up" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer ts.Close()

	addr := strings.TrimPrefix(ts.URL, "http://")
	p := NewCouchDBProbe(addr, 2*time.Second)
	result := p.Execute(context.Background())

	if !result.Healthy {
		t.Errorf("expected healthy, got: %s", result.Message)
	}
}

func TestCouchDBProbe_Unhealthy_ServiceUnavailable(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "maintenance_mode"})
	}))
	defer ts.Close()

	addr := strings.TrimPrefix(ts.URL, "http://")
	p := NewCouchDBProbe(addr, 2*time.Second)
	result := p.Execute(context.Background())

	if result.Healthy {
		t.Error("expected unhealthy for 503 response")
	}
}

func TestCouchDBProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewCouchDBProbe("127.0.0.1:19999", 500*time.Millisecond)
	result := p.Execute(context.Background())

	if result.Healthy {
		t.Error("expected unhealthy for refused connection")
	}
}

func TestNewCouchDBProbe_DefaultTimeout(t *testing.T) {
	p := NewCouchDBProbe("localhost:5984", 0)
	if p.timeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %v", p.timeout)
	}
}

func TestCouchDBProbe_DurationRecorded(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer ts.Close()

	addr := strings.TrimPrefix(ts.URL, "http://")
	p := NewCouchDBProbe(addr, 2*time.Second)
	start := time.Now()
	result := p.Execute(context.Background())
	elapsed := time.Since(start)

	if !result.Healthy {
		t.Errorf("expected healthy, got: %s", result.Message)
	}
	if elapsed < 0 {
		t.Error("elapsed time should be non-negative")
	}
}
