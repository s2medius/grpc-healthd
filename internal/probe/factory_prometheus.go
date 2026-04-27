package probe

import (
	"fmt"
	"time"

	"github.com/your-org/grpc-healthd/internal/config"
)

func init() {
	registerFactory("prometheus", newPrometheusFromConfig)
}

func newPrometheusFromConfig(cfg config.ProbeConfig) (Prober, error) {
	address := cfg.Address
	if address == "" {
		return nil, fmt.Errorf("prometheus probe %q: address is required", cfg.Name)
	}

	var timeout time.Duration
	if cfg.Timeout != "" {
		var err error
		timeout, err = time.ParseDuration(cfg.Timeout)
		if err != nil {
			return nil, fmt.Errorf("prometheus probe %q: invalid timeout %q: %w", cfg.Name, cfg.Timeout, err)
		}
	}

	metricName := ""
	if v, ok := cfg.Options["metric_name"]; ok {
		metricName = v
	}

	return NewPrometheusProbe(address, metricName, timeout), nil
}
