// Package portexpiry tracks ports that have been open beyond a configurable
// maximum duration and emits expiry notifications when the threshold is crossed.
package portexpiry

import (
	"sync"
	"time"
)

// Entry holds the metadata for a tracked port.
type Entry struct {
	Port      int
	OpenedAt  time.Time
	ExpiresAt time.Time
	Expired   bool
}

// Tracker monitors port open durations and marks ports as expired once they
// exceed the configured maximum lifetime.
type Tracker struct {
	mu      sync.Mutex
	entries map[int]*Entry
	max     time.Duration
	clock   func() time.Time
}

// New returns a Tracker that expires ports open longer than maxAge.
func New(maxAge time.Duration) *Tracker {
	return newWithClock(maxAge, time.Now)
}

func newWithClock(maxAge time.Duration, clock func() time.Time) *Tracker {
	return &Tracker{
		entries: make(map[int]*Entry),
		max:     maxAge,
		clock:   clock,
	}
}

// Record registers a port as open. If the port is already tracked the call is
// a no-op so the original open time is preserved.
func (t *Tracker) Record(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.entries[port]; ok {
		return
	}
	now := t.clock()
	t.entries[port] = &Entry{
		Port:      port,
		OpenedAt:  now,
		ExpiresAt: now.Add(t.max),
	}
}

// Remove stops tracking a port.
func (t *Tracker) Remove(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, port)
}

// Expired returns all ports whose open duration has exceeded the maximum age.
// Each returned entry has its Expired field set to true.
func (t *Tracker) Expired() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.clock()
	var out []Entry
	for _, e := range t.entries {
		if now.After(e.ExpiresAt) {
			e.Expired = true
			out = append(out, *e)
		}
	}
	return out
}

// Snapshot returns a copy of all currently tracked entries regardless of
// expiry status.
func (t *Tracker) Snapshot() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Entry, 0, len(t.entries))
	now := t.clock()
	for _, e := range t.entries {
		copy := *e
		copy.Expired = now.After(e.ExpiresAt)
		out = append(out, copy)
	}
	return out
}
