package config_test

import (
	"testing"

	"github.com/your-org/grpc-healthd/internal/config"
)

func TestProbeConfig_Validate_Valid(t *testing.T) {
	cases := []config.ProbeConfig{
		{Name: "tcp-check", Type: config.ProbeTCP, Address: "localhost:80"},
		{Name: "http-check", Type: config.ProbeHTTP, Address: "http://example.com"},
		{Name: "dns-check", Type: config.ProbeDNS, Address: "example.com"},
		{Name: "exec-check", Type: config.ProbeExec, Command: "/bin/true"},
		{Name: "grpc-check", Type: config.ProbeGRPC, Address: "localhost:50051"},
		{Name: "tls-check", Type: config.ProbeTLS, Address: "localhost:443"},
	}
	for _, c := range cases {
		if msg := c.Validate(); msg != "" {
			t.Errorf("%s: unexpected error: %s", c.Name, msg)
		}
	}
}

func TestProbeConfig_Validate_MissingName(t *testing.T) {
	p := config.ProbeConfig{Type: config.ProbeTCP, Address: "localhost:80"}
	if msg := p.Validate(); msg == "" {
		t.Fatal("expected validation error for missing name")
	}
}

func TestProbeConfig_Validate_MissingAddress(t *testing.T) {
	for _, pt := range []config.ProbeType{config.ProbeTCP, config.ProbeHTTP, config.ProbeDNS, config.ProbeGRPC, config.ProbeTLS} {
		p := config.ProbeConfig{Name: "x", Type: pt}
		if msg := p.Validate(); msg == "" {
			t.Errorf("type %s: expected error for missing address", pt)
		}
	}
}

func TestProbeConfig_Validate_ExecMissingCommand(t *testing.T) {
	p := config.ProbeConfig{Name: "e", Type: config.ProbeExec}
	if msg := p.Validate(); msg == "" {
		t.Fatal("expected error for exec with no command")
	}
}

func TestProbeConfig_Validate_UnknownType(t *testing.T) {
	p := config.ProbeConfig{Name: "x", Type: "ftp", Address: "localhost"}
	if msg := p.Validate(); msg == "" {
		t.Fatal("expected error for unknown probe type")
	}
}
