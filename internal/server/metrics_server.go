package server

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsServer serves Prometheus metrics over HTTP.
type MetricsServer struct {
	server *http.Server
}

// NewMetricsServer creates a new MetricsServer listening on addr.
func NewMetricsServer(addr string) *MetricsServer {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return &MetricsServer{
		server: &http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

// ListenAndServe starts the HTTP metrics server.
func (m *MetricsServer) ListenAndServe() error {
	return m.server.ListenAndServe()
}

// Shutdown gracefully stops the metrics server.
func (m *MetricsServer) Shutdown(ctx context.Context) error {
	return m.server.Shutdown(ctx)
}

// Addr returns the configured listen address.
func (m *MetricsServer) Addr() string {
	return m.server.Addr
}
