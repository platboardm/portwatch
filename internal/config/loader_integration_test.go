package config_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

// TestLoadExampleConfig ensures the bundled example config parses correctly.
func TestLoadExampleConfig(t *testing.T) {
	cfg, err := config.Load("example_config.yaml")
	if err != nil {
		t.Fatalf("failed to load example config: %v", err)
	}
	if len(cfg.Targets) == 0 {
		t.Fatal("example config should contain at least one target")
	}
	for _, tgt := range cfg.Targets {
		if tgt.Host == "" {
			t.Errorf("target %q has empty host", tgt.Name)
		}
		if tgt.Port < 1 || tgt.Port > 65535 {
			t.Errorf("target %q has invalid port %d", tgt.Name, tgt.Port)
		}
		if tgt.Interval <= 0 {
			t.Errorf("target %q has non-positive interval", tgt.Name)
		}
		if tgt.Timeout <= 0 {
			t.Errorf("target %q has non-positive timeout", tgt.Name)
		}
	}
}
