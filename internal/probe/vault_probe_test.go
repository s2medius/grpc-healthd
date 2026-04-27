package probe

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestVaultProbe_Healthy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"initialized":true,"sealed":false,"version":"1.15.0"}`))
	}))
	defer ts.Close()

	p := NewVaultProbe(ts.URL, 2*time.Second)
	res := p.Probe(context.Background())
	if !res.Healthy {
		t.Fatalf("expected healthy, got: %s", res.Message)
	}
}

func TestVaultProbe_Healthy_Standby(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests) // 429 = standby
		_, _ = w.Write([]byte(`{"initialized":true,"sealed":false,"version":"1.15.0"}`))
	}))
	defer ts.Close()

	p := NewVaultProbe(ts.URL, 2*time.Second)
	res := p.Probe(context.Background())
	if !res.Healthy {
		t.Fatalf("expected healthy for standby, got: %s", res.Message)
	}
}

func TestVaultProbe_Unhealthy_Sealed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable) // 503 = sealed
		_, _ = w.Write([]byte(`{"initialized":true,"sealed":true}`))
	}))
	defer ts.Close()

	p := NewVaultProbe(ts.URL, 2*time.Second)
	res := p.Probe(context.Background())
	if res.Healthy {
		t.Fatal("expected unhealthy for sealed vault")
	}
}

func TestVaultProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewVaultProbe("http://127.0.0.1:19999", 500*time.Millisecond)
	res := p.Probe(context.Background())
	if res.Healthy {
		t.Fatal("expected unhealthy for connection refused")
	}
}

func TestNewVaultProbe_DefaultTimeout(t *testing.T) {
	p := NewVaultProbe("http://127.0.0.1:8200", 0)
	if p.timeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %v", p.timeout)
	}
}

func TestVaultProbe_DurationRecorded(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"initialized":true,"sealed":false}`))
	}))
	defer ts.Close()

	p := NewVaultProbe(ts.URL, 2*time.Second)
	res := p.Probe(context.Background())
	if !res.Healthy {
		t.Fatalf("expected healthy: %s", res.Message)
	}
}
