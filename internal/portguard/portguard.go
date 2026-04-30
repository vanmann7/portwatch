// Package portguard enforces access policies on port events, blocking
// or flagging ports that violate configured allow/deny rules and thresholds.
package portguard

import (
	"fmt"
	"sync"
)

// Verdict describes the outcome of a guard evaluation.
type Verdict int

const (
	Allow Verdict = iota
	Deny
	Warn
)

//go:generate stringer -type=Verdict

// Rule defines a single guard rule.
type Rule struct {
	Low    int
	High   int
	Verdict Verdict
}

// Guard evaluates port numbers against an ordered list of rules.
type Guard struct {
	mu      sync.RWMutex
	rules   []Rule
	default_ Verdict
}

// New creates a Guard with the given rules and a fallback default verdict.
func New(rules []Rule, defaultVerdict Verdict) *Guard {
	return &Guard{
		rules:    rules,
		default_: defaultVerdict,
	}
}

// Evaluate returns the verdict for the given port.
func (g *Guard) Evaluate(port int) Verdict {
	g.mu.RLock()
	defer g.mu.RUnlock()
	for _, r := range g.rules {
		if port >= r.Low && port <= r.High {
			return r.Verdict
		}
	}
	return g.default_
}

// AddRule appends a new rule at the end of the evaluation chain.
func (g *Guard) AddRule(r Rule) error {
	if r.Low < 1 || r.High > 65535 || r.Low > r.High {
		return fmt.Errorf("portguard: invalid port range %d-%d", r.Low, r.High)
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	g.rules = append(g.rules, r)
	return nil
}

// Rules returns a snapshot of the current rule list.
func (g *Guard) Rules() []Rule {
	g.mu.RLock()
	defer g.mu.RUnlock()
	out := make([]Rule, len(g.rules))
	copy(out, g.rules)
	return out
}
