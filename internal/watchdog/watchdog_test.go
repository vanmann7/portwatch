package watchdog_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watchdog"
)

func TestHealthyAfterPing(t *testing.T) {
	w := watchdog.New(100*time.Millisecond, 200*time.Millisecond, nil, nil)
	w.Ping()
	if got := w.Status(); got != watchdog.StatusHealthy {
		t.Fatalf("expected Healthy, got %v", got)
	}
}

func TestStaleAfterTimeout(t *testing.T) {
	var called int32
	w := watchdog.New(50*time.Millisecond, 200*time.Millisecond, func() {
		atomic.StoreInt32(&called, 1)
	}, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	go w.Run(ctx, 20*time.Millisecond)

	time.Sleep(120*time.Millisecond)
	if got := w.Status(); got != watchdog.StatusStale {
		t.Fatalf("expected Stale, got %v", got)
	}
	time.Sleep(30*time.Millisecond)
	if atomic.LoadInt32(&called) != 1 {
		t.Fatal("expected onStale callback to have been called")
	}
}

func TestDeadAfterTimeout(t *testing.T) {
	var deadCalled int32
	w := watchdog.New(20*time.Millisecond, 60*time.Millisecond, nil, func() {
		atomic.StoreInt32(&deadCalled, 1)
	})

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	go w.Run(ctx, 10*time.Millisecond)

	time.Sleep(120*time.Millisecond)
	if got := w.Status(); got != watchdog.StatusDead {
		t.Fatalf("expected Dead, got %v", got)
	}
	time.Sleep(20*time.Millisecond)
	if atomic.LoadInt32(&deadCalled) != 1 {
		t.Fatal("expected onDead callback to have been called")
	}
}

func TestPingResetsToHealthy(t *testing.T) {
	w := watchdog.New(30*time.Millisecond, 100*time.Millisecond, nil, nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go w.Run(ctx, 10*time.Millisecond)

	time.Sleep(50*time.Millisecond)
	w.Ping()
	time.Sleep(10*time.Millisecond)

	if got := w.Status(); got != watchdog.StatusHealthy {
		t.Fatalf("expected Healthy after ping, got %v", got)
	}
}

func TestRunStopsOnContextCancel(t *testing.T) {
	w := watchdog.New(time.Second, 2*time.Second, nil, nil)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		w.Run(ctx, 10*time.Millisecond)
		close(done)
	}()
	cancel()
	select {
	case <-done:
	case <-time.After(200*time.Millisecond):
		t.Fatal("Run did not stop after context cancellation")
	}
}
