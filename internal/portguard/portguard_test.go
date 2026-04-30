package portguard_test

import (
	"testing"

	"github.com/example/portwatch/internal/portguard"
)

func TestAllowByDefault(t *testing.T) {
	g := portguard.New(nil, portguard.Allow)
	if v := g.Evaluate(8080); v != portguard.Allow {
		t.Fatalf("expected Allow, got %v", v)
	}
}

func TestDenyByDefault(t *testing.T) {
	g := portguard.New(nil, portguard.Deny)
	if v := g.Evaluate(443); v != portguard.Deny {
		t.Fatalf("expected Deny, got %v", v)
	}
}

func TestRuleMatchesSinglePort(t *testing.T) {
	rules := []portguard.Rule{{Low: 22, High: 22, Verdict: portguard.Deny}}
	g := portguard.New(rules, portguard.Allow)
	if v := g.Evaluate(22); v != portguard.Deny {
		t.Fatalf("expected Deny for port 22, got %v", v)
	}
	if v := g.Evaluate(23); v != portguard.Allow {
		t.Fatalf("expected Allow for port 23, got %v", v)
	}
}

func TestRuleMatchesRange(t *testing.T) {
	rules := []portguard.Rule{{Low: 1024, High: 49151, Verdict: portguard.Warn}}
	g := portguard.New(rules, portguard.Allow)
	for _, p := range []int{1024, 8080, 49151} {
		if v := g.Evaluate(p); v != portguard.Warn {
			t.Fatalf("port %d: expected Warn, got %v", p, v)
		}
	}
}

func TestFirstRuleWins(t *testing.T) {
	rules := []portguard.Rule{
		{Low: 80, High: 80, Verdict: portguard.Allow},
		{Low: 1, High: 1024, Verdict: portguard.Deny},
	}
	g := portguard.New(rules, portguard.Deny)
	if v := g.Evaluate(80); v != portguard.Allow {
		t.Fatalf("expected Allow for port 80, got %v", v)
	}
}

func TestAddRuleValid(t *testing.T) {
	g := portguard.New(nil, portguard.Deny)
	if err := g.AddRule(portguard.Rule{Low: 443, High: 443, Verdict: portguard.Allow}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v := g.Evaluate(443); v != portguard.Allow {
		t.Fatalf("expected Allow, got %v", v)
	}
}

func TestAddRuleInvalidRange(t *testing.T) {
	g := portguard.New(nil, portguard.Allow)
	if err := g.AddRule(portguard.Rule{Low: 500, High: 100, Verdict: portguard.Deny}); err == nil {
		t.Fatal("expected error for inverted range")
	}
}

func TestRulesSnapshot(t *testing.T) {
	original := []portguard.Rule{
		{Low: 22, High: 22, Verdict: portguard.Deny},
		{Low: 80, High: 80, Verdict: portguard.Allow},
	}
	g := portguard.New(original, portguard.Warn)
	snap := g.Rules()
	if len(snap) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(snap))
	}
}
