package probe

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/yourorg/grpc-healthd/internal/metrics"
	"golang.org/x/net/websocket"
)

// WebSocketProbe checks health by establishing a WebSocket connection.
type WebSocketProbe struct {
	address string
	timeout time.Duration
}

// NewWebSocketProbe creates a new WebSocketProbe.
// address should be a ws:// or wss:// URL.
func NewWebSocketProbe(address string, timeout time.Duration) *WebSocketProbe {
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	return &WebSocketProbe{
		address: address,
		timeout: timeout,
	}
}

// Probe attempts a WebSocket handshake and returns the result.
func (p *WebSocketProbe) Probe(ctx context.Context) Result {
	start := time.Now()

	origin := "http://localhost"
	dialer := websocket.Dial

	type dialResult struct {
		conn *websocket.Conn
		err  error
	}

	ch := make(chan dialResult, 1)
	go func() {
		conn, err := dialer(p.address, "", origin)
		ch <- dialResult{conn, err}
	}()

	var err error
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case res := <-ch:
		if res.err == nil && res.conn != nil {
			_ = res.conn.Close()
		}
		err = res.err
	}

	duration := time.Since(start)
	if err != nil {
		metrics.RecordProbe(p.address, "websocket", false, duration)
		return Result{Status: StatusUnhealthy, Error: fmt.Errorf("websocket probe failed: %w", err)}
	}

	metrics.RecordProbe(p.address, "websocket", true, duration)
	return Result{Status: StatusHealthy}
}

// compile-time check
var _ interface {
	Probe(context.Context) Result
} = (*WebSocketProbe)(nil)

// httpOriginHeader is used internally for handshake origin.
var _ = http.Header{}
