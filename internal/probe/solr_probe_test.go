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

func startFakeSolr(t *testing.T, status string, httpStatus int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatus)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": status})
	}))
}

func TestSolrProbe_Healthy(t *testing.T) {
	srv := startFakeSolr(t, "OK", http.StatusOK)
	defer srv.Close()

	addr := strings.TrimPrefix(srv.URL, "http://")
	p := NewSolrProbe(addr, time.Second)
	res := p.Check(context.Background())

	if res.Status != Healthy {
		t.Fatalf("expected Healthy, got %s (err: %v)", res.Status, res.Error)
	}
}

func TestSolrProbe_Unhealthy_BadStatus(t *testing.T) {
	srv := startFakeSolr(t, "FAILED", http.StatusOK)
	defer srv.Close()

	addr := strings.TrimPrefix(srv.URL, "http://")
	p := NewSolrProbe(addr, time.Second)
	res := p.Check(context.Background())

	if res.Status != Unhealthy {
		t.Fatalf("expected Unhealthy, got %s", res.Status)
	}
	if res.Error == nil {
		t.Fatal("expected non-nil error")
	}
}

func TestSolrProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewSolrProbe("127.0.0.1:19999", 200*time.Millisecond)
	res := p.Check(context.Background())

	if res.Status != Unhealthy {
		t.Fatalf("expected Unhealthy, got %s", res.Status)
	}
}

func TestNewSolrProbe_DefaultTimeout(t *testing.T) {
	p := NewSolrProbe("localhost:8983", 0)
	if p.timeout != 5*time.Second {
		t.Fatalf("expected default timeout 5s, got %s", p.timeout)
	}
}

func TestSolrProbe_DurationRecorded(t *testing.T) {
	srv := startFakeSolr(t, "OK", http.StatusOK)
	defer srv.Close()

	addr := strings.TrimPrefix(srv.URL, "http://")
	p := NewSolrProbe(addr, time.Second)

	// Should not panic and should complete without error.
	res := p.Check(context.Background())
	if res.Status != Healthy {
		t.Fatalf("expected Healthy, got %s", res.Status)
	}
}
