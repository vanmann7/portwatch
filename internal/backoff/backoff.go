// Package backoff provides exponential backoff with jitter for retry delays.
package backoff

import (
	"math"
	"math/rand"
	"time"
)

// Strategy defines how delays grow between retry attempts.
type Strategy struct {
	// Base is the initial delay before the first retry.
	Base time.Duration
	// Max caps the delay regardless of the exponent.
	Max time.Duration
	// Factor is the multiplicative growth factor (e.g. 2.0 for doubling).
	Factor float64
	// Jitter adds a random fraction of the computed delay to spread load.
	Jitter bool
}

// Default returns a Strategy suitable for most port-scan retry scenarios.
func Default() Strategy {
	return Strategy{
		Base:   250 * time.Millisecond,
		Max:    30 * time.Second,
		Factor: 2.0,
		Jitter: true,
	}
}

// Delay returns the wait duration for the given attempt number (0-indexed).
// Attempt 0 returns Base; subsequent attempts grow exponentially up to Max.
func (s Strategy) Delay(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}

	exp := math.Pow(s.Factor, float64(attempt))
	delay := float64(s.Base) * exp

	if s.Jitter {
		// Add up to 50 % random jitter.
		delay += rand.Float64() * delay * 0.5 //nolint:gosec
	}

	result := time.Duration(delay)
	if result > s.Max {
		result = s.Max
	}
	return result
}

// Reset returns the base delay, useful for restarting a backoff sequence.
func (s Strategy) Reset() time.Duration {
	return s.Base
}
