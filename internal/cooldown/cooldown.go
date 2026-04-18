// Package cooldown provides a per-port cooldown tracker that prevents
// repeated alerts from firing within a configurable quiet period.
package cooldown

import (
	"sync"
	"time"
)

// Tracker records the last alert time for each port and decides whether
// a new alert is allowed based on a minimum interval.
type Tracker struct {
	mu       sync.Mutex
	last     map[int]time.Time
	interval time.Duration
	now      func() time.Time
}

// New returns a Tracker with the given cooldown interval.
func New(interval time.Duration) *Tracker {
	return &Tracker{
		last:     make(map[int]time.Time),
		interval: interval,
		now:      time.Now,
	}
}

// Allow returns true if the port has not been alerted within the cooldown
// interval. Calling Allow records the current time for the port when it
// returns true.
func (t *Tracker) Allow(port int) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	if last, ok := t.last[port]; ok {
		if t.now().Sub(last) < t.interval {
			return false
		}
	}
	t.last[port] = t.now()
	return true
}

// Reset clears the cooldown record for a specific port.
func (t *Tracker) Reset(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, port)
}

// ResetAll clears all cooldown records.
func (t *Tracker) ResetAll() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.last = make(map[int]time.Time)
}

// Count returns the number of ports currently tracked.
func (t *Tracker) Count() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.last)
}
