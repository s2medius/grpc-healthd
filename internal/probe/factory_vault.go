package probe

import (
	"fmt"
	"time"

	"github.com/grpc-healthd/internal/config"
)

func init() {
	registerFactory("vault", newVaultFromConfig)
}

func newVaultFromConfig(cfg config.ProbeConfig) (Prober, error) {
	if cfg.Address == "" {
		return nil, fmt.Errorf("vault probe %q: address is required", cfg.Name)
	}

	var timeout time.Duration
	if cfg.Timeout != "" {
		var err error
		timeout, err = time.ParseDuration(cfg.Timeout)
		if err != nil {
			return nil, fmt.Errorf("vault probe %q: invalid timeout %q: %w", cfg.Name, cfg.Timeout, err)
		}
	}

	return NewVaultProbe(cfg.Address, timeout), nil
}
