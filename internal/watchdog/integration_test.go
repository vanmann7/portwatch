package watchdog_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/watchdog"
)

// TestWatchdogFullLifecycle verifies the full healthy→stale→dead progression.
func TestWatchdogFullLifecycle(t *testing.T) {
	var stages []watchdog.Status
	var mu sync.Mutex
	record := func(s watchdog.Status) {
		mu.Lock()
		stages = append(stages, s)
		mu.Unlock()
	}

	w := watchdog.New(
		40*time.Millisecond,
		90*time.Millisecond,
		func() { record(watchdog.StatusStale) },
		func() { record(watchdog.StatusDead) },
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go w.Run(ctx, 10*time.Millisecond)

	// Initially healthy
	if w.Status() != watchdog.StatusHealthy {
		t.Fatal("should start healthy")
	}

	// Wait past stale threshold
	time.Sleep(70*time.Millisecond)
	if w.Status() != watchdog.StatusStale {
		t.Fatalf("expected stale, got %v", w.Status())
	}

	// Wait past dead threshold
	time.Sleep(60*time.Millisecond)
	if w.Status() != watchdog.StatusDead {
		t.Fatalf("expected dead, got %v", w.Status())
	}

	// Give callbacks time to fire
	time.Sleep(20*time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if len(stages) < 2 {
		t.Fatalf("expected at least 2 stage callbacks, got %d", len(stages))
	}
}
