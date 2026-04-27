// Package portmemo provides a short-term memory store for port events,
// allowing components to recall the most recent event seen for a given port.
package portmemo

import (
	"sync"
	"time"
)

// Entry holds the most recent event recorded for a port.
type Entry struct {
	Port      int
	Event     string // "opened" or "closed"
	RecordedAt time.Time
}

// Memo stores the last-seen event for each port.
type Memo struct {
	mu      sync.RWMutex
	entries map[int]Entry
	clock   func() time.Time
}

// New returns a new Memo.
func New() *Memo {
	return newWithClock(time.Now)
}

func newWithClock(clock func() time.Time) *Memo {
	return &Memo{
		entries: make(map[int]Entry),
		clock:   clock,
	}
}

// Record stores the most recent event for the given port, overwriting any
// previous entry.
func (m *Memo) Record(port int, event string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[port] = Entry{
		Port:       port,
		Event:      event,
		RecordedAt: m.clock(),
	}
}

// Get returns the most recent Entry for port and true, or a zero Entry and
// false if the port has not been recorded.
func (m *Memo) Get(port int) (Entry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.entries[port]
	return e, ok
}

// Forget removes the entry for port, if any.
func (m *Memo) Forget(port int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, port)
}

// Snapshot returns a copy of all current entries keyed by port.
func (m *Memo) Snapshot() map[int]Entry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make(map[int]Entry, len(m.entries))
	for k, v := range m.entries {
		out[k] = v
	}
	return out
}

// Clear removes all entries.
func (m *Memo) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = make(map[int]Entry)
}
