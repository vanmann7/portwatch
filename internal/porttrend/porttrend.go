// Package porttrend tracks the rate of change for observed ports over time,
// providing a simple rising/falling/stable trend indicator.
package porttrend

import (
	"sync"
	"time"
)

// Trend represents the direction of change for a port's activity.
type Trend int

const (
	Stable  Trend = iota // no net change in the observation window
	Rising               // more opens than closes
	Falling              // more closes than opens
)

func (t Trend) String() string {
	switch t {
	case Rising:
		return "rising"
	case Falling:
		return "falling"
	default:
		return "stable"
	}
}

type event struct {
	at     time.Time
	opened bool
}

// Tracker records open/close events per port and reports trend.
type Tracker struct {
	mu     sync.Mutex
	events map[int][]event
	window time.Duration
	now    func() time.Time
}

// New returns a Tracker that considers events within the given window.
func New(window time.Duration) *Tracker {
	return &Tracker{
		events: make(map[int][]event),
		window: window,
		now:    time.Now,
	}
}

// RecordOpen records an open event for the given port.
func (t *Tracker) RecordOpen(port int) {
	t.record(port, true)
}

// RecordClose records a close event for the given port.
func (t *Tracker) RecordClose(port int) {
	t.record(port, false)
}

func (t *Tracker) record(port int, opened bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	t.evict(port, now)
	t.events[port] = append(t.events[port], event{at: now, opened: opened})
}

// Trend returns the current trend for the given port.
func (t *Tracker) Trend(port int) Trend {
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	t.evict(port, now)
	var opens, closes int
	for _, e := range t.events[port] {
		if e.opened {
			opens++
		} else {
			closes++
		}
	}
	switch {
	case opens > closes:
		return Rising
	case closes > opens:
		return Falling
	default:
		return Stable
	}
}

// Reset clears all recorded events for the given port.
func (t *Tracker) Reset(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.events, port)
}

// evict removes events outside the window; caller must hold mu.
func (t *Tracker) evict(port int, now time.Time) {
	cutoff := now.Add(-t.window)
	evs := t.events[port]
	i := 0
	for i < len(evs) && evs[i].at.Before(cutoff) {
		i++
	}
	t.events[port] = evs[i:]
}
