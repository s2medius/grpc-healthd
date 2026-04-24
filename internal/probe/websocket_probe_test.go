package probe_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/grpc-healthd/internal/probe"
	"golang.org/x/net/websocket"
)

func startFakeWebSocket(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.Handle("/ws", websocket.Handler(func(ws *websocket.Conn) {
		_ = ws.Close()
	}))
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

func TestWebSocketProbe_Healthy(t *testing.T) {
	srv := startFakeWebSocket(t)
	addr := "ws://" + strings.TrimPrefix(srv.URL, "http://") + "/ws"

	p := probe.NewWebSocketProbe(addr, time.Second)
	result := p.Probe(context.Background())

	if result.Status != probe.StatusHealthy {
		t.Errorf("expected healthy, got %s: %v", result.Status, result.Error)
	}
}

func TestWebSocketProbe_Unhealthy_ConnectionRefused(t *testing.T) {
	p := probe.NewWebSocketProbe("ws://127.0.0.1:19999/ws", 500*time.Millisecond)
	result := p.Probe(context.Background())

	if result.Status != probe.StatusUnhealthy {
		t.Errorf("expected unhealthy, got %s", result.Status)
	}
}

func TestNewWebSocketProbe_DefaultTimeout(t *testing.T) {
	p := probe.NewWebSocketProbe("ws://localhost/ws", 0)
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestWebSocketProbe_DurationRecorded(t *testing.T) {
	srv := startFakeWebSocket(t)
	addr := "ws://" + strings.TrimPrefix(srv.URL, "http://") + "/ws"

	p := probe.NewWebSocketProbe(addr, time.Second)
	start := time.Now()
	p.Probe(context.Background())
	elapsed := time.Since(start)

	if elapsed > 2*time.Second {
		t.Errorf("probe took too long: %v", elapsed)
	}
}

func TestWebSocketProbe_ContextCancel(t *testing.T) {
	p := probe.NewWebSocketProbe("ws://10.255.255.1:9999/ws", 5*time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	result := p.Probe(ctx)
	if result.Status != probe.StatusUnhealthy {
		t.Errorf("expected unhealthy on context cancel, got %s", result.Status)
	}
}
