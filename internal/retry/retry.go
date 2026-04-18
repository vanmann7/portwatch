package retry

import (
	"context"
	"time"
)

// Policy defines how retries are performed.
type Policy struct {
	MaxAttempts int
	Delay       time.Duration
	Multiplier  float64
}

// Default returns a Policy with sensible defaults.
func Default() Policy {
	return Policy{
		MaxAttempts: 3,
		Delay:       200 * time.Millisecond,
		Multiplier:  2.0,
	}
}

// Retryer executes a function with retry logic.
type Retryer struct {
	policy Policy
}

// New creates a Retryer with the given Policy.
func New(p Policy) *Retryer {
	return &Retryer{policy: p}
}

// Do calls fn up to MaxAttempts times, backing off between attempts.
// It stops early if ctx is cancelled or fn returns nil.
func (r *Retryer) Do(ctx context.Context, fn func() error) error {
	delay := r.policy.Delay
	var err error
	for attempt := 0; attempt < r.policy.MaxAttempts; attempt++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		err = fn()
		if err == nil {
			return nil
		}
		if attempt < r.policy.MaxAttempts-1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
			delay = time.Duration(float64(delay) * r.policy.Multiplier)
		}
	}
	return err
}
