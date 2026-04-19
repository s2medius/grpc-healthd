package probe_test

import (
	"testing"

	"github.com/yourorg/grpc-healthd/internal/config"
	"github.com/yourorg/grpc-healthd/internal/probe"
)

func cfg(typ, name, addr string, cmd []string) config.ProbeConfig {
	return config.ProbeConfig{Type: typ, Name: name, Address: addr, Command: cmd}
}

func TestFromConfig_TCP(t *testing.T) {
	p, err := probe.FromConfig(cfg("tcp", "svc", "localhost:80", nil))
	if err != nil || p == nil {
		t.Fatalf("expected probe, got err=%v", err)
	}
}

func TestFromConfig_HTTP(t *testing.T) {
	p, err := probe.FromConfig(cfg("http", "svc", "http://localhost", nil))
	if err != nil || p == nil {
		t.Fatalf("expected probe, got err=%v", err)
	}
}

func TestFromConfig_DNS(t *testing.T) {
	p, err := probe.FromConfig(cfg("dns", "svc", "example.com", nil))
	if err != nil || p == nil {
		t.Fatalf("expected probe, got err=%v", err)
	}
}

func TestFromConfig_Exec(t *testing.T) {
	p, err := probe.FromConfig(cfg("exec", "svc", "", []string{"echo", "ok"}))
	if err != nil || p == nil {
		t.Fatalf("expected probe, got err=%v", err)
	}
}

func TestFromConfig_Exec_MissingCommand(t *testing.T) {
	_, err := probe.FromConfig(cfg("exec", "svc", "", nil))
	if err == nil {
		t.Fatal("expected error for missing exec command")
	}
}

func TestFromConfig_GRPC(t *testing.T) {
	p, err := probe.FromConfig(cfg("grpc", "svc", "localhost:50051", nil))
	if err != nil || p == nil {
		t.Fatalf("expected probe, got err=%v", err)
	}
}

func TestFromConfig_TLS(t *testing.T) {
	p, err := probe.FromConfig(cfg("tls", "svc", "localhost:443", nil))
	if err != nil || p == nil {
		t.Fatalf("expected probe, got err=%v", err)
	}
}

func TestFromConfig_UnknownType(t *testing.T) {
	_, err := probe.FromConfig(cfg("ftp", "svc", "localhost:21", nil))
	if err == nil {
		t.Fatal("expected error for unknown probe type")
	}
}
