package portgroup_test

import (
	"testing"

	"github.com/user/portwatch/internal/portgroup"
)

func TestLookupSinglePort(t *testing.T) {
	g, err := portgroup.New([]portgroup.Rule{
		{Name: "web", Ports: []string{"80", "443"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := g.Lookup(80); got != "web" {
		t.Errorf("Lookup(80) = %q, want %q", got, "web")
	}
	if got := g.Lookup(443); got != "web" {
		t.Errorf("Lookup(443) = %q, want %q", got, "web")
	}
}

func TestLookupRange(t *testing.T) {
	g, err := portgroup.New([]portgroup.Rule{
		{Name: "ephemeral", Ports: []string{"32768-32770"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, p := range []int{32768, 32769, 32770} {
		if got := g.Lookup(p); got != "ephemeral" {
			t.Errorf("Lookup(%d) = %q, want %q", p, got, "ephemeral")
		}
	}
}

func TestLookupNoMatch(t *testing.T) {
	g, err := portgroup.New([]portgroup.Rule{
		{Name: "db", Ports: []string{"5432"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := g.Lookup(9999); got != "" {
		t.Errorf("Lookup(9999) = %q, want empty string", got)
	}
}

func TestGroupsReturnsNames(t *testing.T) {
	g, err := portgroup.New([]portgroup.Rule{
		{Name: "web", Ports: []string{"80"}},
		{Name: "db", Ports: []string{"5432"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	names := g.Groups()
	if len(names) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(names))
	}
}

func TestInvalidPortReturnsError(t *testing.T) {
	_, err := portgroup.New([]portgroup.Rule{
		{Name: "bad", Ports: []string{"not-a-port"}},
	})
	if err == nil {
		t.Error("expected error for invalid port, got nil")
	}
}

func TestInvalidRangeReturnsError(t *testing.T) {
	_, err := portgroup.New([]portgroup.Rule{
		{Name: "bad", Ports: []string{"9000-8000"}},
	})
	if err == nil {
		t.Error("expected error for inverted range, got nil")
	}
}

func TestFirstGroupWins(t *testing.T) {
	g, err := portgroup.New([]portgroup.Rule{
		{Name: "first", Ports: []string{"8080"}},
		{Name: "second", Ports: []string{"8080"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := g.Lookup(8080); got != "first" {
		t.Errorf("Lookup(8080) = %q, want %q", got, "first")
	}
}
