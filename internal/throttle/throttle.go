// Package throttle provides a simple time-based throttle that limits
// how frequently an action can be triggered within a sliding window.
package throttle

import (
	"sync"
	"time"
)

// Throttle limits calls to at most N times per window duration.
type Throttle struct {
	mu       sync.Mutex
	max      int
	window   time.Duration
	timestamps []time.Time
	now      func() time.Time
}

// New creates a Throttle that allows at most max calls per window.
func New(max int, window time.Duration) *Throttle {
	return &Throttle{
		max:    max,
		window: window,
		now:    time.Now,
	}
}

// Allow reports whether the current call is within the allowed rate.
// It records the call if allowed.
func (t *Throttle) Allow() bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	cutoff := now.Add(-t.window)

	// evict old timestamps
	valid := t.timestamps[:0]
	for _, ts := range t.timestamps {
		if ts.After(cutoff) {
			valid = append(valid, ts)
		}
	}
	t.timestamps = valid

	if len(t.timestamps) >= t.max {
		return false
	}

	t.timestamps = append(t.timestamps, now)
	return true
}

// Reset clears all recorded timestamps.
func (t *Throttle) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.timestamps = t.timestamps[:0]
}

// Count returns the number of calls recorded within the current window.
func (t *Throttle) Count() int {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	cutoff := now.Add(-t.window)
	count := 0
	for _, ts := range t.timestamps {
		if ts.After(cutoff) {
			count++
		}
	}
	return count
}
