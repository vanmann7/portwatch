package snapshot_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

func TestSchedulerEmitsImmediately(t *testing.T) {
	called := 0
	scan := func(ctx context.Context) ([]int, error) {
		called++
		return []int{80, 443}, nil
	}

	s := snapshot.New(scan, 10*time.Second)
	ctx, cancel := context.WithCancel(context.Background())

	go s.Run(ctx)

	select {
	case r := <-s.Results():
		if r.Err != nil {
			t.Fatalf("unexpected error: %v", r.Err)
		}
		if len(r.Ports) != 2 {
			t.Fatalf("expected 2 ports, got %d", len(r.Ports))
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for first result")
	}
	cancel()
}

func TestSchedulerPropagatesError(t *testing.T) {
	scanErr := errors.New("scan failed")
	scan := func(ctx context.Context) ([]int, error) {
		return nil, scanErr
	}

	s := snapshot.New(scan, 10*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go s.Run(ctx)

	select {
	case r := <-s.Results():
		if !errors.Is(r.Err, scanErr) {
			t.Fatalf("expected scan error, got %v", r.Err)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for result")
	}
}

func TestSchedulerClosesChannelOnCancel(t *testing.T) {
	scan := func(ctx context.Context) ([]int, error) {
		return nil, nil
	}

	s := snapshot.New(scan, 50*time.Millisecond)
	ctx, cancel := context.WithCancel(context.Background())

	go s.Run(ctx)
	// drain the first immediate result
	<-s.Results()
	cancel()

	// channel should close shortly after cancel
	select {
	case _, ok := <-s.Results():
		if ok {
			// drain extra ticks that may have fired before cancel
			for range s.Results() {
			}
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("channel not closed after cancel")
	}
}
