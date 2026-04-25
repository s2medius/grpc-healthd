package probe

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func jsonHandler(body map[string]interface{}, status int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(body)
	}
}

func TestHTTPJSONProbe_Healthy(t *testing.T) {
	ts := httptest.NewServer(jsonHandler(map[string]interface{}{"status": "ok"}, http.StatusOK))
	defer ts.Close()

	p := NewHTTPJSONProbe(ts.URL, "status", "ok", 3*time.Second)
	res := p.Probe(context.Background())
	if !res.Healthy {
		t.Fatalf("expected healthy, got: %s", res.Message)
	}
}

func TestHTTPJSONProbe_Unhealthy_WrongValue(t *testing.T) {
	ts := httptest.NewServer(jsonHandler(map[string]interface{}{"status": "degraded"}, http.StatusOK))
	defer ts.Close()

	p := NewHTTPJSONProbe(ts.URL, "status", "ok", 3*time.Second)
	res := p.Probe(context.Background())
	if res.Healthy {
		t.Fatal("expected unhealthy due to wrong value")
	}
}

func TestHTTPJSONProbe_Unhealthy_MissingKey(t *testing.T) {
	ts := httptest.NewServer(jsonHandler(map[string]interface{}{"health": "ok"}, http.StatusOK))
	defer ts.Close()

	p := NewHTTPJSONProbe(ts.URL, "status", "ok", 3*time.Second)
	res := p.Probe(context.Background())
	if res.Healthy {
		t.Fatal("expected unhealthy due to missing key")
	}
}

func TestHTTPJSONProbe_Unhealthy_ServerError(t *testing.T) {
	ts := httptest.NewServer(jsonHandler(map[string]interface{}{"status": "ok"}, http.StatusInternalServerError))
	defer ts.Close()

	p := NewHTTPJSONProbe(ts.URL, "status", "ok", 3*time.Second)
	res := p.Probe(context.Background())
	if res.Healthy {
		t.Fatal("expected unhealthy due to 500 status")
	}
}

func TestHTTPJSONProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewHTTPJSONProbe("http://127.0.0.1:19999", "status", "ok", time.Second)
	res := p.Probe(context.Background())
	if res.Healthy {
		t.Fatal("expected unhealthy due to connection refused")
	}
}

func TestNewHTTPJSONProbe_DefaultTimeout(t *testing.T) {
	p := NewHTTPJSONProbe("http://example.com", "status", "ok", 0)
	if p.timeout != 5*time.Second {
		t.Fatalf("expected default timeout 5s, got %v", p.timeout)
	}
}

func TestHTTPJSONProbe_Unhealthy_InvalidJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("not-json"))
	}))
	defer ts.Close()

	p := NewHTTPJSONProbe(ts.URL, "status", "ok", 3*time.Second)
	res := p.Probe(context.Background())
	if res.Healthy {
		t.Fatal("expected unhealthy due to invalid JSON body")
	}
}
