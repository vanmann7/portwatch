// Package portschedule provides a scheduler that controls when individual
// ports are eligible to be re-scanned, enforcing per-port scan intervals.
package portschedule

import (
	"sync"
	"time"
)

// Entry tracks the last scan time for a single port.
type Entry struct {
	Port     int
	LastScan time.Time
}

// Scheduler decides whether a port is due for re-scanning based on a
// minimum interval between successive scans of the same port.
type Scheduler struct {
	mu       sync.Mutex
	entries  map[int]time.Time
	interval time.Duration
	now      func() time.Time
}

// New returns a Scheduler that enforces the given minimum interval between
// scans of the same port.
func New(interval time.Duration) *Scheduler {
	return &Scheduler{
		entries:  make(map[int]time.Time),
		interval: interval,
		now:      time.Now,
	}
}

// Due reports whether port is eligible for scanning. A port is due when it
// has never been scanned or its last scan occurred at least interval ago.
func (s *Scheduler) Due(port int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	last, seen := s.entries[port]
	if !seen {
		return true
	}
	return s.now().Sub(last) >= s.interval
}

// Record marks port as having been scanned right now.
func (s *Scheduler) Record(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[port] = s.now()
}

// Reset removes the scheduling record for port, making it immediately due.
func (s *Scheduler) Reset(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, port)
}

// Snapshot returns a copy of all tracked entries at the current moment.
func (s *Scheduler) Snapshot() []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Entry, 0, len(s.entries))
	for port, t := range s.entries {
		out = append(out, Entry{Port: port, LastScan: t})
	}
	return out
}
