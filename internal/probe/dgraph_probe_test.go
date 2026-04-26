package probe

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDgraphProbe_Healthy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"instance":"zero","address":"localhost:5080","status":"healthy"}]`))
	}))
	defer ts.Close()

	p := NewDgraphProbe(ts.Listener.Addr().String(), time.Second)
	res := p.Probe(context.Background())
	if !res.Healthy {
		t.Fatalf("expected healthy, got err: %v", res.Err)
	}
}

func TestDgraphProbe_Unhealthy_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	p := NewDgraphProbe(ts.Listener.Addr().String(), time.Second)
	res := p.Probe(context.Background())
	if res.Healthy {
		t.Fatal("expected unhealthy for 503 response")
	}
}

func TestDgraphProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewDgraphProbe("127.0.0.1:19999", 300*time.Millisecond)
	res := p.Probe(context.Background())
	if res.Healthy {
		t.Fatal("expected unhealthy for refused connection")
	}
}

func TestNewDgraphProbe_DefaultTimeout(t *testing.T) {
	p := NewDgraphProbe("localhost:8080", 0)
	if p.timeout != 5*time.Second {
		t.Fatalf("expected default timeout 5s, got %v", p.timeout)
	}
}

func TestDgraphProbe_DurationRecorded(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	p := NewDgraphProbe(ts.Listener.Addr().String(), time.Second)
	res := p.Probe(context.Background())
	if res.Duration <= 0 {
		t.Fatal("expected positive duration")
	}
}
