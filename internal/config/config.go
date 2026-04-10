package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Target represents a single port/host to monitor.
type Target struct {
	Name     string        `yaml:"name"`
	Host     string        `yaml:"host"`
	Port     int           `yaml:"port"`
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
}

// Config holds the full portwatch configuration.
type Config struct {
	Targets []Target `yaml:"targets"`
}

// Validate checks that all targets have required fields and sane values.
func (c *Config) Validate() error {
	if len(c.Targets) == 0 {
		return fmt.Errorf("config: no targets defined")
	}
	for i, t := range c.Targets {
		if t.Host == "" {
			return fmt.Errorf("config: target[%d] missing host", i)
		}
		if t.Port < 1 || t.Port > 65535 {
			return fmt.Errorf("config: target[%d] port %d out of range", i, t.Port)
		}
		if t.Interval <= 0 {
			c.Targets[i].Interval = 30 * time.Second
		}
		if t.Timeout <= 0 {
			c.Targets[i].Timeout = 5 * time.Second
		}
		if t.Name == "" {
			c.Targets[i].Name = fmt.Sprintf("%s:%d", t.Host, t.Port)
		}
	}
	return nil
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: reading file: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config: parsing yaml: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}
