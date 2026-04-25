// Package portlifetime tracks how long each port has been continuously open.
package portlifetime

import (
	"sync"
	"time"
)

// Entry holds the open timestamp and last-seen timestamp for a port.
type Entry struct {
	OpenedAt  time.Time
	LastSeen  time.Time
}

// Tracker records when ports were first opened and computes their lifetime.
type Tracker struct {
	mu      sync.RWMutex
	entries map[int]Entry
	now     func() time.Time
}

// New returns a Tracker using the real clock.
func New() *Tracker {
	return &Tracker{
		entries: make(map[int]Entry),
		now:     time.Now,
	}
}

// newWithClock returns a Tracker using a custom clock (for testing).
func newWithClock(clock func() time.Time) *Tracker {
	return &Tracker{
		entries: make(map[int]Entry),
		now:     clock,
	}
}

// Open records that a port was observed open. If the port is already tracked
// its LastSeen timestamp is updated; otherwise a new entry is created.
func (t *Tracker) Open(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	if e, ok := t.entries[port]; ok {
		e.LastSeen = now
		t.entries[port] = e
		return
	}
	t.entries[port] = Entry{OpenedAt: now, LastSeen: now}
}

// Close removes a port from tracking.
func (t *Tracker) Close(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, port)
}

// Lifetime returns how long the port has been open and whether it is tracked.
func (t *Tracker) Lifetime(port int) (time.Duration, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	e, ok := t.entries[port]
	if !ok {
		return 0, false
	}
	return t.now().Sub(e.OpenedAt), true
}

// Snapshot returns a copy of all tracked entries.
func (t *Tracker) Snapshot() map[int]Entry {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make(map[int]Entry, len(t.entries))
	for k, v := range t.entries {
		out[k] = v
	}
	return out
}

// Reset clears all tracked ports.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = make(map[int]Entry)
}
