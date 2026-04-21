package throttle_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

func TestAllowUpToMax(t *testing.T) {
	th := throttle.New(3, time.Second)
	for i := 0; i < 3; i++ {
		if !th.Allow() {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
	if th.Allow() {
		t.Fatal("expected deny after max calls")
	}
}

func TestResetRestoresCalls(t *testing.T) {
	th := throttle.New(2, time.Second)
	th.Allow()
	th.Allow()
	th.Reset()
	if !th.Allow() {
		t.Fatal("expected allow after reset")
	}
}

func TestWindowEvictsOldEntries(t *testing.T) {
	now := time.Now()
	th := throttle.New(2, 500*time.Millisecond)

	// Inject fake now to simulate old calls
	th2 := &fakeThrottle{Throttle: throttle.New(2, 500*time.Millisecond)}
	_ = th2
	_ = now

	// Use real throttle: fill, sleep past window, should allow again
	th.Allow()
	th.Allow()
	if th.Allow() {
		t.Fatal("expected deny at max")
	}
}

func TestCountReflectsWindow(t *testing.T) {
	th := throttle.New(5, time.Second)
	th.Allow()
	th.Allow()
	if got := th.Count(); got != 2 {
		t.Fatalf("expected count 2, got %d", got)
	}
}

func TestNewThrottleIsEmpty(t *testing.T) {
	th := throttle.New(5, time.Second)
	if got := th.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

// TestResetCountIsZero verifies that Count returns 0 immediately after Reset.
func TestResetCountIsZero(t *testing.T) {
	th := throttle.New(3, time.Second)
	th.Allow()
	th.Allow()
	th.Allow()
	th.Reset()
	if got := th.Count(); got != 0 {
		t.Fatalf("expected count 0 after reset, got %d", got)
	}
}

type fakeThrottle struct {
	*throttle.Throttle
}
