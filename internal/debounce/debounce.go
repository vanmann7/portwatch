// Package debounce provides a simple debouncer that delays action execution
// until a quiet period has elapsed, preventing alert floods during rapid port churn.
package debounce

import (
	"sync"
	"time"
)

// Debouncer delays calls to fn until after wait has elapsed since the last call.
type Debouncer struct {
	wait  time.Duration
	fn    func()
	mu    sync.Mutex
	timer *time.Timer
}

// New creates a new Debouncer that will call fn after wait duration of inactivity.
func New(wait time.Duration, fn func()) *Debouncer {
	if wait <= 0 {
		wait = 100 * time.Millisecond
	}
	return &Debouncer{wait: wait, fn: fn}
}

// Trigger schedules fn to be called after the debounce window.
// If Trigger is called again before the window expires, the timer resets.
func (d *Debouncer) Trigger() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.wait, func() {
		d.fn()
	})
}

// Flush cancels any pending timer and calls fn immediately if one was pending.
// Returns true if a pending call was flushed.
func (d *Debouncer) Flush() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.timer == nil {
		return false
	}
	stopped := d.timer.Stop()
	d.timer = nil
	if stopped {
		go d.fn()
		return true
	}
	return false
}

// Stop cancels any pending timer without calling fn.
func (d *Debouncer) Stop() {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.timer != nil {
		d.timer.Stop()
		d.timer = nil
	}
}
