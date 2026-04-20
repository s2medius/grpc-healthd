package probe

import (
	"fmt"
	"time"

	"github.com/yourusername/grpc-healthd/internal/config"
)

// FromConfig constructs a Probe from a ProbeConfig.
func FromConfig(cfg config.ProbeConfig) (Probe, error) {
	var timeout time.Duration
	if cfg.Timeout != "" {
		var err error
		timeout, err = time.ParseDuration(cfg.Timeout)
		if err != nil {
			return nil, fmt.Errorf("probe %q: invalid timeout %q: %w", cfg.Name, cfg.Timeout, err)
		}
	}

	switch cfg.Type {
	case config.ProbeTypeTCP:
		return NewTCPProbe(cfg.Address, timeout), nil
	case config.ProbeTypeHTTP:
		return NewHTTPProbe(cfg.Address, timeout), nil
	case config.ProbeTypeDNS:
		return NewDNSProbe(cfg.Address, timeout), nil
	case config.ProbeTypeExec:
		if len(cfg.Command) == 0 {
			return nil, fmt.Errorf("probe %q: exec requires at least one command argument", cfg.Name)
		}
		return NewExecProbe(cfg.Command[0], cfg.Command[1:], timeout), nil
	case config.ProbeTypeGRPC:
		return NewGRPCProbe(cfg.Address, timeout), nil
	case config.ProbeTypeTLS:
		return NewTLSProbe(cfg.Address, timeout), nil
	case config.ProbeTypeICMP:
		return NewICMPProbe(cfg.Address, timeout), nil
	case config.ProbeTypeRedis:
		return NewRedisProbe(cfg.Address, timeout), nil
	case config.ProbeTypePostgres:
		return NewPostgresProbe(cfg.Address, timeout), nil
	case config.ProbeTypeMySQL:
		return NewMySQLProbe(cfg.Address, timeout), nil
	default:
		return nil, fmt.Errorf("probe %q: unsupported type %q", cfg.Name, cfg.Type)
	}
}
