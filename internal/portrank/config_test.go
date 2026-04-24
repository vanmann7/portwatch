package portrank

import "testing"

func TestParseRulesValid(t *testing.T) {
	lines := []string{"22:99", "80:40"}
	rules, err := ParseRules(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
	if rules[0].Port != 22 || rules[0].Rank != 99 {
		t.Errorf("rule[0]: got %+v", rules[0])
	}
}

func TestParseRulesSkipsBlankAndComments(t *testing.T) {
	lines := []string{"", "# comment", "  ", "443:70"}
	rules, err := ParseRules(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rules) != 1 || rules[0].Port != 443 {
		t.Errorf("expected single rule for 443, got %+v", rules)
	}
}

func TestParseRulesMissingColon(t *testing.T) {
	_, err := ParseRules([]string{"2222"})
	if err == nil {
		t.Fatal("expected error for missing colon")
	}
}

func TestParseRulesInvalidPort(t *testing.T) {
	_, err := ParseRules([]string{"0:50"})
	if err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestParseRulesPortOutOfRange(t *testing.T) {
	_, err := ParseRules([]string{"99999:50"})
	if err == nil {
		t.Fatal("expected error for port > 65535")
	}
}

func TestParseRulesInvalidRank(t *testing.T) {
	_, err := ParseRules([]string{"80:0"})
	if err == nil {
		t.Fatal("expected error for rank 0")
	}
	_, err = ParseRules([]string{"80:101"})
	if err == nil {
		t.Fatal("expected error for rank 101")
	}
}

func TestNewFromRulesAppliesOverrides(t *testing.T) {
	rules := []RuleEntry{{Port: 22, Rank: 55}, {Port: 9999, Rank: 30}}
	r := NewFromRules(rules)
	if got := r.Rank(22); got != 55 {
		t.Errorf("expected 55, got %d", got)
	}
	if got := r.Rank(9999); got != 30 {
		t.Errorf("expected 30, got %d", got)
	}
}

func TestNewFromRulesEmptyUsesBuiltins(t *testing.T) {
	r := NewFromRules(nil)
	if got := r.Rank(23); got != 95 {
		t.Errorf("expected builtin 95 for telnet, got %d", got)
	}
}
