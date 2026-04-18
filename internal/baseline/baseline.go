// Package baseline manages the trusted port baseline used to distinguish
// expected from unexpected open ports.
package baseline

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// Baseline holds the set of ports considered "known good".
type Baseline struct {
	mu    sync.RWMutex
	ports map[int]struct{}
	path  string
}

// New returns an empty Baseline backed by the given file path.
func New(path string) *Baseline {
	return &Baseline{ports: make(map[int]struct{}), path: path}
}

// Set replaces the current baseline with the provided ports.
func (b *Baseline) Set(ports []int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.ports = make(map[int]struct{}, len(ports))
	for _, p := range ports {
		b.ports[p] = struct{}{}
	}
}

// Contains reports whether port is in the baseline.
func (b *Baseline) Contains(port int) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	_, ok := b.ports[port]
	return ok
}

// Ports returns a sorted snapshot of the current baseline ports.
func (b *Baseline) Ports() []int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]int, 0, len(b.ports))
	for p := range b.ports {
		out = append(out, p)
	}
	return out
}

// Save persists the baseline to disk as JSON.
func (b *Baseline) Save() error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	ports := make([]int, 0, len(b.ports))
	for p := range b.ports {
		ports = append(ports, p)
	}
	data, err := json.Marshal(ports)
	if err != nil {
		return err
	}
	return os.WriteFile(b.path, data, 0o600)
}

// Load reads the baseline from disk. Missing file is not an error.
func (b *Baseline) Load() error {
	data, err := os.ReadFile(b.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	var ports []int
	if err := json.Unmarshal(data, &ports); err != nil {
		return err
	}
	b.Set(ports)
	return nil
}
