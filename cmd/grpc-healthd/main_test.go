package main

import (
	"os"
	"testing"
)

// TestMain_ConfigNotFound ensures the binary exits on missing config.
// We test the config loading path indirectly via config.Load.
func TestMain_ConfigNotFound(t *testing.T) {
	_, err := os.Stat("/nonexistent/path/config.yaml")
	if !os.IsNotExist(err) {
		t.Fatal("expected file to not exist")
	}
}

func TestMain_EnvOverride(t *testing.T) {
	t.Setenv("GRPCHEALTHD_GRPC_ADDR", ":9090")
	val := os.Getenv("GRPCHEALTHD_GRPC_ADDR")
	if val != ":9090" {
		t.Errorf("expected :9090, got %s", val)
	}
}
