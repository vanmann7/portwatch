package circuitbreaker

import (
	"testing"
	"time"
)

func TestNewBreakerIsClosed(t *testing.T) {
	b := New(3, time.Second)
	if b.State() != StateClosed {
		t.Fatalf("expected StateClosed, got %v", b.State())
	}
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestBreakerOpensAfterMaxFailures(t *testing.T) {
	b := New(3, time.Second)
	b.RecordFailure()
	b.RecordFailure()
	if b.State() != StateClosed {
		t.Fatal("should still be closed after 2 failures")
	}
	b.RecordFailure()
	if b.State() != StateOpen {
		t.Fatalf("expected StateOpen after 3 failures, got %v", b.State())
	}
	if err := b.Allow(); err != ErrOpen {
		t.Fatalf("expected ErrOpen, got %v", err)
	}
}

func TestBreakerHalfOpenAfterReset(t *testing.T) {
	now := time.Now()
	b := New(1, 50*time.Millisecond)
	b.now = func() time.Time { return now }

	b.RecordFailure() // trips breaker
	if b.State() != StateOpen {
		t.Fatal("expected StateOpen")
	}

	// advance clock past reset window
	b.now = func() time.Time { return now.Add(100 * time.Millisecond) }
	if err := b.Allow(); err != nil {
		t.Fatalf("expected nil in half-open, got %v", err)
	}
	if b.State() != StateHalfOpen {
		t.Fatalf("expected StateHalfOpen, got %v", b.State())
	}
}

func TestSuccessClosesBreakerFromHalfOpen(t *testing.T) {
	now := time.Now()
	b := New(1, 10*time.Millisecond)
	b.now = func() time.Time { return now }

	b.RecordFailure()
	b.now = func() time.Time { return now.Add(20 * time.Millisecond) }
	_ = b.Allow() // transitions to half-open

	b.RecordSuccess()
	if b.State() != StateClosed {
		t.Fatalf("expected StateClosed after success, got %v", b.State())
	}
	if b.Failures() != 0 {
		t.Fatalf("expected 0 failures after success, got %d", b.Failures())
	}
}

func TestFailureInHalfOpenReopens(t *testing.T) {
	now := time.Now()
	b := New(1, 10*time.Millisecond)
	b.now = func() time.Time { return now }

	b.RecordFailure()
	b.now = func() time.Time { return now.Add(20 * time.Millisecond) }
	_ = b.Allow() // half-open

	b.RecordFailure() // fail again
	if b.State() != StateOpen {
		t.Fatalf("expected StateOpen after second failure, got %v", b.State())
	}
}

func TestFailuresCounterAccumulates(t *testing.T) {
	b := New(10, time.Second)
	for i := 0; i < 5; i++ {
		b.RecordFailure()
	}
	if b.Failures() != 5 {
		t.Fatalf("expected 5 failures, got %d", b.Failures())
	}
}
