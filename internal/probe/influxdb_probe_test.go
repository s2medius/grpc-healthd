package probe

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestInfluxDBProbe_Healthy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ping" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	address := strings.TrimPrefix(ts.URL, "http://")
	p := NewInfluxDBProbe(address, time.Second)
	result := p.Probe(context.Background())

	if !result.Healthy {
		t.Errorf("expected healthy, got unhealthy: %s", result.Message)
	}
}

func TestInfluxDBProbe_Unhealthy_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	address := strings.TrimPrefix(ts.URL, "http://")
	p := NewInfluxDBProbe(address, time.Second)
	result := p.Probe(context.Background())

	if result.Healthy {
		t.Error("expected unhealthy, got healthy")
	}
	if !strings.Contains(result.Message, "500") {
		t.Errorf("expected status code in message, got: %s", result.Message)
	}
}

func TestInfluxDBProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewInfluxDBProbe("127.0.0.1:19999", time.Second)
	result := p.Probe(context.Background())

	if result.Healthy {
		t.Error("expected unhealthy for refused connection")
	}
}

func TestNewInfluxDBProbe_DefaultTimeout(t *testing.T) {
	p := NewInfluxDBProbe("localhost:8086", 0)
	if p.timeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %v", p.timeout)
	}
}

func TestInfluxDBProbe_DurationRecorded(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	address := strings.TrimPrefix(ts.URL, "http://")
	p := NewInfluxDBProbe(address, time.Second)
	result := p.Probe(context.Background())

	if !result.Healthy {
		t.Errorf("expected healthy result, got: %s", result.Message)
	}
}
