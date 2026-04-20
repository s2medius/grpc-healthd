package config

import "fmt"

// ProbeType enumerates the supported probe kinds.
type ProbeType string

const (
	ProbeTypeTCP      ProbeType = "tcp"
	ProbeTypeHTTP     ProbeType = "http"
	ProbeTypeDNS      ProbeType = "dns"
	ProbeTypeExec     ProbeType = "exec"
	ProbeTypeGRPC     ProbeType = "grpc"
	ProbeTypeTLS      ProbeType = "tls"
	ProbeTypeICMP     ProbeType = "icmp"
	ProbeTypeRedis    ProbeType = "redis"
	ProbeTypePostgres ProbeType = "postgres"
	ProbeTypeMySQL    ProbeType = "mysql"
)

// ProbeConfig holds configuration for a single probe.
type ProbeConfig struct {
	Name     string        `yaml:"name"`
	Type     ProbeType     `yaml:"type"`
	Address  string        `yaml:"address"`
	Interval string        `yaml:"interval"`
	Timeout  string        `yaml:"timeout"`
	Command  []string      `yaml:"command,omitempty"`
}

// Validate returns an error if the ProbeConfig is missing required fields.
func (p ProbeConfig) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("probe name is required")
	}
	switch p.Type {
	case ProbeTypeTCP, ProbeTypeHTTP, ProbeTypeDNS, ProbeTypeGRPC,
		ProbeTypeTLS, ProbeTypeICMP, ProbeTypeRedis, ProbeTypePostgres, ProbeTypeMySQL:
		if p.Address == "" {
			return fmt.Errorf("probe %q: address is required for type %q", p.Name, p.Type)
		}
	case ProbeTypeExec:
		if len(p.Command) == 0 {
			return fmt.Errorf("probe %q: command is required for type exec", p.Name)
		}
	case "":
		return fmt.Errorf("probe %q: type is required", p.Name)
	default:
		return fmt.Errorf("probe %q: unknown type %q", p.Name, p.Type)
	}
	return nil
}
