// Package history records scan events over time for trend analysis.
package history

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry represents a single scan event.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Opened    []uint16  `json:"opened,omitempty"`
	Closed    []uint16  `json:"closed,omitempty"`
}

// History holds a bounded list of scan entries.
type History struct {
	mu      sync.Mutex
	entries []Entry
	maxSize int
}

// New creates a History that retains at most maxSize entries.
func New(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &History{maxSize: maxSize}
}

// Record appends a new entry, evicting the oldest if at capacity.
func (h *History) Record(opened, closed []uint16) {
	h.mu.Lock()
	defer h.mu.Unlock()
	e := Entry{Timestamp: time.Now().UTC(), Opened: opened, Closed: closed}
	if len(h.entries) >= h.maxSize {
		h.entries = h.entries[1:]
	}
	h.entries = append(h.entries, e)
}

// Entries returns a copy of all recorded entries.
func (h *History) Entries() []Entry {
	h.mu.Lock()
	defer h.mu.Unlock()
	out := make([]Entry, len(h.entries))
	copy(out, h.entries)
	return out
}

// Since returns all entries recorded at or after the given time.
func (h *History) Since(t time.Time) []Entry {
	h.mu.Lock()
	defer h.mu.Unlock()
	var out []Entry
	for _, e := range h.entries {
		if !e.Timestamp.Before(t) {
			out = append(out, e)
		}
	}
	return out
}

// Save persists the history to a JSON file.
func (h *History) Save(path string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(h.entries)
}

// Load replaces current entries from a JSON file.
func (h *History) Load(path string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	var entries []Entry
	if err := json.NewDecoder(f).Decode(&entries); err != nil {
		return err
	}
	h.entries = entries
	return nil
}
