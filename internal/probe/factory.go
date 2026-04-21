package probe

import (
	"fmt"

	"github.com/yourorg/grpc-healthd/internal/config"
)

// FromConfig constructs a Probe from a ProbeConfig.
func FromConfig(cfg config.ProbeConfig) (Probe, error) {
	switch cfg.Type {
	case "tcp":
		return NewTCPProbe(cfg.Address, cfg.Timeout), nil
	case "http":
		return NewHTTPProbe(cfg.Address, cfg.Timeout), nil
	case "dns":
		return NewDNSProbe(cfg.Address, cfg.Timeout), nil
	case "exec":
		if len(cfg.Command) == 0 {
			return nil, fmt.Errorf("exec probe %q: command is required", cfg.Name)
		}
		return NewExecProbe(cfg.Command[0], cfg.Command[1:], cfg.Timeout), nil
	case "grpc":
		return NewGRPCProbe(cfg.Address, cfg.Timeout), nil
	case "tls":
		return NewTLSProbe(cfg.Address, cfg.Timeout), nil
	case "icmp":
		return NewICMPProbe(cfg.Address, cfg.Timeout), nil
	case "redis":
		return NewRedisProbe(cfg.Address, cfg.Timeout), nil
	case "postgres":
		return NewPostgresProbe(cfg.Address, cfg.Timeout), nil
	case "mysql":
		return NewMySQLProbe(cfg.Address, cfg.Timeout), nil
	case "mongodb":
		return NewMongoDBProbe(cfg.Address, cfg.Timeout), nil
	case "kafka":
		return NewKafkaProbe(cfg.Address, cfg.Timeout), nil
	case "rabbitmq":
		return NewRabbitMQProbe(cfg.Address, cfg.Timeout), nil
	default:
		return nil, fmt.Errorf("unknown probe type: %q", cfg.Type)
	}
}
