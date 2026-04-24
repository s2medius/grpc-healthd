package probe

import (
	"testing"
	"time"

	"github.com/grpc-healthd/internal/config"
)

func TestFromConfig_FTP(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "ftp-check",
		Type:    "ftp",
		Address: "localhost:21",
	}
	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
	ftp, ok := p.(*FTPProbe)
	if !ok {
		t.Fatalf("expected *FTPProbe, got %T", p)
	}
	if ftp.address != "localhost:21" {
		t.Errorf("expected address localhost:21, got %s", ftp.address)
	}
	if ftp.timeout != defaultTimeout {
		t.Errorf("expected default timeout, got %v", ftp.timeout)
	}
}

func TestFromConfig_FTP_WithTimeout(t *testing.T) {
	cfg := config.ProbeConfig{
		Name:    "ftp-check",
		Type:    "ftp",
		Address: "ftp.example.com:21",
		Timeout: 3 * time.Second,
	}
	p, err := FromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ftp := p.(*FTPProbe)
	if ftp.timeout != 3*time.Second {
		t.Errorf("expected 3s timeout, got %v", ftp.timeout)
	}
}
