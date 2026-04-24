package config

import (
	"fmt"
	"time"
)

// ProbeConfig holds the configuration for a single health probe.
type ProbeConfig struct {
	Name    string        `yaml:"name"`
	Type    string        `yaml:"type"`
	Address string        `yaml:"address"`
	Timeout time.Duration `yaml:"timeout"`
	Command string        `yaml:"command"`
	Args    []string      `yaml:"args"`
	Interval time.Duration `yaml:"interval"`
}

// knownProbeTypes lists all supported probe type identifiers.
var knownProbeTypes = map[string]struct{}{
	"tcp": {}, "http": {}, "https": {}, "dns": {},
	"exec": {}, "grpc": {}, "tls": {}, "icmp": {},
	"redis": {}, "postgres": {}, "mysql": {}, "mongodb": {},
	"kafka": {}, "rabbitmq": {}, "elasticsearch": {}, "etcd": {},
	"nats": {}, "memcached": {}, "consul": {}, "amqp": {},
	"smtp": {}, "http2": {}, "websocket": {}, "ftp": {},
}

// Validate checks that the ProbeConfig has all required fields.
func (p ProbeConfig) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("probe name is required")
	}
	if _, ok := knownProbeTypes[p.Type]; !ok {
		return fmt.Errorf("unknown probe type %q for probe %q", p.Type, p.Name)
	}
	if p.Type == "exec" {
		if p.Command == "" {
			return fmt.Errorf("probe %q: command is required for exec type", p.Name)
		}
		return nil
	}
	if p.Address == "" {
		return fmt.Errorf("probe %q: address is required for type %q", p.Name, p.Type)
	}
	return nil
}
