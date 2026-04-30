// Package portpolicy evaluates ports against a set of named allow/deny rules
// and returns a verdict with the matching rule name.
package portpolicy

import (
	"fmt"
	"sync"
)

// Action represents the policy decision for a port.
type Action int

const (
	Allow Action = iota
	Deny
)

func (a Action) String() string {
	if a == Allow {
		return "allow"
	}
	return "deny"
}

// Rule is a single policy entry.
type Rule struct {
	Name   string
	Action Action
	Min    int
	Max    int
}

// Verdict is the result of evaluating a port against the policy.
type Verdict struct {
	Action    Action
	MatchedBy string // name of the rule that matched, or "default"
}

// Policy holds an ordered list of rules and a default action.
type Policy struct {
	mu            sync.RWMutex
	rules         []Rule
	defaultAction Action
}

// New creates a Policy with the given default action.
func New(defaultAction Action) *Policy {
	return &Policy{defaultAction: defaultAction}
}

// AddRule appends a named rule to the policy.
func (p *Policy) AddRule(r Rule) error {
	if r.Min < 1 || r.Max > 65535 || r.Min > r.Max {
		return fmt.Errorf("portpolicy: invalid range %d-%d for rule %q", r.Min, r.Max, r.Name)
	}
	if r.Name == "" {
		return fmt.Errorf("portpolicy: rule name must not be empty")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.rules = append(p.rules, r)
	return nil
}

// Evaluate returns a Verdict for the given port number.
func (p *Policy) Evaluate(port int) Verdict {
	p.mu.RLock()
	defer p.mu.RUnlock()
	for _, r := range p.rules {
		if port >= r.Min && port <= r.Max {
			return Verdict{Action: r.Action, MatchedBy: r.Name}
		}
	}
	return Verdict{Action: p.defaultAction, MatchedBy: "default"}
}

// Rules returns a snapshot of the current rule list.
func (p *Policy) Rules() []Rule {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]Rule, len(p.rules))
	copy(out, p.rules)
	return out
}
