package probe

import (
	"fmt"
	"time"

	"github.com/grpc-healthd/internal/config"
)

func init() {
	registerFactory("timescaledb", newTimescaleDBFromConfig)
}

func newTimescaleDBFromConfig(cfg config.ProbeConfig) (Prober, error) {
	var timeout time.Duration
	if cfg.Timeout != "" {
		var err error
		timeout, err = time.ParseDuration(cfg.Timeout)
		if err != nil {
			return nil, fmt.Errorf("timescaledb probe %q: invalid timeout %q: %w", cfg.Name, cfg.Timeout, err)
		}
	}
	return NewTimescaleDBProbe(cfg.Address, timeout), nil
}
