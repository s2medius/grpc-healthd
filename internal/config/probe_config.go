package config

// ProbeType enumerates supported probe kinds.
type ProbeType string

const (
	ProbeTCP  ProbeType = "tcp"
	ProbeHTTP ProbeType = "http"
	ProbeDNS  ProbeType = "dns"
	ProbeExec ProbeType = "exec"
	ProbeGRPC ProbeType = "grpc"
	ProbeTLS  ProbeType = "tls"
)

// ProbeConfig holds configuration for a single probe.
type ProbeConfig struct {
	Name       string        `toml:"name"`
	Type       ProbeType     `toml:"type"`
	Address    string        `toml:"address"`
	Command    string        `toml:"command"`
	Args       []string      `toml:"args"`
	TimeoutSec int           `toml:"timeout_sec"`
	IntervalSec int          `toml:"interval_sec"`
	SkipTLSVerify bool       `toml:"skip_tls_verify"`
}

// Validate returns an error string if the ProbeConfig is invalid, empty string otherwise.
func (p ProbeConfig) Validate() string {
	if p.Name == "" {
		return "probe name must not be empty"
	}
	switch p.Type {
	case ProbeTCP, ProbeHTTP, ProbeDNS, ProbeGRPC, ProbeTLS:
		if p.Address == "" {
			return "probe address must not be empty for type " + string(p.Type)
		}
	case ProbeExec:
		if p.Command == "" {
			return "probe command must not be empty for exec type"
		}
	default:
		return "unknown probe type: " + string(p.Type)
	}
	return ""
}
