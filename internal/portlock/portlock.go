// Package portlock provides a mechanism to lock (pin) specific ports,
// preventing them from being reported as unexpected changes during scans.
package portlock

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Locker tracks ports that have been explicitly locked by the operator.
type Locker struct {
	mu      sync.RWMutex
	locked  map[int]string // port -> reason
	path    string
}

// New returns a new Locker. If path is non-empty the lock list is persisted
// to disk across restarts.
func New(path string) *Locker {
	return &Locker{
		locked: make(map[int]string),
		path:   path,
	}
}

// Lock pins the given port with an optional human-readable reason.
func (l *Locker) Lock(port int, reason string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.locked[port] = reason
}

// Unlock removes a previously locked port. Returns false if the port was
// not locked.
func (l *Locker) Unlock(port int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.locked[port]; !ok {
		return false
	}
	delete(l.locked, port)
	return true
}

// IsLocked reports whether port is currently locked.
func (l *Locker) IsLocked(port int) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	_, ok := l.locked[port]
	return ok
}

// Reason returns the reason associated with a locked port and whether it
// was found.
func (l *Locker) Reason(port int) (string, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	r, ok := l.locked[port]
	return r, ok
}

// Snapshot returns a copy of the current lock table.
func (l *Locker) Snapshot() map[int]string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make(map[int]string, len(l.locked))
	for k, v := range l.locked {
		out[k] = v
	}
	return out
}

// Save persists the lock table to the configured path.
func (l *Locker) Save() error {
	if l.path == "" {
		return nil
	}
	l.mu.RLock()
	defer l.mu.RUnlock()
	data, err := json.Marshal(l.locked)
	if err != nil {
		return fmt.Errorf("portlock: marshal: %w", err)
	}
	return os.WriteFile(l.path, data, 0o644)
}

// Load restores the lock table from the configured path. A missing file is
// silently ignored.
func (l *Locker) Load() error {
	if l.path == "" {
		return nil
	}
	data, err := os.ReadFile(l.path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("portlock: read: %w", err)
	}
	var m map[int]string
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("portlock: unmarshal: %w", err)
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.locked = m
	return nil
}
