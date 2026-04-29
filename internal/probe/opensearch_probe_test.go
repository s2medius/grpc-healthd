package probe

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func startFakeOpenSearch(t *testing.T, status string, httpStatus int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/_cluster/health" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(httpStatus)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": status})
	}))
}

func TestOpenSearchProbe_Healthy_Green(t *testing.T) {
	srv := startFakeOpenSearch(t, "green", http.StatusOK)
	defer srv.Close()

	p := NewOpenSearchProbe(srv.Listener.Addr().String(), 0)
	res := p.Probe(context.Background())
	if res.Status != StatusHealthy {
		t.Fatalf("expected healthy, got %s: %v", res.Status, res.Error)
	}
}

func TestOpenSearchProbe_Healthy_Yellow(t *testing.T) {
	srv := startFakeOpenSearch(t, "yellow", http.StatusOK)
	defer srv.Close()

	p := NewOpenSearchProbe(srv.Listener.Addr().String(), 0)
	res := p.Probe(context.Background())
	if res.Status != StatusHealthy {
		t.Fatalf("expected healthy, got %s: %v", res.Status, res.Error)
	}
}

func TestOpenSearchProbe_Unhealthy_Red(t *testing.T) {
	srv := startFakeOpenSearch(t, "red", http.StatusOK)
	defer srv.Close()

	p := NewOpenSearchProbe(srv.Listener.Addr().String(), 0)
	res := p.Probe(context.Background())
	if res.Status != StatusUnhealthy {
		t.Fatalf("expected unhealthy, got %s", res.Status)
	}
}

func TestOpenSearchProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewOpenSearchProbe("127.0.0.1:19200", 500*time.Millisecond)
	res := p.Probe(context.Background())
	if res.Status != StatusUnhealthy {
		t.Fatalf("expected unhealthy, got %s", res.Status)
	}
}

func TestNewOpenSearchProbe_DefaultTimeout(t *testing.T) {
	p := NewOpenSearchProbe("localhost:9200", 0)
	if p.timeout != DefaultTimeout {
		t.Fatalf("expected default timeout %v, got %v", DefaultTimeout, p.timeout)
	}
}

func TestOpenSearchProbe_DurationRecorded(t *testing.T) {
	srv := startFakeOpenSearch(t, "green", http.StatusOK)
	defer srv.Close()

	p := NewOpenSearchProbe(srv.Listener.Addr().String(), 0)
	res := p.Probe(context.Background())
	if res.Duration <= 0 {
		t.Fatalf("expected positive duration, got %v", res.Duration)
	}
}
