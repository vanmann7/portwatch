package portpolicy_test

import (
	"testing"

	"github.com/user/portwatch/internal/portpolicy"
)

func TestAllowByDefault(t *testing.T) {
	p := portpolicy.New(portpolicy.Allow)
	v := p.Evaluate(9999)
	if v.Action != portpolicy.Allow {
		t.Fatalf("expected Allow, got %s", v.Action)
	}
	if v.MatchedBy != "default" {
		t.Fatalf("expected default, got %q", v.MatchedBy)
	}
}

func TestDenyByDefault(t *testing.T) {
	p := portpolicy.New(portpolicy.Deny)
	v := p.Evaluate(1234)
	if v.Action != portpolicy.Deny {
		t.Fatalf("expected Deny, got %s", v.Action)
	}
}

func TestRuleMatchesSinglePort(t *testing.T) {
	p := portpolicy.New(portpolicy.Allow)
	_ = p.AddRule(portpolicy.Rule{Name: "block-telnet", Action: portpolicy.Deny, Min: 23, Max: 23})
	v := p.Evaluate(23)
	if v.Action != portpolicy.Deny {
		t.Fatalf("expected Deny, got %s", v.Action)
	}
	if v.MatchedBy != "block-telnet" {
		t.Fatalf("expected block-telnet, got %q", v.MatchedBy)
	}
}

func TestRuleMatchesRange(t *testing.T) {
	p := portpolicy.New(portpolicy.Deny)
	_ = p.AddRule(portpolicy.Rule{Name: "allow-web", Action: portpolicy.Allow, Min: 80, Max: 443})
	for _, port := range []int{80, 200, 443} {
		v := p.Evaluate(port)
		if v.Action != portpolicy.Allow {
			t.Errorf("port %d: expected Allow, got %s", port, v.Action)
		}
	}
}

func TestFirstRuleWins(t *testing.T) {
	p := portpolicy.New(portpolicy.Allow)
	_ = p.AddRule(portpolicy.Rule{Name: "deny-ssh", Action: portpolicy.Deny, Min: 22, Max: 22})
	_ = p.AddRule(portpolicy.Rule{Name: "allow-ssh", Action: portpolicy.Allow, Min: 22, Max: 22})
	v := p.Evaluate(22)
	if v.MatchedBy != "deny-ssh" {
		t.Fatalf("expected first rule to win, got %q", v.MatchedBy)
	}
}

func TestInvalidRuleReturnsError(t *testing.T) {
	p := portpolicy.New(portpolicy.Allow)
	err := p.AddRule(portpolicy.Rule{Name: "bad", Action: portpolicy.Deny, Min: 500, Max: 100})
	if err == nil {
		t.Fatal("expected error for inverted range")
	}
}

func TestEmptyNameReturnsError(t *testing.T) {
	p := portpolicy.New(portpolicy.Allow)
	err := p.AddRule(portpolicy.Rule{Name: "", Action: portpolicy.Deny, Min: 22, Max: 22})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRulesSnapshot(t *testing.T) {
	p := portpolicy.New(portpolicy.Allow)
	_ = p.AddRule(portpolicy.Rule{Name: "r1", Action: portpolicy.Deny, Min: 22, Max: 22})
	_ = p.AddRule(portpolicy.Rule{Name: "r2", Action: portpolicy.Allow, Min: 80, Max: 80})
	rules := p.Rules()
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
}
