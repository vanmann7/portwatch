// Package portmatcher provides pattern-based port matching using glob-style
// rules and named sets, allowing flexible port selection for filtering or
// alerting configurations.
package portmatcher

import (
	"fmt"
	"strconv"
	"strings"
)

// Matcher evaluates whether a port number matches any of the configured rules.
type Matcher struct {
	rules []rule
}

type rule struct {
	low  int
	high int
}

// New creates a Matcher from a slice of rule strings. Each entry may be a
// single port ("80"), a range ("1000-2000"), or a comma-separated list
// ("80,443,8080").
func New(entries []string) (*Matcher, error) {
	m := &Matcher{}
	for _, entry := range entries {
		parts := strings.Split(entry, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			r, err := parseRule(part)
			if err != nil {
				return nil, fmt.Errorf("portmatcher: invalid rule %q: %w", part, err)
			}
			m.rules = append(m.rules, r)
		}
	}
	return m, nil
}

// Match reports whether port is covered by any configured rule.
func (m *Matcher) Match(port int) bool {
	for _, r := range m.rules {
		if port >= r.low && port <= r.high {
			return true
		}
	}
	return false
}

// Len returns the number of parsed rules.
func (m *Matcher) Len() int { return len(m.rules) }

func parseRule(s string) (rule, error) {
	if idx := strings.Index(s, "-"); idx > 0 {
		lo, err := parsePort(s[:idx])
		if err != nil {
			return rule{}, err
		}
		hi, err := parsePort(s[idx+1:])
		if err != nil {
			return rule{}, err
		}
		if lo > hi {
			return rule{}, fmt.Errorf("low %d > high %d", lo, hi)
		}
		return rule{low: lo, high: hi}, nil
	}
	p, err := parsePort(s)
	if err != nil {
		return rule{}, err
	}
	return rule{low: p, high: p}, nil
}

func parsePort(s string) (int, error) {
	v, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0, fmt.Errorf("not a number: %s", s)
	}
	if v < 1 || v > 65535 {
		return 0, fmt.Errorf("port %d out of range [1, 65535]", v)
	}
	return v, nil
}
