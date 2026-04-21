package presencetracker

import (
	"testing"
	"time"
)

// fixedClock returns a function that returns t unchanged, useful for
// deterministic tests that need to control the current time.
func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestOpenedRecordsFirstSeen(t *testing.T) {
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	tr := New(fixedClock(base))
	tr.Opened(8080)
	snap := tr.Snapshot()
	if snap[8080] != base {
		t.Fatalf("expected %v, got %v", base, snap[8080])
	}
}

func TestOpenedIsIdempotent(t *testing.T) {
	calls := 0
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	clock := func() time.Time {
		calls++
		return now.Add(time.Duration(calls) * time.Second)
	}
	tr := New(clock)
	tr.Opened(443)
	tr.Opened(443) // second call must not overwrite first-seen
	up, ok := tr.Uptime(443)
	if !ok {
		t.Fatal("port should be tracked")
	}
	// uptime is measured from the first Opened call; calls==3 now (two Opened + one Uptime)
	if up < time.Second {
		t.Fatalf("expected uptime >= 1s, got %v", up)
	}
}

func TestClosedRemovesPort(t *testing.T) {
	tr := New(nil)
	tr.Opened(22)
	tr.Closed(22)
	_, ok := tr.Uptime(22)
	if ok {
		t.Fatal("port should not be tracked after Closed")
	}
}

func TestUptimeUnknownPort(t *testing.T) {
	tr := New(nil)
	up, ok := tr.Uptime(9999)
	if ok || up != 0 {
		t.Fatalf("expected (0, false), got (%v, %v)", up, ok)
	}
}

func TestSnapshotIsCopy(t *testing.T) {
	base := time.Now()
	tr := New(fixedClock(base))
	tr.Opened(80)
	snap := tr.Snapshot()
	delete(snap, 80)
	if _, ok := tr.Uptime(80); !ok {
		t.Fatal("original tracker should not be affected by mutation of snapshot")
	}
}

func TestResetClearsAll(t *testing.T) {
	tr := New(nil)
	tr.Opened(80)
	tr.Opened(443)
	tr.Reset()
	if len(tr.Snapshot()) != 0 {
		t.Fatal("expected empty snapshot after Reset")
	}
}

func TestNilClockUsesRealClock(t *testing.T) {
	tr := New(nil)
	tr.Opened(8443)
	time.Sleep(2 * time.Millisecond)
	up, ok := tr.Uptime(8443)
	if !ok {
		t.Fatal("port should be tracked")
	}
	if up < time.Millisecond {
		t.Fatalf("expected non-trivial uptime, got %v", up)
	}
}
