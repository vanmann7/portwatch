// Package portcache provides a short-lived in-memory cache for port scan
// results, reducing redundant probes within a configurable TTL window.
package portcache

import (
	"sync"
	"time"
)

// Entry holds a cached scan result for a single port.
type Entry struct {
	Open      bool
	CachedAt  time.Time
	ExpiresAt time.Time
}

// Cache stores port scan results with a fixed TTL.
type Cache struct {
	mu      sync.RWMutex
	entries map[uint16]Entry
	ttl     time.Duration
	now     func() time.Time
}

// New returns a Cache with the given TTL.
func New(ttl time.Duration) *Cache {
	return newWithClock(ttl, time.Now)
}

func newWithClock(ttl time.Duration, now func() time.Time) *Cache {
	return &Cache{
		entries: make(map[uint16]Entry),
		ttl:     ttl,
		now:     now,
	}
}

// Set stores the open/closed state for a port, replacing any existing entry.
func (c *Cache) Set(port uint16, open bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	at := c.now()
	c.entries[port] = Entry{
		Open:      open,
		CachedAt:  at,
		ExpiresAt: at.Add(c.ttl),
	}
}

// Get returns the cached entry for a port and whether it is still valid.
func (c *Cache) Get(port uint16) (Entry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[port]
	if !ok {
		return Entry{}, false
	}
	if c.now().After(e.ExpiresAt) {
		return Entry{}, false
	}
	return e, true
}

// Evict removes all entries whose TTL has elapsed.
func (c *Cache) Evict() {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := c.now()
	for port, e := range c.entries {
		if now.After(e.ExpiresAt) {
			delete(c.entries, port)
		}
	}
}

// Len returns the number of entries currently held (including expired ones
// not yet evicted).
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
