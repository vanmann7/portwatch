package dispatcher_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/dispatcher"
)

func TestDispatcherCallsAllHandlers(t *testing.T) {
	d := dispatcher.New[int]()

	var countA, countB atomic.Int64
	d.Register(func(_ context.Context, _ int) error { countA.Add(1); return nil })
	d.Register(func(_ context.Context, _ int) error { countB.Add(1); return nil })

	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)

	ctx := context.Background()
	errCh := d.Dispatch(ctx, ch)
	drain(t, errCh)

	if countA.Load() != 3 || countB.Load() != 3 {
		t.Fatalf("expected 3 calls each, got %d and %d", countA.Load(), countB.Load())
	}
}

func TestDispatcherCollectsHandlerErrors(t *testing.T) {
	d := dispatcher.New[string]()
	boom := errors.New("boom")
	d.Register(func(_ context.Context, _ string) error { return boom })

	ch := make(chan string, 1)
	ch <- "event"
	close(ch)

	var errs []error
	for err := range d.Dispatch(context.Background(), ch) {
		errs = append(errs, err)
	}

	if len(errs) == 0 {
		t.Fatal("expected at least one error")
	}
	if !errors.Is(errs[0], boom) {
		t.Fatalf("unexpected error: %v", errs[0])
	}
}

func TestDispatcherStopsOnContextCancel(t *testing.T) {
	d := dispatcher.New[int]()
	var calls atomic.Int64
	d.Register(func(_ context.Context, _ int) error { calls.Add(1); return nil })

	ch := make(chan int) // never closed
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	drain(t, d.Dispatch(ctx, ch))

	if calls.Load() != 0 {
		t.Fatalf("expected 0 calls, got %d", calls.Load())
	}
}

func TestDispatcherLen(t *testing.T) {
	d := dispatcher.New[int]()
	if d.Len() != 0 {
		t.Fatalf("expected 0, got %d", d.Len())
	}
	d.Register(func(_ context.Context, _ int) error { return nil })
	d.Register(func(_ context.Context, _ int) error { return nil })
	if d.Len() != 2 {
		t.Fatalf("expected 2, got %d", d.Len())
	}
}

func TestDispatcherNoHandlersNoErrors(t *testing.T) {
	d := dispatcher.New[int]()
	ch := make(chan int, 2)
	ch <- 10
	ch <- 20
	close(ch)

	var count int
	for range d.Dispatch(context.Background(), ch) {
		count++
	}
	if count != 0 {
		t.Fatalf("expected 0 errors, got %d", count)
	}
}

func drain(t *testing.T, ch <-chan error) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		for range ch {
		}
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for dispatcher to finish")
	}
}
