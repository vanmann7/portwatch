package portstate

import "fmt"

// String returns a human-readable name for the State.
func (s State) String() string {
	switch s {
	case Unknown:
		return "unknown"
	case Open:
		return "open"
	case Closed:
		return "closed"
	default:
		return fmt.Sprintf("State(%d)", int(s))
	}
}

// IsOpen reports whether the state is Open.
func (s State) IsOpen() bool { return s == Open }

// IsClosed reports whether the state is Closed.
func (s State) IsClosed() bool { return s == Closed }
