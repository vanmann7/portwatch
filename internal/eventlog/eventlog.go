// Package eventlog provides a structured, in-memory ring buffer for
// recording port change events with timestamps and metadata.
package eventlog

import (
	"sync"
	"time"
)

// EventKind describes the type of port change event.
type EventKind string

const (
	EventOpened EventKind = "opened"
	EventClosed EventKind = "closed"

	// DefaultCapacity is the default maximum number of events retained.
	DefaultCapacity = 256
)

// Entry represents a single recorded port event.
type Entry struct {
	Port      int
	Kind      EventKind
	Service   string
	Timestamp time.Time
}

// EventLog is a fixed-capacity ring buffer of port events.
type EventLog struct {
	mu       sync.Mutex
	entries  []Entry
	cap      int
	head     int
	count    int
}

// New returns an EventLog with the given capacity.
// If capacity is <= 0, DefaultCapacity is used.
func New(capacity int) *EventLog {
	if capacity <= 0 {
		capacity = DefaultCapacity
	}
	return &EventLog{
		entries: make([]Entry, capacity),
		cap:     capacity,
	}
}

// Record appends a new event entry to the log.
func (l *EventLog) Record(port int, kind EventKind, service string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.entries[l.head] = Entry{
		Port:      port,
		Kind:      kind,
		Service:   service,
		Timestamp: time.Now(),
	}
	l.head = (l.head + 1) % l.cap
	if l.count < l.cap {
		l.count++
	}
}

// Entries returns a snapshot of all recorded events in chronological order.
func (l *EventLog) Entries() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.count == 0 {
		return nil
	}
	out := make([]Entry, l.count)
	start := (l.head - l.count + l.cap) % l.cap
	for i := 0; i < l.count; i++ {
		out[i] = l.entries[(start+i)%l.cap]
	}
	return out
}

// Len returns the number of entries currently stored.
func (l *EventLog) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.count
}

// Clear removes all entries from the log.
func (l *EventLog) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.head = 0
	l.count = 0
}
