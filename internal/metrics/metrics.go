package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ProbeTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "grpc_healthd",
			Name:      "probe_total",
			Help:      "Total number of probes executed.",
		},
		[]string{"target", "status"},
	)

	ProbeDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "grpc_healthd",
			Name:      "probe_duration_seconds",
			Help:      "Duration of probe executions in seconds.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"target"},
	)

	ProbeUp = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "grpc_healthd",
			Name:      "probe_up",
			Help:      "Whether the last probe succeeded (1) or failed (0).",
		},
		[]string{"target"},
	)
)

// RecordProbe records the result of a probe for a given target.
func RecordProbe(target string, healthy bool, durationSeconds float64) {
	status := "failure"
	up := 0.0
	if healthy {
		status = "success"
		up = 1.0
	}
	ProbeTotal.WithLabelValues(target, status).Inc()
	ProbeDuration.WithLabelValues(target).Observe(durationSeconds)
	ProbeUp.WithLabelValues(target).Set(up)
}
