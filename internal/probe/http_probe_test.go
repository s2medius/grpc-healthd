package probe

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPProbe_Healthy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	p := NewHTTPProbe(ts.URL, 0)
	result := p.Check()

	if result.Status != StatusHealthy {
		t.Errorf("expected healthy, got %s: %s", result.Status, result.Message)
	}
}

func TestHTTPProbe_Unhealthy_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	p := NewHTTPProbe(ts.URL, 0)
	result := p.Check()

	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", result.Status)
	}
}

func TestHTTPProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewHTTPProbe("http://127.0.0.1:19999", 500*time.Millisecond)
	result := p.Check()

	if result.Status != StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", result.Status)
	}
}

func TestNewHTTPProbe_DefaultTimeout(t *testing.T) {
	p := NewHTTPProbe("http://example.com", 0)
	if p.Timeout != DefaultTimeout {
		t.Errorf("expected default timeout %v, got %v", DefaultTimeout, p.Timeout)
	}
}

func TestHTTPProbe_DurationRecorded(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	p := NewHTTPProbe(ts.URL, 0)
	result := p.Check()

	if result.Duration <= 0 {
		t.Errorf("expected positive duration, got %v", result.Duration)
	}
}
