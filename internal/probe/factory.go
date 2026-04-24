package probe

import (
	"fmt"

	"github.com/grpc-healthd/internal/config"
)

// Probe is the interface implemented by all health probes.
type Probe interface {
	Check() Result
}

// FromConfig constructs the appropriate Probe from a ProbeConfig.
func FromConfig(cfg config.ProbeConfig) (Probe, error) {
	switch cfg.Type {
	case "tcp":
		return NewTCPProbe(cfg.Address, cfg.Timeout), nil
	case "http", "https":
		return NewHTTPProbe(cfg.Address, cfg.Timeout), nil
	case "dns":
		return NewDNSProbe(cfg.Address, cfg.Timeout), nil
	case "exec":
		return NewExecProbe(cfg.Command, cfg.Args, cfg.Timeout), nil
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
	case "elasticsearch":
		return NewElasticsearchProbe(cfg.Address, cfg.Timeout), nil
	case "etcd":
		return NewEtcdProbe(cfg.Address, cfg.Timeout), nil
	case "nats":
		return NewNATSProbe(cfg.Address, cfg.Timeout), nil
	case "memcached":
		return NewMemcachedProbe(cfg.Address, cfg.Timeout), nil
	case "consul":
		return NewConsulProbe(cfg.Address, cfg.Timeout), nil
	case "amqp":
		return NewAMQPProbe(cfg.Address, cfg.Timeout), nil
	case "smtp":
		return NewSMTPProbe(cfg.Address, cfg.Timeout), nil
	case "http2":
		return NewHTTP2Probe(cfg.Address, cfg.Timeout), nil
	case "websocket":
		return NewWebSocketProbe(cfg.Address, cfg.Timeout), nil
	case "ftp":
		return NewFTPProbe(cfg.Address, cfg.Timeout), nil
	default:
		return nil, fmt.Errorf("unsupported probe type: %q", cfg.Type)
	}
}
