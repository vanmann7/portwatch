// Package snapshot provides periodic port snapshot scheduling.
package snapshot

import (
	"context"
	"time"
)

// Result holds the outcome of a single scheduled snapshot.
type Result struct {
	Ports []int
	At    time.Time
	Err   error
}

// ScanFunc is the signature of a function that scans ports.
type ScanFunc func(ctx context.Context) ([]int, error)

// Scheduler triggers a ScanFunc at a fixed interval and sends results on a channel.
type Scheduler struct {
	scan     ScanFunc
	interval time.Duration
	results  chan Result
}

// New creates a new Scheduler.
func New(scan ScanFunc, interval time.Duration) *Scheduler {
	return &Scheduler{
		scan:     scan,
		interval: interval,
		results:  make(chan Result, 8),
	}
}

// Results returns the read-only channel of snapshot results.
func (s *Scheduler) Results() <-chan Result {
	return s.results
}

// Run starts the scheduler loop. It blocks until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) {
	defer close(s.results)

	// Run immediately, then on each tick.
	s.emit(ctx)

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.emit(ctx)
		}
{
	portsAt: time.Now(), Ports Err: err}
	select {
	case s.results <- r:
	case <-ctx.Done():
	}
}
