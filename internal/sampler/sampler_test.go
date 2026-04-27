package sampler_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/sampler"
)

const tickInterval = 20 * time.Millisecond

func TestSamplerEmitsResult(t *testing.T) {
	called := make(chan struct{}, 1)
	scan := func(_ context.Context) ([]int, error) {
		called <- struct{}{}
		return []int{8080, 9090}, nil
	}

	s := sampler.New(tickInterval, 0, scan)
	ctx, cancel := context.WithTimeout(context.Background(), 5*tickInterval)
	defer cancel()

	ch := s.Run(ctx)
	select {
	case r := <-ch:
		if r.Err != nil {
			t.Fatalf("unexpected error: %v", r.Err)
		}
		if len(r.Ports) != 2 {
			t.Fatalf("expected 2 ports, got %d", len(r.Ports))
		}
		if r.At.IsZero() {
			t.Fatal("expected non-zero timestamp")
		}
	case <-time.After(10 * tickInterval):
		t.Fatal("timed out waiting for first result")
	}
}

func TestSamplerPropagatesError(t *testing.T) {
	sentinel := errors.New("scan failed")
	scan := func(_ context.Context) ([]int, error) {
		return nil, sentinel
	}

	s := sampler.New(tickInterval, 0, scan)
	ctx, cancel := context.WithTimeout(context.Background(), 5*tickInterval)
	defer cancel()

	r := <-s.Run(ctx)
	if !errors.Is(r.Err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", r.Err)
	}
}

func TestSamplerClosesChannelOnCancel(t *testing.T) {
	scan := func(_ context.Context) ([]int, error) { return nil, nil }
	s := sampler.New(tickInterval, 0, scan)

	ctx, cancel := context.WithCancel(context.Background())
	ch := s.Run(ctx)
	cancel()

	// Drain until closed; must not block.
	timer := time.NewTimer(10 * tickInterval)
	defer timer.Stop()
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return // channel closed as expected
			}
		case <-timer.C:
			t.Fatal("channel was not closed after context cancellation")
		}
	}
}

func TestSamplerTicksMultipleTimes(t *testing.T) {
	var count atomic.Int32
	scan := func(_ context.Context) ([]int, error) {
		count.Add(1)
		return nil, nil
	}

	s := sampler.New(tickInterval, 0, scan)
	ctx, cancel := context.WithTimeout(context.Background(), 6*tickInterval)
	defer cancel()

	ch := s.Run(ctx)
	for range ch {
	}

	if got := count.Load(); got < 2 {
		t.Fatalf("expected at least 2 ticks, got %d", got)
	}
}

func TestSamplerResultTimestampIsRecent(t *testing.T) {
	scan := func(_ context.Context) ([]int, error) {
		return []int{443}, nil
	}

	before := time.Now()
	s := sampler.New(tickInterval, 0, scan)
	ctx, cancel := context.WithTimeout(context.Background(), 5*tickInterval)
	defer cancel()

	r := <-s.Run(ctx)
	after := time.Now()

	if r.At.Before(before) || r.At.After(after) {
		t.Fatalf("timestamp %v is outside expected range [%v, %v]", r.At, before, after)
	}
}
