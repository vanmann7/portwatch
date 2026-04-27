// Package portage tracks how long each port has been continuously open
// and classifies ports into age buckets (new, established, long-running).
package portage

import (
	"sync"
	"time"
)

// Bucket classifies a port by how long it has been open.
type Bucket string

const (
	BucketNew         Bucket = "new"          // open < 5 minutes
	BucketEstablished Bucket = "established"  // open 5–60 minutes
	BucketLongRunning Bucket = "long-running" // open > 60 minutes
)

// Entry holds age metadata for a single port.
type Entry struct {
	Port     int
	OpenedAt time.Time
	Age      time.Duration
	Bucket   Bucket
}

type clock func() time.Time

// Tracker records when each port was first seen and classifies its age.
type Tracker struct {
	mu      sync.Mutex
	openedAt map[int]time.Time
	now      clock
}

// New returns a new Tracker using the real wall clock.
func New() *Tracker {
	return newWithClock(time.Now)
}

func newWithClock(c clock) *Tracker {
	return &Tracker{
		openedAt: make(map[int]time.Time),
		now:      c,
	}
}

// Record marks port as open at the current time (idempotent).
func (t *Tracker) Record(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, exists := t.openedAt[port]; !exists {
		t.openedAt[port] = t.now()
	}
}

// Remove forgets a port (called when it closes).
func (t *Tracker) Remove(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.openedAt, port)
}

// Get returns the Entry for a port, and whether it was found.
func (t *Tracker) Get(port int) (Entry, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	op, ok := t.openedAt[port]
	if !ok {
		return Entry{}, false
	}
	age := t.now().Sub(op)
	return Entry{
		Port:     port,
		OpenedAt: op,
		Age:      age,
		Bucket:   classify(age),
	}, true
}

// Snapshot returns Entries for all tracked ports.
func (t *Tracker) Snapshot() []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Entry, 0, len(t.openedAt))
	for port, op := range t.openedAt {
		age := t.now().Sub(op)
		out = append(out, Entry{
			Port:     port,
			OpenedAt: op,
			Age:      age,
			Bucket:   classify(age),
		})
	}
	return out
}

func classify(age time.Duration) Bucket {
	switch {
	case age < 5*time.Minute:
		return BucketNew
	case age < 60*time.Minute:
		return BucketEstablished
	default:
		return BucketLongRunning
	}
}
