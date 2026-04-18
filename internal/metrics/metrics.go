// Package metrics tracks runtime statistics for portwatch scans.
package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time view of scan metrics.
type Snapshot struct {
	TotalScans   int
	OpenedPorts  int
	ClosedPorts  int
	LastScanAt   time.Time
	LastDuration time.Duration
}

// Tracker accumulates metrics across scan cycles.
type Tracker struct {
	mu           sync.Mutex
	totalScans   int
	openedPorts  int
	closedPorts  int
	lastScanAt   time.Time
	lastDuration time.Duration
}

// New returns a new Tracker.
func New() *Tracker {
	return &Tracker{}
}

// RecordScan records the result of a single scan cycle.
func (t *Tracker) RecordScan(opened, closed int, duration time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.totalScans++
	t.openedPorts += opened
	t.closedPorts += closed
	t.lastScanAt = time.Now()
	t.lastDuration = duration
}

// Snapshot returns a copy of the current metrics.
func (t *Tracker) Snapshot() Snapshot {
	t.mu.Lock()
	defer t.mu.Unlock()
	return Snapshot{
		TotalScans:   t.totalScans,
		OpenedPorts:  t.openedPorts,
		ClosedPorts:  t.closedPorts,
		LastScanAt:   t.lastScanAt,
		LastDuration: t.lastDuration,
	}
}

// Reset clears all accumulated metrics.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	*t = Tracker{}
}
