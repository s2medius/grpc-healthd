package probe

import (
	"fmt"
	"time"

	"github.com/your-org/grpc-healthd/internal/config"
)

func init() {
	registerFactory("valkey", newValkeyFromConfig)
}

// newValkeyFromConfig constructs a ValkeyProbe from a ProbeConfig.
func newValkeyFromConfig(cfg config.ProbeConfig) (Prober, error) {
	var timeout time.Duration
	if cfg.Timeout != "" {
		var err error
		timeout, err = time.ParseDuration(cfg.Timeout)
		if err != nil {
			return nil, fmt.Errorf("valkey probe %q: invalid timeout %q: %w", cfg.Name, cfg.Timeout, err)
		}
	}
	return NewValkeyProbe(cfg.Address, timeout), nil
}
