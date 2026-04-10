package config

import (
	"os"
	"testing"
	"time"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "portwatch-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoad_Valid(t *testing.T) {
	yaml := `
targets:
  - name: local-http
    host: localhost
    port: 8080
    interval: 10s
    timeout: 3s
`
	path := writeTemp(t, yaml)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Targets) != 1 {
		t.Fatalf("expected 1 target, got %d", len(cfg.Targets))
	}
	if cfg.Targets[0].Interval != 10*time.Second {
		t.Errorf("expected interval 10s, got %v", cfg.Targets[0].Interval)
	}
}

func TestLoad_Defaults(t *testing.T) {
	yaml := `
targets:
  - host: example.com
    port: 443
`
	path := writeTemp(t, yaml)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tgt := cfg.Targets[0]
	if tgt.Interval != 30*time.Second {
		t.Errorf("expected default interval 30s, got %v", tgt.Interval)
	}
	if tgt.Timeout != 5*time.Second {
		t.Errorf("expected default timeout 5s, got %v", tgt.Timeout)
	}
	if tgt.Name != "example.com:443" {
		t.Errorf("unexpected default name: %s", tgt.Name)
	}
}

func TestLoad_NoTargets(t *testing.T) {
	path := writeTemp(t, "targets: []\n")
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for empty targets")
	}
}

func TestLoad_InvalidPort(t *testing.T) {
	yaml := `
targets:
  - host: localhost
    port: 99999
`
	path := writeTemp(t, yaml)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid port")
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/portwatch.yaml")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
