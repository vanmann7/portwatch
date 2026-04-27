// Package portmap maintains a live mapping of open ports to their
// associated metadata, providing a unified view of the current port
// landscape observed by portwatch.
package portmap

import (
	"fmt"
	"sync"
	"time"
)

// Entry holds the metadata associated with a single open port.
type Entry struct {
	Port      int
	Protocol  string
	Service   string
	FirstSeen time.Time
	LastSeen  time.Time
	Hits      int
}

// PortMap is a thread-safe registry of currently open ports and their
// associated metadata. It is updated on every scan cycle.
type PortMap struct {
	mu      sync.RWMutex
	entries map[int]*Entry
}

// New returns an empty PortMap.
func New() *PortMap {
	return &PortMap{
		entries: make(map[int]*Entry),
	}
}

// Record inserts or updates the entry for the given port.
// If the port has not been seen before, FirstSeen is set to now.
// LastSeen and Hits are always updated.
func (m *PortMap) Record(port int, protocol, service string, now time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if e, ok := m.entries[port]; ok {
		e.LastSeen = now
		e.Hits++
		// Update service/protocol if a more specific value is provided.
		if protocol != "" {
			e.Protocol = protocol
		}
		if service != "" {
			e.Service = service
		}
		return
	}

	m.entries[port] = &Entry{
		Port:      port,
		Protocol:  protocol,
		Service:   service,
		FirstSeen: now,
		LastSeen:  now,
		Hits:      1,
	}
}

// Remove deletes the entry for the given port. It is a no-op if the
// port is not present.
func (m *PortMap) Remove(port int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.entries, port)
}

// Get returns the entry for the given port and whether it was found.
func (m *PortMap) Get(port int) (Entry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.entries[port]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// Ports returns a sorted snapshot of all currently tracked port numbers.
func (m *PortMap) Ports() []int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ports := make([]int, 0, len(m.entries))
	for p := range m.entries {
		ports = append(ports, p)
	}
	sortInts(ports)
	return ports
}

// Snapshot returns a copy of all entries keyed by port number.
func (m *PortMap) Snapshot() map[int]Entry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make(map[int]Entry, len(m.entries))
	for p, e := range m.entries {
		out[p] = *e
	}
	return out
}

// Len returns the number of ports currently tracked.
func (m *PortMap) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.entries)
}

// String returns a human-readable summary of the port map.
func (m *PortMap) String() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return fmt.Sprintf("PortMap{count: %d}", len(m.entries))
}

// sortInts performs an in-place insertion sort on a small slice of ints.
// For the typical number of open ports this is fast enough without
// importing sort.
func sortInts(s []int) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
