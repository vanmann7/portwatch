// Package eventfilter provides a composable event-filtering stage that
// evaluates a chain of predicates against a pipeline event and drops
// events that fail any predicate.
package eventfilter

import "github.com/user/portwatch/internal/pipeline"

// Predicate is a function that returns true when an event should be kept.
type Predicate func(e pipeline.Event) bool

// Filter holds an ordered list of predicates.
type Filter struct {
	predicates []Predicate
}

// New returns a Filter that applies each predicate in order.
// An empty predicate list passes all events through.
func New(predicates ...Predicate) *Filter {
	return &Filter{predicates: predicates}
}

// Allow returns true when the event satisfies every predicate.
func (f *Filter) Allow(e pipeline.Event) bool {
	for _, p := range f.predicates {
		if !p(e) {
			return false
		}
	}
	return true
}

// Add appends a predicate to the filter chain.
func (f *Filter) Add(p Predicate) {
	f.predicates = append(f.predicates, p)
}

// MinPort returns a Predicate that passes events whose port is >= min.
func MinPort(min int) Predicate {
	return func(e pipeline.Event) bool {
		return e.Port >= min
	}
}

// MaxPort returns a Predicate that passes events whose port is <= max.
func MaxPort(max int) Predicate {
	return func(e pipeline.Event) bool {
		return e.Port <= max
	}
}

// OnlyOpened returns a Predicate that passes only opened-port events.
func OnlyOpened() Predicate {
	return func(e pipeline.Event) bool {
		return e.Type == pipeline.EventOpened
	}
}

// OnlyClosed returns a Predicate that passes only closed-port events.
func OnlyClosed() Predicate {
	return func(e pipeline.Event) bool {
		return e.Type == pipeline.EventClosed
	}
}
