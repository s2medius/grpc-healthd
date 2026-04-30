package probe_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/yourusername/grpc-healthd/internal/probe"
)

func TestGraphQLProbe_Healthy(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]string{"__typename": "Query"},
		})
	}))
	defer ts.Close()

	p := probe.NewGraphQLProbe(ts.URL, 5*time.Second)
	if got := p.Check(context.Background()); got != probe.StatusHealthy {
		t.Errorf("expected Healthy, got %s", got)
	}
}

func TestGraphQLProbe_Unhealthy_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	p := probe.NewGraphQLProbe(ts.URL, 5*time.Second)
	if got := p.Check(context.Background()); got != probe.StatusUnhealthy {
		t.Errorf("expected Unhealthy, got %s", got)
	}
}

func TestGraphQLProbe_Unhealthy_GraphQLErrors(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": []map[string]string{{"message": "some error"}},
		})
	}))
	defer ts.Close()

	p := probe.NewGraphQLProbe(ts.URL, 5*time.Second)
	if got := p.Check(context.Background()); got != probe.StatusUnhealthy {
		t.Errorf("expected Unhealthy, got %s", got)
	}
}

func TestGraphQLProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := probe.NewGraphQLProbe("http://127.0.0.1:1", 500*time.Millisecond)
	if got := p.Check(context.Background()); got != probe.StatusUnhealthy {
		t.Errorf("expected Unhealthy, got %s", got)
	}
}

func TestNewGraphQLProbe_DefaultTimeout(t *testing.T) {
	p := probe.NewGraphQLProbe("http://localhost:8080/graphql", 0)
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestGraphQLProbe_DurationRecorded(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]string{"__typename": "Query"},
		})
	}))
	defer ts.Close()

	p := probe.NewGraphQLProbe(ts.URL, 5*time.Second)
	// Should not panic and should record metrics
	p.Check(context.Background())
}
