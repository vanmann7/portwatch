// Package presencetracker records how long each port has been continuously
// open, allowing downstream components to annotate events with an uptime
// duration and to detect ports that appear only transiently.
package presencetracker

import (
	"sync"
	"time"
)

// Tracker records the first-seen timestamp for each open port.
type Tracker struct {
	mu    sync.Mutex
	first map[int]time.Time
	now   func() time.Time
}

// New returns an initialised Tracker. If now is nil the real clock is used.
func New(now func() time.Time) *Tracker {
	if now == nil {
		now = time.Now
	}
	return &Tracker{
		first: make(map[int]time.Time),
		now:   now,
	}
}

// Opened records port as open starting at the current time.
// If the port is already tracked the call is a no-op.
func (t *Tracker) Opened(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.first[port]; !ok {
		t.first[port] = t.now()
	}
}

// Closed removes port from the tracker.
func (t *Tracker) Closed(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.first, port)
}

// Uptime returns how long port has been continuously open and true.
// If the port is not tracked it returns 0 and false.
func (t *Tracker) Uptime(port int) (time.Duration, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	ts, ok := t.first[port]
	if !ok {
		return 0, false
	}
	return t.now().Sub(ts), true
}

// Snapshot returns a copy of the current port → first-seen map.
func (t *Tracker) Snapshot() map[int]time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make(map[int]time.Time, len(t.first))
	for k, v := range t.first {
		out[k] = v
	}
	return out
}

// Reset clears all tracked ports.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.first = make(map[int]time.Time)
}
