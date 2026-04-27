package portquota

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllowUpToMax(t *testing.T) {
	tr := newWithClock(3, time.Minute, fixedClock(time.Now()))
	for i := 0; i < 3; i++ {
		if !tr.Allow(80) {
			t.Fatalf("expected Allow=true on call %d", i+1)
		}
	}
	if tr.Allow(80) {
		t.Fatal("expected Allow=false after quota exhausted")
	}
}

func TestAllowResetsAfterWindow(t *testing.T) {
	now := time.Now()
	clock := func() time.Time { return now }
	tr := newWithClock(2, time.Second, clock)

	tr.Allow(443)
	tr.Allow(443)
	if tr.Allow(443) {
		t.Fatal("expected quota to be exhausted")
	}

	// Advance past the window.
	now = now.Add(2 * time.Second)
	if !tr.Allow(443) {
		t.Fatal("expected Allow=true after window reset")
	}
}

func TestRemainingDecrementsOnAllow(t *testing.T) {
	tr := newWithClock(5, time.Minute, fixedClock(time.Now()))
	if r := tr.Remaining(22); r != 5 {
		t.Fatalf("expected 5, got %d", r)
	}
	tr.Allow(22)
	tr.Allow(22)
	if r := tr.Remaining(22); r != 3 {
		t.Fatalf("expected 3, got %d", r)
	}
}

func TestRemainingZeroWhenExhausted(t *testing.T) {
	tr := newWithClock(2, time.Minute, fixedClock(time.Now()))
	tr.Allow(8080)
	tr.Allow(8080)
	tr.Allow(8080) // over quota
	if r := tr.Remaining(8080); r != 0 {
		t.Fatalf("expected 0, got %d", r)
	}
}

func TestResetAllowsImmediately(t *testing.T) {
	tr := newWithClock(1, time.Hour, fixedClock(time.Now()))
	tr.Allow(9000)
	if tr.Allow(9000) {
		t.Fatal("expected quota exhausted before reset")
	}
	tr.Reset(9000)
	if !tr.Allow(9000) {
		t.Fatal("expected Allow=true after Reset")
	}
}

func TestResetAllClearsAllPorts(t *testing.T) {
	tr := newWithClock(1, time.Hour, fixedClock(time.Now()))
	tr.Allow(1)
	tr.Allow(2)
	tr.ResetAll()
	if !tr.Allow(1) || !tr.Allow(2) {
		t.Fatal("expected all ports to be reset")
	}
}

func TestIndependentPortsDoNotInterfere(t *testing.T) {
	tr := newWithClock(1, time.Minute, fixedClock(time.Now()))
	tr.Allow(80)
	if !tr.Allow(443) {
		t.Fatal("expected port 443 to be unaffected by port 80 quota")
	}
}
