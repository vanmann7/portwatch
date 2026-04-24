// Package portwindow provides a sliding-window tracker that records
// which ports were seen open during a configurable time window.
package portwindow

import (
	"sync"
	"time"
)

// Entry holds a port number and the time it was last observed.
type Entry struct {
	Port      int
	ObservedAt time.Time
}

// Window tracks ports seen within a rolling duration.
type Window struct {
	mu       sync.Mutex
	duration time.Duration
	entries  []Entry
	now      func() time.Time
}

// New creates a Window that retains observations within d.
func New(d time.Duration) *Window {
	return &Window{
		duration: d,
		now:      time.Now,
	}
}

// Record adds port to the window, updating its timestamp if already present.
func (w *Window) Record(port int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict()
	for i, e := range w.entries {
		if e.Port == port {
			w.entries[i].ObservedAt = w.now()
			return
		}
	}
	w.entries = append(w.entries, Entry{Port: port, ObservedAt: w.now()})
}

// Contains reports whether port has been observed within the current window.
func (w *Window) Contains(port int) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict()
	for _, e := range w.entries {
		if e.Port == port {
			return true
		}
	}
	return false
}

// Ports returns a snapshot of all ports currently in the window.
func (w *Window) Ports() []int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict()
	out := make([]int, len(w.entries))
	for i, e := range w.entries {
		out[i] = e.Port
	}
	return out
}

// Reset clears all recorded observations.
func (w *Window) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entries = w.entries[:0]
}

// evict removes entries older than the window duration. Must be called with mu held.
func (w *Window) evict() {
	cutoff := w.now().Add(-w.duration)
	filtered := w.entries[:0]
	for _, e := range w.entries {
		if e.ObservedAt.After(cutoff) {
			filtered = append(filtered, e)
		}
	}
	w.entries = filtered
}
