package portguard_test

import (
	"strings"
	"testing"

	"github.com/example/portwatch/internal/portguard"
)

func TestParseRulesValid(t *testing.T) {
	input := `
# guard rules
allow 80
deny  22
warn  1024-49151
`
	rules, err := portguard.ParseRules(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(rules))
	}
	if rules[0].Verdict != portguard.Allow || rules[0].Low != 80 {
		t.Errorf("rule 0 mismatch: %+v", rules[0])
	}
	if rules[1].Verdict != portguard.Deny || rules[1].Low != 22 {
		t.Errorf("rule 1 mismatch: %+v", rules[1])
	}
	if rules[2].Verdict != portguard.Warn || rules[2].Low != 1024 || rules[2].High != 49151 {
		t.Errorf("rule 2 mismatch: %+v", rules[2])
	}
}

func TestParseRulesSkipsBlankAndComments(t *testing.T) {
	input := "\n# comment\n\nallow 443\n"
	rules, err := portguard.ParseRules(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
}

func TestParseRulesMalformedLine(t *testing.T) {
	_, err := portguard.ParseRules(strings.NewReader("allow\n"))
	if err == nil {
		t.Fatal("expected error for malformed line")
	}
}

func TestParseRulesUnknownVerdict(t *testing.T) {
	_, err := portguard.ParseRules(strings.NewReader("block 80\n"))
	if err == nil {
		t.Fatal("expected error for unknown verdict")
	}
}

func TestParseRulesInvalidPort(t *testing.T) {
	_, err := portguard.ParseRules(strings.NewReader("allow abc\n"))
	if err == nil {
		t.Fatal("expected error for invalid port")
	}
}

func TestParseRulesRoundTrip(t *testing.T) {
	input := "deny 22\nallow 80-443\n"
	rules, err := portguard.ParseRules(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	g := portguard.New(rules, portguard.Allow)
	if v := g.Evaluate(22); v != portguard.Deny {
		t.Fatalf("port 22: expected Deny, got %v", v)
	}
	if v := g.Evaluate(80); v != portguard.Allow {
		t.Fatalf("port 80: expected Allow, got %v", v)
	}
}
