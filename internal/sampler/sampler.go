// Package sampler provides periodic port-scan sampling with configurable
// interval and jitter to avoid thundering-herd effects when multiple
// portwatch instances run on the same host.
package sampler

import (
	"context"
	"math/rand"
	"time"
)

// ScanFunc is the function signature expected by the Sampler to perform a
// single port scan and return the set of open ports.
type ScanFunc func(ctx context.Context) ([]int, error)

// Sampler fires a ScanFunc at a regular interval, optionally adding a random
// jitter up to JitterMax before each tick to spread load.
type Sampler struct {
	interval  time.Duration
	jitterMax time.Duration
	scan      ScanFunc
}

// New creates a Sampler with the given interval, jitter ceiling, and scan
// function. A zero jitterMax disables jitter entirely.
func New(interval, jitterMax time.Duration, scan ScanFunc) *Sampler {
	return &Sampler{
		interval:  interval,
		jitterMax: jitterMax,
		scan:      scan,
	}
}

// Result carries the outcome of a single sample tick.
type Result struct {
	Ports []int
	Err   error
	At    time.Time
}

// Run starts the sampling loop and sends results to the returned channel.
// The channel is closed when ctx is cancelled. The first sample fires after
// the initial interval (plus jitter).
func (s *Sampler) Run(ctx context.Context) <-chan Result {
	ch := make(chan Result)
	go func() {
		defer close(ch)
		for {
			wait := s.interval
			if s.jitterMax > 0 {
				//nolint:gosec // non-cryptographic jitter is intentional
				wait += time.Duration(rand.Int63n(int64(s.jitterMax)))
			}
			select {
			case <-ctx.Done():
				return
			case <-time.After(wait):
			}
			ports, err := s.scan(ctx)
			select {
			case ch <- Result{Ports: ports, Err: err, At: time.Now()}:
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch
}
