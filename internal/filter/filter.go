// Package filter provides port filtering utilities for portwatch.
package filter

import "fmt"

// Rule describes a single port filter rule.
type Rule struct {
	Low  int
	High int
}

// Filter holds a set of rules that define which ports to include or exclude.
type Filter struct {
	rules []Rule
}

// New creates a Filter from a slice of range strings like "22", "8000-9000".
func New(ranges []string) (*Filter, error) {
	f := &Filter{}
	for _, r := range ranges {
		rule, err := parseRule(r)
		if err != nil {
			return nil, err
		}
		f.rules = append(f.rules, rule)
	}
	return f, nil
}

// Allow reports whether port p is covered by at least one rule.
// If the Filter has no rules every port is allowed.
func (f *Filter) Allow(p int) bool {
	if len(f.rules) == 0 {
		return true
	}
	for _, r := range f.rules {
		if p >= r.Low && p <= r.High {
			return true
		}
	}
	return false
}

// parseRule converts "22" or "8000-9000" into a Rule.
func parseRule(s string) (Rule, error) {
	var lo, hi int
	if _, err := fmt.Sscanf(s, "%d-%d", &lo, &hi); err == nil {
		if lo < 1 || hi > 65535 || lo > hi {
			return Rule{}, fmt.Errorf("filter: invalid range %q", s)
		}
		return Rule{Low: lo, High: hi}, nil
	}
	if _, err := fmt.Sscanf(s, "%d", &lo); err == nil {
		if lo < 1 || lo > 65535 {
			return Rule{}, fmt.Errorf("filter: invalid port %q", s)
		}
		return Rule{Low: lo, High: lo}, nil
	}
	return Rule{}, fmt.Errorf("filter: cannot parse rule %q", s)
}
