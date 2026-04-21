package escalation

import (
	"testing"
	"time"
)

func defaultCfg() Config {
	return Config{
		Window:        200 * time.Millisecond,
		ElevatedAfter: 3,
		CriticalAfter: 5,
	}
}

func TestNormalOnFirstTrigger(t *testing.T) {
	e := New(defaultCfg())
	if got := e.Record(80); got != LevelNormal {
		t.Fatalf("expected Normal, got %s", got)
	}
}

func TestElevatedAfterThreshold(t *testing.T) {
	e := New(defaultCfg())
	var got Level
	for i := 0; i < 3; i++ {
		got = e.Record(443)
	}
	if got != LevelElevated {
		t.Fatalf("expected Elevated after 3 triggers, got %s", got)
	}
}

func TestCriticalAfterThreshold(t *testing.T) {
	e := New(defaultCfg())
	var got Level
	for i := 0; i < 5; i++ {
		got = e.Record(22)
	}
	if got != LevelCritical {
		t.Fatalf("expected Critical after 5 triggers, got %s", got)
	}
}

func TestWindowEvictsOldTriggers(t *testing.T) {
	cfg := Config{
		Window:        50 * time.Millisecond,
		ElevatedAfter: 3,
		CriticalAfter: 5,
	}
	e := New(cfg)
	for i := 0; i < 4; i++ {
		e.Record(8080)
	}
	// Wait for the window to expire.
	time.Sleep(80 * time.Millisecond)
	got := e.Record(8080)
	if got != LevelNormal {
		t.Fatalf("expected Normal after window expiry, got %s", got)
	}
}

func TestResetClearsHistory(t *testing.T) {
	e := New(defaultCfg())
	for i := 0; i < 5; i++ {
		e.Record(9000)
	}
	e.Reset(9000)
	if got := e.Record(9000); got != LevelNormal {
		t.Fatalf("expected Normal after reset, got %s", got)
	}
}

func TestLevelString(t *testing.T) {
	cases := []struct {
		level Level
		want  string
	}{
		{LevelNormal, "normal"},
		{LevelElevated, "elevated"},
		{LevelCritical, "critical"},
	}
	for _, tc := range cases {
		if got := tc.level.String(); got != tc.want {
			t.Errorf("Level(%d).String() = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestDefaultsAppliedWhenZero(t *testing.T) {
	e := New(Config{}) // all zero values
	if e.cfg.Window != time.Minute {
		t.Errorf("expected default window 1m, got %s", e.cfg.Window)
	}
	if e.cfg.ElevatedAfter != 3 {
		t.Errorf("expected ElevatedAfter=3, got %d", e.cfg.ElevatedAfter)
	}
	if e.cfg.CriticalAfter != 7 {
		t.Errorf("expected CriticalAfter=7, got %d", e.cfg.CriticalAfter)
	}
}
