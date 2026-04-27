// Package portstate tracks the current open/closed state of individual ports
// and emits transition events when a port changes state.
package portstate

import (
	"sync"
	"time"
)

// State represents whether a port is open or closed.
type State int

const (
	Unknown State = iota
	Open
	Closed
)

// Transition describes a port changing from one state to another.
type Transition struct {
	Port      int
	Prev      State
	Next      State
	ChangedAt time.Time
}

// Tracker maintains the last-known state for each port.
type Tracker struct {
	mu    sync.RWMutex
	states map[int]State
	clock  func() time.Time
}

// New returns a new Tracker.
func New() *Tracker {
	return newWithClock(time.Now)
}

func newWithClock(clock func() time.Time) *Tracker {
	return &Tracker{
		states: make(map[int]State),
		clock:  clock,
	}
}

// Update sets the state for port and returns a Transition if the state changed.
// If the port was not previously tracked, Prev is Unknown.
func (t *Tracker) Update(port int, next State) (Transition, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	prev, ok := t.states[port]
	if !ok {
		prev = Unknown
	}
	if prev == next {
		return Transition{}, false
	}
	t.states[port] = next
	return Transition{
		Port:      port,
		Prev:      prev,
		Next:      next,
		ChangedAt: t.clock(),
	}, true
}

// Get returns the current state of a port.
func (t *Tracker) Get(port int) State {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.states[port]
}

// Snapshot returns a copy of all tracked port states.
func (t *Tracker) Snapshot() map[int]State {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make(map[int]State, len(t.states))
	for p, s := range t.states {
		out[p] = s
	}
	return out
}

// Reset clears all tracked state.
func (t *Tracker) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.states = make(map[int]State)
}
