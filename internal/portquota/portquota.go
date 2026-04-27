// Package portquota enforces per-port alert quotas over a rolling time window.
// When a port exceeds its quota the Tracker returns false, suppressing further
// alerts until the window resets.
package portquota

import (
	"sync"
	"time"
)

// entry tracks how many alerts have been issued for a single port within the
// current window.
type entry struct {
	count     int
	windowEnd time.Time
}

// Tracker enforces a maximum number of alerts per port per time window.
type Tracker struct {
	mu      sync.Mutex
	entries map[int]*entry
	max     int
	window  time.Duration
	now     func() time.Time
}

// New returns a Tracker that allows at most max alerts per port within window.
func New(max int, window time.Duration) *Tracker {
	return newWithClock(max, window, time.Now)
}

func newWithClock(max int, window time.Duration, now func() time.Time) *Tracker {
	return &Tracker{
		entries: make(map[int]*entry),
		max:     max,
		window:  window,
		now:     now,
	}
}

// Allow reports whether an alert for port should be emitted. It increments the
// counter for the port and returns false once the quota is exhausted.
func (t *Tracker) Allow(port int) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	e, ok := t.entries[port]
	if !ok || now.After(e.windowEnd) {
		t.entries[port] = &entry{count: 1, windowEnd: now.Add(t.window)}
		return true
	}
	if e.count >= t.max {
		return false
	}
	e.count++
	return true
}

// Reset clears the quota state for port, allowing alerts to flow immediately.
func (t *Tracker) Reset(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.entries, port)
}

// ResetAll clears all quota state.
func (t *Tracker) ResetAll() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = make(map[int]*entry)
}

// Remaining returns how many alerts are still permitted for port in the current
// window. It returns max when no alerts have been issued yet.
func (t *Tracker) Remaining(port int) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	e, ok := t.entries[port]
	if !ok || now.After(e.windowEnd) {
		return t.max
	}
	rem := t.max - e.count
	if rem < 0 {
		return 0
	}
	return rem
}
