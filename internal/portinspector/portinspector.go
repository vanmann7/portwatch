// Package portinspector provides a unified view of a port's current status,
// combining state, lifetime, rank, and label into a single inspectable record.
package portinspector

import (
	"fmt"
	"time"
)

// Record holds all known metadata about a single port at inspection time.
type Record struct {
	Port     int
	State    string
	Label    string
	Rank     int
	Lifetime time.Duration
	OpenedAt time.Time
}

// String returns a human-readable summary of the record.
func (r Record) String() string {
	if r.State == "closed" {
		return fmt.Sprintf("port %d [%s] rank=%d label=%q", r.Port, r.State, r.Rank, r.Label)
	}
	return fmt.Sprintf("port %d [%s] rank=%d label=%q lifetime=%s",
		r.Port, r.State, r.Rank, r.Label, r.Lifetime.Truncate(time.Second))
}

// StateProvider returns the current state string for a port.
type StateProvider interface {
	Get(port int) (string, bool)
}

// LabelProvider returns the label for a port.
type LabelProvider interface {
	Label(port int) string
}

// RankProvider returns the numeric rank for a port.
type RankProvider interface {
	Rank(port int) int
}

// LifetimeProvider returns the elapsed lifetime for a port.
type LifetimeProvider interface {
	Lifetime(port int) (time.Duration, bool)
	OpenedAt(port int) (time.Time, bool)
}

// Inspector combines multiple providers to produce a unified Record.
type Inspector struct {
	state    StateProvider
	labels   LabelProvider
	ranks    RankProvider
	lifetime LifetimeProvider
}

// New creates a new Inspector backed by the given providers.
func New(s StateProvider, l LabelProvider, r RankProvider, lt LifetimeProvider) *Inspector {
	return &Inspector{state: s, labels: l, ranks: r, lifetime: lt}
}

// Inspect returns a Record for the given port.
// If the port has no known state, the second return value is false.
func (ins *Inspector) Inspect(port int) (Record, bool) {
	state, ok := ins.state.Get(port)
	if !ok {
		return Record{}, false
	}
	rec := Record{
		Port:  port,
		State: state,
		Label: ins.labels.Label(port),
		Rank:  ins.ranks.Rank(port),
	}
	if d, ok := ins.lifetime.Lifetime(port); ok {
		rec.Lifetime = d
	}
	if t, ok := ins.lifetime.OpenedAt(port); ok {
		rec.OpenedAt = t
	}
	return rec, true
}
