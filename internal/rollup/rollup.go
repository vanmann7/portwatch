// Package rollup groups rapid-fire port events into a single batched
// notification, reducing alert noise during large topology changes.
package rollup

import (
	"sync"
	"time"
)

// Event represents a single port-change event to be rolled up.
type Event struct {
	Port   int
	Action string // "opened" or "closed"
}

// Batch is a collection of events flushed after the window expires.
type Batch struct {
	Events []Event
}

// Roller accumulates events within a sliding window and flushes them
// as a single Batch when the window closes or Flush is called.
type Roller struct {
	mu      sync.Mutex
	window  time.Duration
	pending []Event
	timer   *time.Timer
	out     chan Batch
}

// New creates a Roller that batches events over the given window duration.
// Batches are delivered on the returned channel.
func New(window time.Duration) (*Roller, <-chan Batch) {
	ch := make(chan Batch, 16)
	return &Roller{
		window: window,
		out:    ch,
	}, ch
}

// Add queues an event. If no timer is running, one is started for the
// configured window; subsequent adds reset the timer (sliding window).
func (r *Roller) Add(e Event) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.pending = append(r.pending, e)

	if r.timer != nil {
		r.timer.Reset(r.window)
		return
	}

	r.timer = time.AfterFunc(r.window, func() {
		r.mu.Lock()
		batch := Batch{Events: r.pending}
		r.pending = nil
		r.timer = nil
		r.mu.Unlock()
		r.out <- batch
	})
}

// Flush immediately emits any pending events as a Batch, cancelling
// the in-flight timer. It is a no-op when the queue is empty.
func (r *Roller) Flush() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.pending) == 0 {
		return
	}

	if r.timer != nil {
		r.timer.Stop()
		r.timer = nil
	}

	batch := Batch{Events: r.pending}
	r.pending = nil
	r.out <- batch
}

// Stop cancels any pending timer and closes the output channel.
func (r *Roller) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.timer != nil {
		r.timer.Stop()
		r.timer = nil
	}

	r.pending = nil
	close(r.out)
}
