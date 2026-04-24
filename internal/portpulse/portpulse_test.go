package portpulse

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestRecordFirstHit(t *testing.T) {
	tr := New()
	tr.clock = fixedClock(time.Unix(1000, 0))
	tr.Record(80)

	e, ok := tr.Get(80)
	if !ok {
		t.Fatal("expected entry for port 80")
	}
	if e.Hits != 1 {
		t.Fatalf("expected 1 hit, got %d", e.Hits)
	}
	if !e.FirstSeen.Equal(time.Unix(1000, 0)) {
		t.Fatalf("unexpected FirstSeen: %v", e.FirstSeen)
	}
}

func TestRecordAccumulatesHits(t *testing.T) {
	tr := New()
	tr.Record(443)
	tr.Record(443)
	tr.Record(443)

	e, ok := tr.Get(443)
	if !ok {
		t.Fatal("expected entry for port 443")
	}
	if e.Hits != 3 {
		t.Fatalf("expected 3 hits, got %d", e.Hits)
	}
}

func TestGetUnknownPortReturnsFalse(t *testing.T) {
	tr := New()
	_, ok := tr.Get(9999)
	if ok {
		t.Fatal("expected false for unrecorded port")
	}
}

func TestSnapshotReturnsAllEntries(t *testing.T) {
	tr := New()
	tr.Record(22)
	tr.Record(80)
	tr.Record(80)

	snap := tr.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap))
	}
}

func TestResetClearsEntries(t *testing.T) {
	tr := New()
	tr.Record(8080)
	tr.Reset()

	snap := tr.Snapshot()
	if len(snap) != 0 {
		t.Fatalf("expected empty snapshot after reset, got %d", len(snap))
	}
}

func TestLastSeenUpdatesOnSubsequentRecord(t *testing.T) {
	tr := New()
	t1 := time.Unix(1000, 0)
	t2 := time.Unix(2000, 0)

	tr.clock = fixedClock(t1)
	tr.Record(53)

	tr.clock = fixedClock(t2)
	tr.Record(53)

	e, _ := tr.Get(53)
	if !e.FirstSeen.Equal(t1) {
		t.Fatalf("FirstSeen should remain %v, got %v", t1, e.FirstSeen)
	}
	if !e.LastSeen.Equal(t2) {
		t.Fatalf("LastSeen should be %v, got %v", t2, e.LastSeen)
	}
}
