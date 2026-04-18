// Package dedupe provides event deduplication to suppress repeated
// notifications for the same port change within a configurable window.
package dedupe

import (
	"sync"
	"time"
)

// key uniquely identifies a port event.
type key struct {
	port  int
	event string
}

// Deduper suppresses duplicate port events within a time window.
type Deduper struct {
	mu     sync.Mutex
	window time.Duration
	seen   map[key]time.Time
	now    func() time.Time
}

// New returns a Deduper that suppresses repeated events within window.
func New(window time.Duration) *Deduper {
	return &Deduper{
		window: window,
		seen:   make(map[key]time.Time),
		now:    time.Now,
	}
}

// IsDuplicate returns true if the same (port, event) pair was seen
// within the configured window. If not a duplicate, it records the event.
func (d *Deduper) IsDuplicate(port int, event string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.evict()

	k := key{port: port, event: event}
	if _, ok := d.seen[k]; ok {
		return true
	}
	d.seen[k] = d.now()
	return false
}

// Reset clears all recorded events.
func (d *Deduper) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[key]time.Time)
}

// evict removes entries older than the window. Must be called with lock held.
func (d *Deduper) evict() {
	cutoff := d.now().Add(-d.window)
	for k, t := range d.seen {
		if t.Before(cutoff) {
			delete(d.seen, k)
		}
	}
}
