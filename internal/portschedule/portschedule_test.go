package portschedule

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestDueWhenNeverScanned(t *testing.T) {
	s := New(time.Minute)
	if !s.Due(80) {
		t.Fatal("expected port 80 to be due on first check")
	}
}

func TestNotDueImmediatelyAfterRecord(t *testing.T) {
	now := time.Now()
	s := New(time.Minute)
	s.now = fixedNow(now)
	s.Record(80)
	if s.Due(80) {
		t.Fatal("expected port 80 to not be due immediately after record")
	}
}

func TestDueAfterIntervalElapsed(t *testing.T) {
	base := time.Now()
	s := New(time.Minute)
	s.now = fixedNow(base)
	s.Record(80)

	// Advance clock past the interval.
	s.now = fixedNow(base.Add(time.Minute + time.Second))
	if !s.Due(80) {
		t.Fatal("expected port 80 to be due after interval elapsed")
	}
}

func TestResetMakesPortDue(t *testing.T) {
	now := time.Now()
	s := New(time.Minute)
	s.now = fixedNow(now)
	s.Record(443)
	s.Reset(443)
	if !s.Due(443) {
		t.Fatal("expected port 443 to be due after reset")
	}
}

func TestSnapshotContainsRecordedPorts(t *testing.T) {
	s := New(time.Minute)
	s.Record(22)
	s.Record(80)
	s.Record(443)

	snap := s.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected 3 snapshot entries, got %d", len(snap))
	}

	ports := make(map[int]bool)
	for _, e := range snap {
		ports[e.Port] = true
	}
	for _, p := range []int{22, 80, 443} {
		if !ports[p] {
			t.Errorf("expected port %d in snapshot", p)
		}
	}
}

func TestSnapshotIsIndependent(t *testing.T) {
	s := New(time.Minute)
	s.Record(8080)
	snap1 := s.Snapshot()
	s.Record(9090)
	snap2 := s.Snapshot()

	if len(snap1) != 1 {
		t.Fatalf("snap1: expected 1 entry, got %d", len(snap1))
	}
	if len(snap2) != 2 {
		t.Fatalf("snap2: expected 2 entries, got %d", len(snap2))
	}
}
