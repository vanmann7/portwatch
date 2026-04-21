package labeler_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/labeler"
)

func TestParseRulesSinglePort(t *testing.T) {
	rules, err := labeler.ParseRules([]string{"22:ssh-custom"})
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 1 || rules[0].Low != 22 || rules[0].High != 22 || rules[0].Label != "ssh-custom" {
		t.Fatalf("unexpected rule: %+v", rules[0])
	}
}

func TestParseRulesRange(t *testing.T) {
	rules, err := labeler.ParseRules([]string{"8000-8100:app-tier"})
	if err != nil {
		t.Fatal(err)
	}
	if rules[0].Low != 8000 || rules[0].High != 8100 || rules[0].Label != "app-tier" {
		t.Fatalf("unexpected rule: %+v", rules[0])
	}
}

func TestParseRulesMultiple(t *testing.T) {
	entries := []string{"80:web", "443:tls", "3000-3100:dev"}
	rules, err := labeler.ParseRules(entries)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(rules))
	}
}

func TestParseRulesMissingLabel(t *testing.T) {
	_, err := labeler.ParseRules([]string{"80:"})
	if err == nil {
		t.Fatal("expected error for missing label")
	}
}

func TestParseRulesInvalidPort(t *testing.T) {
	_, err := labeler.ParseRules([]string{"abc:label"})
	if err == nil {
		t.Fatal("expected error for non-numeric port")
	}
}

func TestParseRulesOutOfBounds(t *testing.T) {
	_, err := labeler.ParseRules([]string{"0:zero"})
	if err == nil {
		t.Fatal("expected error for port 0")
	}
	_, err = labeler.ParseRules([]string{"65536:toobig"})
	if err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestParseRulesInvertedRange(t *testing.T) {
	_, err := labeler.ParseRules([]string{"9000-8000:bad"})
	if err == nil {
		t.Fatal("expected error for inverted range")
	}
}
