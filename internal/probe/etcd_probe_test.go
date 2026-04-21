package probe_test

import (
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/grpc-healthd/internal/probe"
)

func startFakeEtcd(t *testing.T, healthy bool) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if healthy {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]string{"health": "true"})
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			_ = json.NewEncoder(w).Encode(map[string]string{"health": "false"})
		}
	}))
	return ts
}

func TestEtcdProbe_Healthy(t *testing.T) {
	ts := startFakeEtcd(t, true)
	defer ts.Close()

	p := probe.NewEtcdProbe(ts.Listener.Addr().String(), 5*time.Second)
	result := p.Probe()
	if !result.Healthy {
		t.Errorf("expected healthy, got unhealthy: %s", result.Message)
	}
}

func TestEtcdProbe_Unhealthy_ServiceUnavailable(t *testing.T) {
	ts := startFakeEtcd(t, false)
	defer ts.Close()

	p := probe.NewEtcdProbe(ts.Listener.Addr().String(), 5*time.Second)
	result := p.Probe()
	if result.Healthy {
		t.Error("expected unhealthy, got healthy")
	}
}

func TestEtcdProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	// Find a free port, then don't listen on it
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	ln.Close()

	p := probe.NewEtcdProbe(addr, 500*time.Millisecond)
	result := p.Probe()
	if result.Healthy {
		t.Error("expected unhealthy for refused connection")
	}
}

func TestNewEtcdProbe_DefaultTimeout(t *testing.T) {
	p := probe.NewEtcdProbe("localhost:2379", 0)
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestEtcdProbe_DurationRecorded(t *testing.T) {
	ts := startFakeEtcd(t, true)
	defer ts.Close()

	p := probe.NewEtcdProbe(ts.Listener.Addr().String(), 5*time.Second)
	result := p.Probe()
	if result.Duration <= 0 {
		t.Errorf("expected positive duration, got %v", result.Duration)
	}
}
