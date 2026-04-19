package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// GRPCConfig holds gRPC server settings.
type GRPCConfig struct {
	Addr string `yaml:"addr"`
}

// MetricsConfig holds Prometheus metrics server settings.
type MetricsConfig struct {
	Addr string `yaml:"addr"`
}

// Config is the top-level configuration structure.
type Config struct {
	GRPC    GRPCConfig    `yaml:"grpc"`
	Metrics MetricsConfig `yaml:"metrics"`
	Probes  []ProbeConfig `yaml:"probes"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		GRPC:    GRPCConfig{Addr: ":50051"},
		Metrics: MetricsConfig{Addr: ":9090"},
	}
}

// Load reads a YAML config file and merges it with defaults.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	for i := range cfg.Probes {
		if cfg.Probes[i].Interval == 0 {
			cfg.Probes[i].Interval = 30 * time.Second
		}
		if cfg.Probes[i].Timeout == 0 {
			cfg.Probes[i].Timeout = 5 * time.Second
		}
	}

	return cfg, nil
}
