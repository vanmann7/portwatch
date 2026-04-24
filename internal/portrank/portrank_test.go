package portrank

import (
	"testing"
)

func TestBuiltinSSHRank(t *testing.T) {
	r := New(nil)
	if got := r.Rank(22); got != 90 {
		t.Fatalf("expected 90, got %d", got)
	}
}

func TestBuiltinTelnetRank(t *testing.T) {
	r := New(nil)
	if got := r.Rank(23); got != 95 {
		t.Fatalf("expected 95, got %d", got)
	}
}

func TestUnknownPortFallback(t *testing.T) {
	r := New(nil)
	if got := r.Rank(9999); got != defaultRank {
		t.Fatalf("expected %d, got %d", defaultRank, got)
	}
}

func TestOverrideTakesPrecedence(t *testing.T) {
	r := New(map[int]int{22: 50})
	if got := r.Rank(22); got != 50 {
		t.Fatalf("expected 50, got %d", got)
	}
}

func TestInvalidOverrideIgnored(t *testing.T) {
	r := New(map[int]int{80: 0, 443: 101})
	if got := r.Rank(80); got != 60 {
		t.Fatalf("expected builtin 60, got %d", got)
	}
	if got := r.Rank(443); got != 65 {
		t.Fatalf("expected builtin 65, got %d", got)
	}
}

func TestSetOverrideRuntime(t *testing.T) {
	r := New(nil)
	r.SetOverride(9000, 77)
	if got := r.Rank(9000); got != 77 {
		t.Fatalf("expected 77, got %d", got)
	}
}

func TestSetOverrideInvalidIgnored(t *testing.T) {
	r := New(nil)
	r.SetOverride(9000, 0)
	if got := r.Rank(9000); got != defaultRank {
		t.Fatalf("expected %d, got %d", defaultRank, got)
	}
}

func TestRemoveOverrideFallsBack(t *testing.T) {
	r := New(map[int]int{22: 50})
	r.RemoveOverride(22)
	if got := r.Rank(22); got != 90 {
		t.Fatalf("expected builtin 90 after removal, got %d", got)
	}
}

func TestTopReturnsHighestRanked(t *testing.T) {
	r := New(nil)
	ports := []int{8080, 22, 9999, 23, 443}
	top := r.Top(ports, 3)
	if len(top) != 3 {
		t.Fatalf("expected 3 results, got %d", len(top))
	}
	// expected order: 23(95), 22(90), 443(65)
	expected := []int{23, 22, 443}
	for i, p := range expected {
		if top[i] != p {
			t.Errorf("top[%d]: expected %d, got %d", i, p, top[i])
		}
	}
}

func TestTopNLargerThanSlice(t *testing.T) {
	r := New(nil)
	ports := []int{22, 80}
	top := r.Top(ports, 10)
	if len(top) != 2 {
		t.Fatalf("expected 2, got %d", len(top))
	}
}

func TestTopEmptySlice(t *testing.T) {
	r := New(nil)
	top := r.Top(nil, 5)
	if len(top) != 0 {
		t.Fatalf("expected empty, got %v", top)
	}
}

func TestTopTieBrokenByLowerPort(t *testing.T) {
	// Give two ports the same rank via overrides
	r := New(map[int]int{100: 50, 200: 50})
	top := r.Top([]int{200, 100}, 2)
	if top[0] != 100 {
		t.Errorf("expected lower port 100 first, got %d", top[0])
	}
}
