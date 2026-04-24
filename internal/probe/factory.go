package probe

import (
	"fmt"
	"time"

	"github.com/yourorg/grpc-healthd/internal/config"
)

const defaultTimeoutStr = "5s"

// FromConfig constructs a Probe from a ProbeConfig.
func FromConfig(cfg config.ProbeConfig) (interface {
	Probe(ctx interface{}) Result
}, error) {
	var timeout time.Duration
	if cfg.Timeout != "" {
		var err error
		timeout, err = time.ParseDuration(cfg.Timeout)
		if err != nil {
			return nil, fmt.Errorf("invalid timeout %q: %w", cfg.Timeout, err)
		}
	}

	switch cfg.Type {
	case "tcp":
		return NewTCPProbe(cfg.Address, timeout), nil
	case "http":
		return NewHTTPProbe(cfg.Address, timeout), nil
	case "dns":
		return NewDNSProbe(cfg.Address, timeout), nil
	case "exec":
		return NewExecProbe(cfg.Command, cfg.Args, timeout), nil
	case "grpc":
		return NewGRPCProbe(cfg.Address, timeout), nil
	case "tls":
		return NewTLSProbe(cfg.Address, timeout), nil
	case "icmp":
		return NewICMPProbe(cfg.Address, timeout), nil
	case "redis":
		return NewRedisProbe(cfg.Address, timeout), nil
	case "postgres":
		return NewPostgresProbe(cfg.Address, timeout), nil
	case "mysql":
		return NewMySQLProbe(cfg.Address, timeout), nil
	case "mongodb":
		return NewMongoDBProbe(cfg.Address, timeout), nil
	case "kafka":
		return NewKafkaProbe(cfg.Address, timeout), nil
	case "rabbitmq":
		return NewRabbitMQProbe(cfg.Address, timeout), nil
	case "elasticsearch":
		return NewElasticsearchProbe(cfg.Address, timeout), nil
	case "etcd":
		return NewEtcdProbe(cfg.Address, timeout), nil
	case "nats":
		return NewNATSProbe(cfg.Address, timeout), nil
	case "memcached":
		return NewMemcachedProbe(cfg.Address, timeout), nil
	case "consul":
		return NewConsulProbe(cfg.Address, timeout), nil
	case "amqp":
		return NewAMQPProbe(cfg.Address, timeout), nil
	case "smtp":
		return NewSMTPProbe(cfg.Address, timeout), nil
	case "http2":
		return NewHTTP2Probe(cfg.Address, timeout), nil
	case "websocket":
		return NewWebSocketProbe(cfg.Address, timeout), nil
	default:
		return nil, fmt.Errorf("unknown probe type: %q", cfg.Type)
	}
}
