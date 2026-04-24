package portmatcher_test

import (
	"testing"

	"github.com/user/portwatch/internal/portmatcher"
)

func TestMatchSinglePort(t *testing.T) {
	m, err := portmatcher.New([]string{"80"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Match(80) {
		t.Error("expected port 80 to match")
	}
	if m.Match(81) {
		t.Error("expected port 81 not to match")
	}
}

func TestMatchRange(t *testing.T) {
	m, err := portmatcher.New([]string{"1000-2000"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, p := range []int{1000, 1500, 2000} {
		if !m.Match(p) {
			t.Errorf("expected port %d to match range 1000-2000", p)
		}
	}
	for _, p := range []int{999, 2001} {
		if m.Match(p) {
			t.Errorf("expected port %d not to match range 1000-2000", p)
		}
	}
}

func TestMatchCommaList(t *testing.T) {
	m, err := portmatcher.New([]string{"22,80,443"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, p := range []int{22, 80, 443} {
		if !m.Match(p) {
			t.Errorf("expected port %d to match", p)
		}
	}
	if m.Match(8080) {
		t.Error("expected port 8080 not to match")
	}
}

func TestMatchMultipleEntries(t *testing.T) {
	m, err := portmatcher.New([]string{"22", "8000-8080"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.Match(22) || !m.Match(8042) {
		t.Error("expected matches for 22 and 8042")
	}
}

func TestMatchNoRules(t *testing.T) {
	m, err := portmatcher.New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Match(80) {
		t.Error("expected no match when no rules defined")
	}
	if m.Len() != 0 {
		t.Errorf("expected Len 0, got %d", m.Len())
	}
}

func TestInvalidPort(t *testing.T) {
	cases := []string{"0", "65536", "abc", "80-70", "-1"}
	for _, c := range cases {
		_, err := portmatcher.New([]string{c})
		if err == nil {
			t.Errorf("expected error for rule %q", c)
		}
	}
}

func TestLenReflectsRuleCount(t *testing.T) {
	m, err := portmatcher.New([]string{"22", "80", "443"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Len() != 3 {
		t.Errorf("expected Len 3, got %d", m.Len())
	}
}
