package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

var errTemp = errors.New("temporary error")

func TestSuccessOnFirstAttempt(t *testing.T) {
	r := New(Default())
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRetriesUpToMax(t *testing.T) {
	p := Policy{MaxAttempts: 3, Delay: time.Millisecond, Multiplier: 1.0}
	r := New(p)
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		return errTemp
	})
	if !errors.Is(err, errTemp) {
		t.Fatalf("expected errTemp, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestSucceedsOnSecondAttempt(t *testing.T) {
	p := Policy{MaxAttempts: 3, Delay: time.Millisecond, Multiplier: 1.0}
	r := New(p)
	calls := 0
	err := r.Do(context.Background(), func() error {
		calls++
		if calls < 2 {
			return errTemp
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

func TestContextCancellation(t *testing.T) {
	p := Policy{MaxAttempts: 5, Delay: 50 * time.Millisecond, Multiplier: 1.0}
	r := New(p)
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	go func() {
		time.Sleep(60 * time.Millisecond)
		cancel()
	}()
	err := r.Do(ctx, func() error {
		calls++
		return errTemp
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if calls >= 5 {
		t.Fatal("expected early exit due to cancellation")
	}
}
