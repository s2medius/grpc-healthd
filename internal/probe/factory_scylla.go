package probe

import (
	"fmt"
	"time"

	"github.com/your-org/grpc-healthd/internal/config"
)

func init() {
	RegisterFactory("scylla", newScyllaFromConfig)
}

func newScyllaFromConfig(cfg config.ProbeConfig) (Probe, error) {
	timeout := 5 * time.Second
	if cfg.Timeout != "" {
		var err error
		timeout, err = time.ParseDuration(cfg.Timeout)
		if err != nil {
			return nil, fmt.Errorf("scylla probe %q: invalid timeout %q: %w", cfg.Name, cfg.Timeout, err)
		}
	}
	return NewScyllaProbe(cfg.Address, timeout), nil
}
