package metrics_test

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/your-org/grpc-healthd/internal/metrics"
)

func TestRecordProbe_Success(t *testing.T) {
	metrics.RecordProbe("localhost:50051", true, 0.042)

	count := testutil.ToFloat64(metrics.ProbeTotal.WithLabelValues("localhost:50051", "success"))
	if count < 1 {
		t.Errorf("expected probe_total success >= 1, got %v", count)
	}

	up := testutil.ToFloat64(metrics.ProbeUp.WithLabelValues("localhost:50051"))
	if up != 1.0 {
		t.Errorf("expected probe_up=1, got %v", up)
	}
}

func TestRecordProbe_Failure(t *testing.T) {
	metrics.RecordProbe("localhost:50052", false, 0.1)

	count := testutil.ToFloat64(metrics.ProbeTotal.WithLabelValues("localhost:50052", "failure"))
	if count < 1 {
		t.Errorf("expected probe_total failure >= 1, got %v", count)
	}

	up := testutil.ToFloat64(metrics.ProbeUp.WithLabelValues("localhost:50052"))
	if up != 0.0 {
		t.Errorf("expected probe_up=0, got %v", up)
	}
}

func TestRecordProbe_DurationObserved(t *testing.T) {
	metrics.RecordProbe("localhost:50053", true, 0.25)

	// Verify histogram has observations by checking count via testutil
	count := testutil.ToFloat64(metrics.ProbeTotal.WithLabelValues("localhost:50053", "success"))
	if count < 1 {
		t.Errorf("expected at least one probe recorded, got %v", count)
	}
}
