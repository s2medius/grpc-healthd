package probe

import (
	"fmt"
	"time"

	"github.com/your-org/grpc-healthd/internal/config"
)

func init() {
	RegisterFactory("splunk", newSplunkFromConfig)
}

func newSplunkFromConfig(cfg config.ProbeConfig) (Probe, error) {
	if cfg.Address == "" {
		return nil, fmt.Errorf("splunk probe %q: address is required", cfg.Name)
	}

	var timeout time.Duration
	if cfg.Timeout != "" {
		var err error
		timeout, err = time.ParseDuration(cfg.Timeout)
		if err != nil {
			return nil, fmt.Errorf("splunk probe %q: invalid timeout %q: %w", cfg.Name, cfg.Timeout, err)
		}
	}

	return NewSplunkProbe(cfg.Address, timeout), nil
}
