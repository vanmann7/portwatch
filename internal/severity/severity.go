// Package severity classifies port change events into severity levels
// based on configurable rules such as well-known port ranges and custom overrides.
package severity

import "fmt"

// Level represents the severity of a port event.
type Level int

const (
	Info    Level = iota // routine or expected change
	Warning              // notable but not critical
	Critical             // high-risk port activity
)

// String returns a human-readable label for the level.
func (l Level) String() string {
	switch l {
	case Info:
		return "INFO"
	case Warning:
		return "WARNING"
	case Critical:
		return "CRITICAL"
	default:
		return fmt.Sprintf("LEVEL(%d)", int(l))
	}
}

// Classifier assigns severity levels to ports.
type Classifier struct {
	overrides map[int]Level
}

// New returns a Classifier with optional port-level overrides.
// The overrides map takes precedence over built-in rules.
func New(overrides map[int]Level) *Classifier {
	if overrides == nil {
		overrides = make(map[int]Level)
	}
	return &Classifier{overrides: overrides}
}

// Classify returns the severity Level for the given port number.
// Precedence: overrides > well-known critical ports > well-known warning ports > Info.
func (c *Classifier) Classify(port int) Level {
	if lvl, ok := c.overrides[port]; ok {
		return lvl
	}
	if isCritical(port) {
		return Critical
	}
	if isWarning(port) {
		return Warning
	}
	return Info
}

// isCritical returns true for ports associated with high-risk services.
func isCritical(port int) bool {
	switch port {
	case 21, 22, 23, 3389, 5900, 4444, 1337:
		return true
	}
	return false
}

// isWarning returns true for ports that warrant attention but are not critical.
func isWarning(port int) bool {
	switch port {
	case 80, 443, 8080, 8443, 3306, 5432, 6379, 27017:
		return true
	}
	// Ephemeral / dynamic range — flag as warning
	if port >= 1024 && port <= 49151 {
		return true
	}
	return false
}
