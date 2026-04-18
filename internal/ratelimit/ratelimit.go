// Package ratelimit provides a simple token-bucket rate limiter
// to prevent alert flooding when many ports change in a short window.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter controls how many events may pass per interval.
type Limiter struct {
	mu       sync.Mutex
	tokens   int
	max      int
	interval time.Duration
	last     time.Time
}

// New creates a Limiter that allows at most max events per interval.
func New(max int, interval time.Duration) *Limiter {
	return &Limiter{
		tokens:   max,
		max:      max,
		interval: interval,
		last:     time.Now(),
	}
}

// Allow reports whether an event should be allowed through.
// It refills tokens proportionally to elapsed time before deciding.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(l.last)

	if elapsed >= l.interval {
		refill := int(elapsed/l.interval) * l.max
		l.tokens += refill
		if l.tokens > l.max {
			l.tokens = l.max
		}
		l.last = now
	}

	if l.tokens <= 0 {
		return false
	}
	l.tokens--
	return true
}

// Reset restores the limiter to its full token capacity.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.tokens = l.max
	l.last = time.Now()
}
