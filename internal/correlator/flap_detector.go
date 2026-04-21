package correlator

import (
	"sync"
	"time"
)

// FlapDetector identifies ports that oscillate between open and closed
// within a short observation window — a common sign of instability.
type FlapDetector struct {
	mu      sync.Mutex
	window  time.Duration
	counts  map[int][]time.Time
	thresh  int
}

// NewFlapDetector returns a FlapDetector that flags a port as flapping
// when it changes state more than thresh times within window.
func NewFlapDetector(window time.Duration, thresh int) *FlapDetector {
	return &FlapDetector{
		window: window,
		thresh: thresh,
		counts: make(map[int][]time.Time),
	}
}

// Observe records a state-change event for port and reports whether
// the port is currently considered to be flapping.
func (f *FlapDetector) Observe(port int) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-f.window)

	times := f.counts[port]
	filtered := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	filtered = append(filtered, now)
	f.counts[port] = filtered

	return len(filtered) > f.thresh
}

// IsFlapping reports whether port is currently flapping without
// recording a new observation.
func (f *FlapDetector) IsFlapping(port int) bool {
	f.mu.Lock()
	defer f.mu.Unlock()

	cutoff := time.Now().Add(-f.window)
	count := 0
	for _, t := range f.counts[port] {
		if t.After(cutoff) {
			count++
		}
	}
	return count > f.thresh
}

// Reset clears all recorded observations for all ports.
func (f *FlapDetector) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.counts = make(map[int][]time.Time)
}
