package probe

import (
	"fmt"
	"time"

	"github.com/yourusername/grpc-healthd/internal/config"
)

func init() {
	Register("graphql", newGraphQLFromConfig)
}

func newGraphQLFromConfig(cfg config.ProbeConfig) (Prober, error) {
	if cfg.Address == "" {
		return nil, fmt.Errorf("graphql probe %q: address is required", cfg.Name)
	}

	var timeout time.Duration
	if cfg.Timeout != "" {
		var err error
		timeout, err = time.ParseDuration(cfg.Timeout)
		if err != nil {
			return nil, fmt.Errorf("graphql probe %q: invalid timeout %q: %w", cfg.Name, cfg.Timeout, err)
		}
	}

	return NewGraphQLProbe(cfg.Address, timeout), nil
}
