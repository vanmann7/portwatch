// Package config handles loading and validating portwatch configuration.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds the runtime configuration for portwatch.
type Config struct {
	// PortRange defines the inclusive range of ports to scan.
	PortRange PortRange `json:"port_range"`
	// Interval is how often to run a scan.
	Interval time.Duration `json:"interval"`
	// StateFile is the path where port state is persisted.
	StateFile string `json:"state_file"`
	// LogFile is an optional path to write alerts; empty means stdout.
	LogFile string `json:"log_file"`
}

// PortRange defines a start/end port boundary.
type PortRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// Default returns a sensible default configuration.
func Default() *Config {
	return &Config{
		PortRange: PortRange{Start: 1, End: 1024},
		Interval:  30 * time.Second,
		StateFile: "/tmp/portwatch.state",
	}
}

// Load reads a JSON config file from path and returns a Config.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: open %q: %w", path, err)
	}
	defer f.Close()

	cfg := Default()
	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, fmt.Errorf("config: decode: %w", err)
	}
	return cfg, cfg.Validate()
}

// Validate checks that the configuration values are sensible.
func (c *Config) Validate() error {
	if c.PortRange.Start < 1 || c.PortRange.Start > 65535 {
		return fmt.Errorf("config: port_range.start must be 1-65535, got %d", c.PortRange.Start)
	}
	if c.PortRange.End < 1 || c.PortRange.End > 65535 {
		return fmt.Errorf("config: port_range.end must be 1-65535, got %d", c.PortRange.End)
	}
	if c.PortRange.Start > c.PortRange.End {
		return fmt.Errorf("config: port_range.start (%d) must be <= end (%d)", c.PortRange.Start, c.PortRange.End)
	}
	if c.Interval < time.Second {
		return fmt.Errorf("config: interval must be >= 1s, got %s", c.Interval)
	}
	if c.StateFile == "" {
		return fmt.Errorf("config: state_file must not be empty")
	}
	return nil
}
