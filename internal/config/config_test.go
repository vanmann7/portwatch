package config_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/config"
)

func writeConfig(t *testing.T, v any) string {
	t.Helper()
	f, err := os.CreateTemp("", "portwatch-cfg-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(v); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestDefault(t *testing.T) {
	cfg := config.Default()
	if cfg.PortRange.Start != 1 || cfg.PortRange.End != 1024 {
		t.Errorf("unexpected default range: %+v", cfg.PortRange)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("unexpected default interval: %s", cfg.Interval)
	}
}

func TestLoadValid(t *testing.T) {
	path := writeConfig(t, map[string]any{
		"port_range": map[string]int{"start": 80, "end": 8080},
		"interval":   int(10 * time.Second),
		"state_file": "/tmp/test.state",
	})
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PortRange.Start != 80 || cfg.PortRange.End != 8080 {
		t.Errorf("unexpected range: %+v", cfg.PortRange)
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/portwatch.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestValidateInvalidRange(t *testing.T) {
	cfg := config.Default()
	cfg.PortRange.Start = 9000
	cfg.PortRange.End = 80
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for inverted range")
	}
}

func TestValidateShortInterval(t *testing.T) {
	cfg := config.Default()
	cfg.Interval = 500 * time.Millisecond
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for short interval")
	}
}

func TestValidateEmptyStateFile(t *testing.T) {
	cfg := config.Default()
	cfg.StateFile = ""
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for empty state_file")
	}
}
