// Package escalation promotes a port-change event to a higher severity
// level when the same port has triggered repeated alerts within a sliding
// time window. This prevents alert fatigue for stable infrastructure while
// still surfacing genuinely persistent anomalies.
package escalation

import (
	"sync"
	"time"
)

// Level represents an escalation tier.
type Level int

const (
	// LevelNormal is the default level assigned to a new event.
	LevelNormal Level = iota
	// LevelElevated is assigned after repeated triggers within the window.
	LevelElevated
	// LevelCritical is assigned when the trigger count exceeds the critical threshold.
	LevelCritical
)

// String returns a human-readable representation of the level.
func (l Level) String() string {
	switch l {
	case LevelElevated:
		return "elevated"
	case LevelCritical:
		return "critical"
	default:
		return "normal"
	}
}

// Config holds the thresholds used to determine escalation.
type Config struct {
	// Window is the duration over which trigger counts are accumulated.
	Window time.Duration
	// ElevatedAfter is the number of triggers that promotes an event to LevelElevated.
	ElevatedAfter int
	// CriticalAfter is the number of triggers that promotes an event to LevelCritical.
	CriticalAfter int
}

type entry struct {
	times []time.Time
}

// Escalator tracks per-port trigger counts and returns the appropriate Level.
type Escalator struct {
	cfg     Config
	mu      sync.Mutex
	entries map[int]*entry
}

// New creates an Escalator with the given Config.
func New(cfg Config) *Escalator {
	if cfg.Window <= 0 {
		cfg.Window = time.Minute
	}
	if cfg.ElevatedAfter <= 0 {
		cfg.ElevatedAfter = 3
	}
	if cfg.CriticalAfter <= 0 {
		cfg.CriticalAfter = 7
	}
	return &Escalator{
		cfg:     cfg,
		entries: make(map[int]*entry),
	}
}

// Record registers a new trigger for port and returns the resulting Level.
func (e *Escalator) Record(port int) Level {
	e.mu.Lock()
	defer e.mu.Unlock()

	now := time.Now()
	en, ok := e.entries[port]
	if !ok {
		en = &entry{}
		e.entries[port] = en
	}

	// Evict observations outside the window.
	cutoff := now.Add(-e.cfg.Window)
	valid := en.times[:0]
	for _, t := range en.times {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	valid = append(valid, now)
	en.times = valid

	count := len(en.times)
	switch {
	case count >= e.cfg.CriticalAfter:
		return LevelCritical
	case count >= e.cfg.ElevatedAfter:
		return LevelElevated
	default:
		return LevelNormal
	}
}

// Reset clears the trigger history for port.
func (e *Escalator) Reset(port int) {
	e.mu.Lock()
	delete(e.entries, port)
	e.mu.Unlock()
}
