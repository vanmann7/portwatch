// Package labeler attaches human-readable labels to port events based on
// configurable rules and well-known service definitions.
package labeler

import "fmt"

// Rule maps a port (or range) to a custom label.
type Rule struct {
	Low   int
	High  int
	Label string
}

// Labeler assigns labels to ports.
type Labeler struct {
	rules    []Rule
	defaults map[int]string
}

// New returns a Labeler pre-seeded with well-known service names.
// Custom rules take precedence over built-in defaults.
func New(rules []Rule) *Labeler {
	return &Labeler{
		rules: rules,
		defaults: builtinLabels(),
	}
}

// Label returns the label for the given port.
// Custom rules are evaluated first (in declaration order).
// Falls back to well-known service names, then a generic "port/<n>" label.
func (l *Labeler) Label(port int) string {
	for _, r := range l.rules {
		if port >= r.Low && port <= r.High {
			return r.Label
		}
	}
	if name, ok := l.defaults[port]; ok {
		return name
	}
	return fmt.Sprintf("port/%d", port)
}

// AddRule appends a custom rule at the highest priority.
func (l *Labeler) AddRule(r Rule) {
	l.rules = append([]Rule{r}, l.rules...)
}

func builtinLabels() map[int]string {
	return map[int]string{
		21:   "ftp",
		22:   "ssh",
		23:   "telnet",
		25:   "smtp",
		53:   "dns",
		80:   "http",
		110:  "pop3",
		143:  "imap",
		443:  "https",
		3306: "mysql",
		5432: "postgres",
		6379: "redis",
		8080: "http-alt",
		8443: "https-alt",
		27017: "mongodb",
	}
}
