// Package suppress provides a suppression list for silencing alerts
// on known or expected ports so operators are not flooded with noise.
package suppress

import (
	"encoding/json"
	"os"
	"sort"
	"sync"
)

// List holds a set of suppressed ports.
type List struct {
	mu    sync.RWMutex
	ports map[int]struct{}
}

// New returns an empty suppression list.
func New() *List {
	return &List{ports: make(map[int]struct{})}
}

// Add adds one or more ports to the suppression list.
func (l *List) Add(ports ...int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, p := range ports {
		l.ports[p] = struct{}{}
	}
}

// Remove removes a port from the suppression list.
func (l *List) Remove(port int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.ports, port)
}

// IsSuppressed reports whether the given port is suppressed.
func (l *List) IsSuppressed(port int) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	_, ok := l.ports[port]
	return ok
}

// Ports returns a sorted snapshot of all suppressed ports.
func (l *List) Ports() []int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]int, 0, len(l.ports))
	for p := range l.ports {
		out = append(out, p)
	}
	sort.Ints(out)
	return out
}

// Save persists the suppression list to a JSON file.
func (l *List) Save(path string) error {
	l.mu.RLock()
	defer l.mu.RUnlock()
	ports := make([]int, 0, len(l.ports))
	for p := range l.ports {
		ports = append(ports, p)
	}
	sort.Ints(ports)
	data, err := json.Marshal(ports)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Load reads a suppression list from a JSON file.
func Load(path string) (*List, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var ports []int
	if err := json.Unmarshal(data, &ports); err != nil {
		return nil, err
	}
	l := New()
	l.Add(ports...)
	return l, nil
}
