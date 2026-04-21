package watchdog

import "fmt"

// String returns a human-readable representation of a Status value.
func (s Status) String() string {
	switch s {
	case StatusHealthy:
		return "healthy"
	case StatusStale:
		return "stale"
	case StatusDead:
		return "dead"
	default:
		return fmt.Sprintf("unknown(%d)", int(s))
	}
}

// IsHealthy reports whether the status represents a live component.
func (s Status) IsHealthy() bool { return s == StatusHealthy }

// IsDegraded reports whether the component is stale or dead.
func (s Status) IsDegraded() bool { return s == StatusStale || s == StatusDead }
