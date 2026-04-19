package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// ProbeConfig defines a single health probe target.
type ProbeConfig struct {
	Name     string        `yaml:"name"`
	Address  string        `yaml:"address"`
	Timeout  time.Duration `yaml:"timeout"`
	Interval time.Duration `yaml:"interval"`
	Service  string        `yaml:"service"` // gRPC health check service name (empty = overall)
}

// ServerConfig holds gRPC and metrics server settings.
type ServerConfig struct {
	GRPCListenAddr    string `yaml:"grpc_listen_addr"`
	MetricsListenAddr string `yaml:"metrics_listen_addr"`
}

// Config is the top-level configuration structure.
type Config struct {
	Server ServerConfig  `yaml:"server"`
	Probes []ProbeConfig `yaml:"probes"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			GRPCListenAddr:    ":50051",
			MetricsListenAddr: ":9090",
		},
	}
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	for i := range cfg.Probes {
		if cfg.Probes[i].Timeout == 0 {
			cfg.Probes[i].Timeout = 5 * time.Second
		}
		if cfg.Probes[i].Interval == 0 {
			cfg.Probes[i].Interval = 30 * time.Second
		}
	}

	return cfg, nil
}
