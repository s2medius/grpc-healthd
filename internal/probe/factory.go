package probe

import (
	"fmt"

	"github.com/yourusername/grpc-healthd/internal/config"
)

// Probe is the interface implemented by all probe types.
type Probe interface {
	Run(ctx interface{ Done() <-chan struct{} }) Status
}

// FromConfig constructs a Probe from the given ProbeConfig.
// Returns an error if the probe type is unknown.
func FromConfig(cfg config.ProbeConfig) (Probe, error) {
	switch cfg.Type {
	case "tcp":
		return NewTCPProbe(cfg), nil
	case "http", "https":
		return NewHTTPProbe(cfg), nil
	case "dns":
		return NewDNSProbe(cfg), nil
	case "exec":
		return NewExecProbe(cfg), nil
	case "grpc":
		return NewGRPCProbe(cfg), nil
	case "tls":
		return NewTLSProbe(cfg), nil
	case "icmp":
		return NewICMPProbe(cfg), nil
	case "redis":
		return NewRedisProbe(cfg), nil
	default:
		return nil, fmt.Errorf("unknown probe type: %q", cfg.Type)
	}
}
