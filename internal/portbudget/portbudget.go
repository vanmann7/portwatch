// Package portbudget enforces a maximum number of concurrently open ports
// within a configurable set of ranges. When the budget is exceeded the
// Exceeded method returns true so callers can emit an alert.
package portbudget

import (
	"fmt"
	"sync"
)

// Budget tracks open ports against a per-range limit.
type Budget struct {
	mu      sync.Mutex
	ranges  []portRange
	open    map[int]struct{}
}

type portRange struct {
	lo, hi int
	max    int
}

// Option configures a Budget.
type Option func(*Budget) error

// WithRange registers a port range [lo, hi] with a maximum allowed count.
func WithRange(lo, hi, max int) Option {
	return func(b *Budget) error {
		if lo < 1 || hi > 65535 || lo > hi {
			return fmt.Errorf("portbudget: invalid range %d-%d", lo, hi)
		}
		if max < 1 {
			return fmt.Errorf("portbudget: max must be >= 1")
		}
		b.ranges = append(b.ranges, portRange{lo: lo, hi: hi, max: max})
		return nil
	}
}

// New creates a Budget with the provided options.
func New(opts ...Option) (*Budget, error) {
	b := &Budget{
		open: make(map[int]struct{}),
	}
	for _, o := range opts {
		if err := o(b); err != nil {
			return nil, err
		}
	}
	return b, nil
}

// Open records a port as open.
func (b *Budget) Open(port int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.open[port] = struct{}{}
}

// Close removes a port from the open set.
func (b *Budget) Close(port int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.open, port)
}

// Exceeded reports whether any registered range has more open ports than its
// configured maximum. It also returns the offending range description.
func (b *Budget) Exceeded() (bool, string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, r := range b.ranges {
		count := 0
		for p := range b.open {
			if p >= r.lo && p <= r.hi {
				count++
			}
		}
		if count > r.max {
			return true, fmt.Sprintf("%d-%d (open: %d, max: %d)", r.lo, r.hi, count, r.max)
		}
	}
	return false, ""
}

// Count returns the number of open ports within the given range.
func (b *Budget) Count(lo, hi int) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	count := 0
	for p := range b.open {
		if p >= lo && p <= hi {
			count++
		}
	}
	return count
}
