// Package portpulse tracks how frequently each port is seen open
// across successive scans, producing a "pulse" score per port.
package portpulse

import (
	"sync"
	"time"
)

// Entry holds pulse data for a single port.
type Entry struct {
	Port      int
	Hits      int
	FirstSeen time.Time
	LastSeen  time.Time
}

// Tracker accumulates hit counts per port.
type Tracker struct {
	mu      sync.Mutex
	entries map[int]*Entry
	clock   func() time.Time
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{
		entries: make(map[int]*Entry),
		clock:   time.Now,
	}
}

// Record registers one observation of the given port.
func (t *Tracker) Record(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.clock()
	e, ok := t.entries[port]
	if !ok {
		t.entries[port] = &Entry{
			Port:      port,
			Hits:      1,
			FirstSeen: now,
			LastSeen:  now,
		}
		return
	}
	e.Hits++
	e.LastSeen = now
}

// Get returns the Entry for a port and whether it exists.
func (t *Tracker) Get(port int) (Entry, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	e, ok := t.entries[port]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// Snapshot returns a copy of all tracked entries.
func (t *Tracker) Snapshot() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Entry, 0, len(t.entries))
	for _, e := range t.entries {
		out = append(out, *e)
	}
	return out
}

// Reset clears all recorded data.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = make(map[int]*Entry)
}
