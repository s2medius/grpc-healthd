package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// ProbeType enumerates supported probe kinds.
type ProbeType string

const (
	ProbeTCP  ProbeType = "tcp"
	ProbeHTTP ProbeType = "http"
	ProbeDNS  ProbeType = "dns"
	ProbeExec ProbeType = "exec"
	ProbeGRPC ProbeType = "grpc"
)

// ProbeConfig holds configuration for a single probe.
type ProbeConfig struct {
	Name     string        `yaml:"name"`
	Type     ProbeType     `yaml:"type"`
	Target   string        `yaml:"target"`
	Service  string        `yaml:"service"`   // gRPC service name
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
	Args     []string      `yaml:"args"`
}

// Config is the top-level daemon configuration.
type Config struct {
	GRPCListenAddr  string        `yaml:"grpc_listen_addr"`
	MetricsAddr     string        `yaml:"metrics_addr"`
	DefaultInterval time.Duration `yaml:"default_interval"`
	DefaultTimeout  time.Duration `yaml:"default_timeout"`
	Probes          []ProbeConfig `yaml:"probes"`
}

// DefaultConfig returns a Config populated with sensible defaults.
func DefaultConfig() Config {
	return Config{
		GRPCListenAddr:  ":50051",
		MetricsAddr:     ":9090",
		DefaultInterval: 15 * time.Second,
		DefaultTimeout:  5 * time.Second,
	}
}

// Load reads a YAML config file and merges it with defaults.
func Load(path string) (Config, error) {
	cfg := DefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	for i := range cfg.Probes {
		if cfg.Probes[i].Interval == 0 {
			cfg.Probes[i].Interval = cfg.DefaultInterval
		}
		if cfg.Probes[i].Timeout == 0 {
			cfg.Probes[i].Timeout = cfg.DefaultTimeout
		}
	}
	return cfg, nil
}
