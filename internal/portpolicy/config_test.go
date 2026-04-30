package portpolicy_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/portpolicy"
)

const sampleConfig = `
# portpolicy rules
deny  block-telnet  23
allow allow-web     80-443
deny  block-high    1024-65535
`

func TestParseRulesValid(t *testing.T) {
	rules, err := portpolicy.ParseRules(strings.NewReader(sampleConfig))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(rules))
	}
	if rules[0].Name != "block-telnet" || rules[0].Action != portpolicy.Deny {
		t.Errorf("rule 0 mismatch: %+v", rules[0])
	}
	if rules[1].Min != 80 || rules[1].Max != 443 {
		t.Errorf("rule 1 range mismatch: %+v", rules[1])
	}
}

func TestParseRulesSkipsBlankAndComments(t *testing.T) {
	input := "\n# comment\n\n"
	rules, err := portpolicy.ParseRules(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 0 {
		t.Fatalf("expected 0 rules, got %d", len(rules))
	}
}

func TestParseRulesMalformedLine(t *testing.T) {
	_, err := portpolicy.ParseRules(strings.NewReader("deny only-two-fields"))
	if err == nil {
		t.Fatal("expected error for malformed line")
	}
}

func TestParseRulesUnknownAction(t *testing.T) {
	_, err := portpolicy.ParseRules(strings.NewReader("permit ssh 22"))
	if err == nil {
		t.Fatal("expected error for unknown action")
	}
}

func TestParseRulesInvalidPort(t *testing.T) {
	_, err := portpolicy.ParseRules(strings.NewReader("deny bad-port abc"))
	if err == nil {
		t.Fatal("expected error for non-numeric port")
	}
}

func TestParseRulesRoundTrip(t *testing.T) {
	rules, err := portpolicy.ParseRules(strings.NewReader(sampleConfig))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	p := portpolicy.New(portpolicy.Allow)
	for _, r := range rules {
		if err := p.AddRule(r); err != nil {
			t.Fatalf("AddRule: %v", err)
		}
	}
	v := p.Evaluate(23)
	if v.Action != portpolicy.Deny || v.MatchedBy != "block-telnet" {
		t.Errorf("unexpected verdict for port 23: %+v", v)
	}
}
