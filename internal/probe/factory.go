package probe

import (
	"fmt"
	"time"

	"github.com/yourorg/grpc-healthd/internal/config"
)

// Probe is the interface implemented by all probe types.
type Prober interface {
	Probe() Result
}

// FromConfig constructs a Prober from a ProbeConfig.
func FromConfig(cfg config.ProbeConfig) (Prober, error) {
	timeout := time.Duration(cfg.TimeoutSeconds) * time.Second
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	switch cfg.Type {
	case "tcp":
		return NewTCPProbe(cfg.Address, timeout), nil
	case "http":
		return NewHTTPProbe(cfg.Address, timeout), nil
	case "dns":
		return NewDNSProbe(cfg.Address, timeout), nil
	case "exec":
		if len(cfg.Command) == 0 {
			return nil, fmt.Errorf("exec probe %q: command must not be empty", cfg.Name)
		}
		return NewExecProbe(cfg.Command[0], cfg.Command[1:], timeout), nil
	case "grpc":
		return NewGRPCProbe(cfg.Address, cfg.Service, timeout), nil
	case "tls":
		return NewTLSProbe(cfg.Address, timeout), nil
	default:
		return nil, fmt.Errorf("unknown probe type %q for probe %q", cfg.Type, cfg.Name)
	}
}
