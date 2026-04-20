// Package aggregator groups port change events within a time window
// and emits a single consolidated summary, reducing alert noise when
// many ports change at once (e.g. during a service restart).
package aggregator

import (
	"sync"
	"time"
)

// Event represents a single port change detected by the scanner.
type Event struct {
	Port   int
	Action string // "opened" or "closed"
}

// Summary is the consolidated output emitted after each window.
type Summary struct {
	Opened []int
	Closed []int
	At     time.Time
}

// Aggregator collects events over a fixed window and emits summaries.
type Aggregator struct {
	mu      sync.Mutex
	window  time.Duration
	opened  map[int]struct{}
	closed  map[int]struct{}
	timer   *time.Timer
	out     chan Summary
	stopped bool
}

// New creates an Aggregator that flushes every window duration.
func New(window time.Duration) *Aggregator {
	return &Aggregator{
		window: window,
		opened: make(map[int]struct{}),
		closed: make(map[int]struct{}),
		out:    make(chan Summary, 8),
	}
}

// Add enqueues an event. The window timer is started on the first event.
func (a *Aggregator) Add(e Event) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.stopped {
		return
	}
	switch e.Action {
	case "opened":
		a.opened[e.Port] = struct{}{}
		delete(a.closed, e.Port)
	case "closed":
		a.closed[e.Port] = struct{}{}
		delete(a.opened, e.Port)
	}
	if a.timer == nil {
		a.timer = time.AfterFunc(a.window, a.flush)
	}
}

// Out returns the channel on which summaries are delivered.
func (a *Aggregator) Out() <-chan Summary { return a.out }

// Flush forces an immediate emit regardless of the window.
func (a *Aggregator) Flush() {
	a.mu.Lock()
	if a.timer != nil {
		a.timer.Stop()
		a.timer = nil
	}
	a.mu.Unlock()
	a.flush()
}

// Stop flushes any pending events and closes the output channel.
func (a *Aggregator) Stop() {
	a.Flush()
	a.mu.Lock()
	defer a.mu.Unlock()
	if !a.stopped {
		a.stopped = true
		close(a.out)
	}
}

func (a *Aggregator) flush() {
	a.mu.Lock()
	if len(a.opened) == 0 && len(a.closed) == 0 {
		a.timer = nil
		a.mu.Unlock()
		return
	}
	s := Summary{At: time.Now()}
	for p := range a.opened {
		s.Opened = append(s.Opened, p)
	}
	for p := range a.closed {
		s.Closed = append(s.Closed, p)
	}
	a.opened = make(map[int]struct{})
	a.closed = make(map[int]struct{})
	a.timer = nil
	a.mu.Unlock()
	if !a.stopped {
		a.out <- s
	}
}
