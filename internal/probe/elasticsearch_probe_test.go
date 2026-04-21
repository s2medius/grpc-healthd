package probe

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func startFakeElasticsearch(t *testing.T, status string, httpStatus int) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/_cluster/health" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatus)
		body := map[string]interface{}{
			"cluster_name": "test-cluster",
			"status":       status,
		}
		_ = json.NewEncoder(w).Encode(body)
	}))
	t.Cleanup(server.Close)
	return server
}

func TestElasticsearchProbe_Healthy_Green(t *testing.T) {
	server := startFakeElasticsearch(t, "green", http.StatusOK)
	addr := strings.TrimPrefix(server.URL, "http://")
	p := NewElasticsearchProbe(addr, 5*time.Second)
	result := p.Check()
	if !result.Healthy {
		t.Errorf("expected healthy, got unhealthy: %s", result.Message)
	}
}

func TestElasticsearchProbe_Healthy_Yellow(t *testing.T) {
	server := startFakeElasticsearch(t, "yellow", http.StatusOK)
	addr := strings.TrimPrefix(server.URL, "http://")
	p := NewElasticsearchProbe(addr, 5*time.Second)
	result := p.Check()
	if !result.Healthy {
		t.Errorf("expected healthy for yellow status, got unhealthy: %s", result.Message)
	}
}

func TestElasticsearchProbe_Unhealthy_Red(t *testing.T) {
	server := startFakeElasticsearch(t, "red", http.StatusOK)
	addr := strings.TrimPrefix(server.URL, "http://")
	p := NewElasticsearchProbe(addr, 5*time.Second)
	result := p.Check()
	if result.Healthy {
		t.Error("expected unhealthy for red cluster status")
	}
}

func TestElasticsearchProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := NewElasticsearchProbe("127.0.0.1:19299", 1*time.Second)
	result := p.Check()
	if result.Healthy {
		t.Error("expected unhealthy when connection refused")
	}
}

func TestNewElasticsearchProbe_DefaultTimeout(t *testing.T) {
	p := NewElasticsearchProbe("localhost:9200", 0)
	ep, ok := p.(*ElasticsearchProbe)
	if !ok {
		t.Fatal("expected *ElasticsearchProbe")
	}
	if ep.timeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %v", ep.timeout)
	}
}

func TestElasticsearchProbe_DurationRecorded(t *testing.T) {
	server := startFakeElasticsearch(t, "green", http.StatusOK)
	addr := strings.TrimPrefix(server.URL, "http://")
	p := NewElasticsearchProbe(addr, 5*time.Second)
	result := p.Check()
	if result.Duration <= 0 {
		t.Error("expected positive duration")
	}
}
