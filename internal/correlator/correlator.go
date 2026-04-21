package correlator

import (
	"fmt"
	"sync"
	"time"
)

// EventKind represents whether a port was opened or closed.
type EventKind string

const (
	Opened EventKind = "opened"
	Closed EventKind = "closed"
)

// Event represents a single port change event with metadata.
type Event struct {
	Port      int
	Kind      EventKind
	Timestamp time.Time
	CorrelID  string
}

// Correlator groups related port events by a shared correlation ID
// derived from the port number and a time-bucketed window.
type Correlator struct {
	mu       sync.Mutex
	window   time.Duration
	groups   map[string][]Event
}

// New returns a Correlator that buckets events within the given window.
func New(window time.Duration) *Correlator {
	return &Correlator{
		window:  window,
		groups:  make(map[string][]Event),
	}
}

// Record adds an event and returns its assigned correlation ID.
func (c *Correlator) Record(port int, kind EventKind) string {
	c.mu.Lock()
	defer c.mu.Unlock()

	bucket := time.Now().Truncate(c.window).Unix()
	id := fmt.Sprintf("%d-%d", port, bucket)

	e := Event{
		Port:      port,
		Kind:      kind,
		Timestamp: time.Now(),
		CorrelID:  id,
	}
	c.groups[id] = append(c.groups[id], e)
	return id
}

// Events returns all events associated with the given correlation ID.
func (c *Correlator) Events(correlID string) []Event {
	c.mu.Lock()
	defer c.mu.Unlock()

	events := c.groups[correlID]
	result := make([]Event, len(events))
	copy(result, events)
	return result
}

// Flush removes and returns all event groups older than the window.
func (c *Correlator) Flush() map[string][]Event {
	c.mu.Lock()
	defer c.mu.Unlock()

	cutoff := time.Now().Truncate(c.window).Unix()
	flushed := make(map[string][]Event)

	for id, events := range c.groups {
		if len(events) == 0 {
			continue
		}
		bucketTime := events[0].Timestamp.Truncate(c.window).Unix()
		if bucketTime < cutoff {
			flushed[id] = events
			delete(c.groups, id)
		}
	}
	return flushed
}

// Reset clears all stored event groups.
func (c *Correlator) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.groups = make(map[string][]Event)
}
