package filter

import (
	"testing"
)

func TestAllowSinglePort(t *testing.T) {
	f, err := New([]string{"22"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Allow(22) {
		t.Error("expected port 22 to be allowed")
	}
	if f.Allow(80) {
		t.Error("expected port 80 to be denied")
	}
}

func TestAllowRange(t *testing.T) {
	f, err := New([]string{"8000-8080"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, p := range []int{8000, 8040, 8080} {
		if !f.Allow(p) {
			t.Errorf("expected port %d to be allowed", p)
		}
	}
	if f.Allow(7999) || f.Allow(8081) {
		t.Error("ports outside range should be denied")
	}
}

func TestAllowNoRules(t *testing.T) {
	f, err := New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, p := range []int{1, 80, 443, 65535} {
		if !f.Allow(p) {
			t.Errorf("expected port %d to be allowed with no rules", p)
		}
	}
}

func TestAllowMultipleRules(t *testing.T) {
	f, err := New([]string{"22", "80", "443", "8000-8100"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, p := range []int{22, 80, 443, 8050} {
		if !f.Allow(p) {
			t.Errorf("expected port %d to be allowed", p)
		}
	}
	if f.Allow(23) || f.Allow(8101) {
		t.Error("unexpected ports allowed")
	}
}

func TestInvalidRules(t *testing.T) {
	cases := []string{"0", "65536", "9000-8000", "abc", ""}
	for _, c := range cases {
		if _, err := New([]string{c}); err == nil {
			t.Errorf("expected error for rule %q", c)
		}
	}
}

func TestAllowRangeBoundaries(t *testing.T) {
	// Verify that range endpoints 1 and 65535 (min/max valid ports) are handled correctly.
	f, err := New([]string{"1-65535"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, p := range []int{1, 1024, 65535} {
		if !f.Allow(p) {
			t.Errorf("expected port %d to be allowed", p)
		}
	}
}
