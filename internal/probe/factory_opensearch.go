package probe

import (
	"fmt"
	"time"

	"github.com/your-org/grpc-healthd/internal/config"
)

func init() {
	registerFactory("opensearch", newOpenSearchFromConfig)
}

func newOpenSearchFromConfig(cfg config.ProbeConfig) (Prober, error) {
	if cfg.Address == "" {
		return nil, fmt.Errorf("opensearch probe %q: address is required", cfg.Name)
	}

	var timeout time.Duration
	if cfg.Timeout != "" {
		var err error
		timeout, err = time.ParseDuration(cfg.Timeout)
		if err != nil {
			return nil, fmt.Errorf("opensearch probe %q: invalid timeout %q: %w", cfg.Name, cfg.Timeout, err)
		}
	}

	return NewOpenSearchProbe(cfg.Address, timeout), nil
}
