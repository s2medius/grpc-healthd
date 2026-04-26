package probe

import (
	"fmt"
	"time"

	"github.com/your-org/grpc-healthd/internal/config"
)

func init() {
	registerFactory("dgraph", newDgraphFromConfig)
}

func newDgraphFromConfig(cfg config.ProbeConfig) (Prober, error) {
	var timeout time.Duration
	if cfg.Timeout != "" {
		var err error
		timeout, err = time.ParseDuration(cfg.Timeout)
		if err != nil {
			return nil, fmt.Errorf("dgraph probe %q: invalid timeout %q: %w", cfg.Name, cfg.Timeout, err)
		}
	}
	return NewDgraphProbe(cfg.Address, timeout), nil
}
