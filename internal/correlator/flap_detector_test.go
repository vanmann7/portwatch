package correlator

import (
	"testing"
	"time"
)

func TestObserveBelowThresholdNotFlapping(t *testing.T) {
	fd := NewFlapDetector(time.Minute, 3)
	fd.Observe(80)
	fd.Observe(80)

	if fd.IsFlapping(80) {
		t.Error("expected port 80 not to be flapping with only 2 observations")
	}
}

func TestObserveExceedsThresholdIsFlapping(t *testing.T) {
	fd := NewFlapDetector(time.Minute, 3)
	for i := 0; i < 4; i++ {
		fd.Observe(443)
	}

	if !fd.IsFlapping(443) {
		t.Error("expected port 443 to be flapping after 4 observations")
	}
}

func TestObserveReturnsTrueWhenFlapping(t *testing.T) {
	fd := NewFlapDetector(time.Minute, 2)
	fd.Observe(8080)
	fd.Observe(8080)
	flapping := fd.Observe(8080)

	if !flapping {
		t.Error("expected Observe to return true when threshold exceeded")
	}
}

func TestOldObservationsEvicted(t *testing.T) {
	fd := NewFlapDetector(10*time.Millisecond, 2)
	fd.Observe(9000)
	fd.Observe(9000)
	fd.Observe(9000)

	time.Sleep(20 * time.Millisecond)

	// After the window expires, a fresh observation should not trigger flapping.
	flapping := fd.Observe(9000)
	if flapping {
		t.Error("expected old observations to be evicted and port not flapping")
	}
}

func TestResetClearsFlappingState(t *testing.T) {
	fd := NewFlapDetector(time.Minute, 1)
	fd.Observe(22)
	fd.Observe(22)
	fd.Reset()

	if fd.IsFlapping(22) {
		t.Error("expected flapping state to be cleared after Reset")
	}
}

func TestUnknownPortNotFlapping(t *testing.T) {
	fd := NewFlapDetector(time.Minute, 3)
	if fd.IsFlapping(12345) {
		t.Error("expected unknown port to not be flapping")
	}
}
