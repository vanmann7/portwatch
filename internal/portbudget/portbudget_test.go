package portbudget

import (
	"testing"
)

func TestBudgetNotExceededWhenEmpty(t *testing.T) {
	b, err := New(WithRange(80, 89, 3))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	exceeded, _ := b.Exceeded()
	if exceeded {
		t.Fatal("expected budget not exceeded for empty open set")
	}
}

func TestBudgetNotExceededAtLimit(t *testing.T) {
	b, _ := New(WithRange(80, 89, 2))
	b.Open(80)
	b.Open(81)
	exceeded, _ := b.Exceeded()
	if exceeded {
		t.Fatal("expected budget not exceeded when exactly at limit")
	}
}

func TestBudgetExceededOverLimit(t *testing.T) {
	b, _ := New(WithRange(80, 89, 2))
	b.Open(80)
	b.Open(81)
	b.Open(82)
	exceeded, desc := b.Exceeded()
	if !exceeded {
		t.Fatal("expected budget exceeded")
	}
	if desc == "" {
		t.Fatal("expected non-empty description")
	}
}

func TestBudgetCloseReducesCount(t *testing.T) {
	b, _ := New(WithRange(80, 89, 2))
	b.Open(80)
	b.Open(81)
	b.Open(82)
	b.Close(82)
	exceeded, _ := b.Exceeded()
	if exceeded {
		t.Fatal("expected budget not exceeded after close")
	}
}

func TestBudgetPortOutsideRangeIgnored(t *testing.T) {
	b, _ := New(WithRange(80, 89, 1))
	b.Open(443) // outside range
	exceeded, _ := b.Exceeded()
	if exceeded {
		t.Fatal("port outside range should not affect budget")
	}
}

func TestBudgetCountReturnsCorrectValue(t *testing.T) {
	b, _ := New()
	b.Open(8080)
	b.Open(8081)
	b.Open(9000)
	if got := b.Count(8080, 8090); got != 2 {
		t.Fatalf("expected count 2, got %d", got)
	}
}

func TestBudgetMultipleRanges(t *testing.T) {
	b, _ := New(
		WithRange(80, 89, 2),
		WithRange(443, 449, 1),
	)
	b.Open(80)
	b.Open(443)
	b.Open(444) // exceeds second range
	exceeded, desc := b.Exceeded()
	if !exceeded {
		t.Fatal("expected second range to be exceeded")
	}
	if desc == "" {
		t.Fatal("expected description for exceeded range")
	}
}

func TestInvalidRangeReturnsError(t *testing.T) {
	_, err := New(WithRange(500, 100, 5))
	if err == nil {
		t.Fatal("expected error for invalid range lo > hi")
	}
}

func TestInvalidMaxReturnsError(t *testing.T) {
	_, err := New(WithRange(80, 90, 0))
	if err == nil {
		t.Fatal("expected error for max < 1")
	}
}
