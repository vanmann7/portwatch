package backoff_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/backoff"
)

func TestDefaultValues(t *testing.T) {
	s := backoff.Default()
	if s.Base != 250*time.Millisecond {
		t.Errorf("expected Base=250ms, got %v", s.Base)
	}
	if s.Max != 30*time.Second {
		t.Errorf("expected Max=30s, got %v", s.Max)
	}
	if s.Factor != 2.0 {
		t.Errorf("expected Factor=2.0, got %v", s.Factor)
	}
}

func TestDelayGrowsExponentially(t *testing.T) {
	s := backoff.Strategy{
		Base:   100 * time.Millisecond,
		Max:    10 * time.Second,
		Factor: 2.0,
		Jitter: false,
	}

	expected := []time.Duration{
		100 * time.Millisecond,
		200 * time.Millisecond,
		400 * time.Millisecond,
		800 * time.Millisecond,
	}

	for attempt, want := range expected {
		got := s.Delay(attempt)
		if got != want {
			t.Errorf("attempt %d: expected %v, got %v", attempt, want, got)
		}
	}
}

func TestDelayCapAtMax(t *testing.T) {
	s := backoff.Strategy{
		Base:   1 * time.Second,
		Max:    3 * time.Second,
		Factor: 4.0,
		Jitter: false,
	}

	// attempt 2 would be 1s * 4^2 = 16s, but must be capped at 3s.
	got := s.Delay(2)
	if got != 3*time.Second {
		t.Errorf("expected delay capped at 3s, got %v", got)
	}
}

func TestDelayWithJitterIsWithinBounds(t *testing.T) {
	s := backoff.Strategy{
		Base:   100 * time.Millisecond,
		Max:    10 * time.Second,
		Factor: 2.0,
		Jitter: true,
	}

	for attempt := 0; attempt < 5; attempt++ {
		got := s.Delay(attempt)
		if got < s.Base || got > s.Max {
			t.Errorf("attempt %d: jitter delay %v out of bounds [%v, %v]", attempt, got, s.Base, s.Max)
		}
	}
}

func TestNegativeAttemptTreatedAsZero(t *testing.T) {
	s := backoff.Strategy{
		Base:   200 * time.Millisecond,
		Max:    5 * time.Second,
		Factor: 2.0,
		Jitter: false,
	}

	if got := s.Delay(-3); got != 200*time.Millisecond {
		t.Errorf("expected Base delay for negative attempt, got %v", got)
	}
}

func TestResetReturnsBase(t *testing.T) {
	s := backoff.Default()
	if got := s.Reset(); got != s.Base {
		t.Errorf("Reset() expected %v, got %v", s.Base, got)
	}
}
