// Package circuitbreaker implements a simple circuit breaker that trips
// after a configurable number of consecutive failures and resets after a
// cooldown period, preventing repeated attempts against a failing resource.
package circuitbreaker

import (
	"errors"
	"sync"
	"time"
)

// ErrOpen is returned when the circuit breaker is in the open (tripped) state.
var ErrOpen = errors.New("circuit breaker is open")

// State represents the current state of the circuit breaker.
type State int

const (
	StateClosed State = iota // normal operation
	StateOpen                // tripped; calls rejected
	StateHalfOpen            // probe call allowed after reset window
)

// Breaker is a circuit breaker that tracks consecutive failures.
type Breaker struct {
	mu           sync.Mutex
	maxFailures  int
	resetAfter   time.Duration
	failures     int
	state        State
	openedAt     time.Time
	now          func() time.Time
}

// New creates a Breaker that opens after maxFailures consecutive failures
// and attempts a half-open probe after resetAfter has elapsed.
func New(maxFailures int, resetAfter time.Duration) *Breaker {
	return &Breaker{
		maxFailures: maxFailures,
		resetAfter:  resetAfter,
		now:         time.Now,
	}
}

// Allow reports whether a call should be attempted. It returns ErrOpen when
// the breaker is tripped and the reset window has not yet elapsed.
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateClosed:
		return nil
	case StateOpen:
		if b.now().Sub(b.openedAt) >= b.resetAfter {
			b.state = StateHalfOpen
			return nil
		}
		return ErrOpen
	case StateHalfOpen:
		return nil
	}
	return nil
}

// RecordSuccess resets the failure counter and closes the circuit.
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure increments the failure counter and opens the circuit when
// the threshold is reached.
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.maxFailures {
		b.state = StateOpen
		b.openedAt = b.now()
	}
}

// State returns the current circuit breaker state.
func (b *Breaker) State() State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

// Failures returns the current consecutive failure count.
func (b *Breaker) Failures() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.failures
}
