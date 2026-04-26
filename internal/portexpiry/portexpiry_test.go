package portexpiry

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestRecordSetsOpenedAt(t *testing.T) {
	tr := newWithClock(time.Hour, fixedClock(epoch))
	tr.Record(80)
	snap := tr.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(snap))
	}
	if snap[0].OpenedAt != epoch {
		t.Errorf("unexpected OpenedAt: %v", snap[0].OpenedAt)
	}
	if snap[0].ExpiresAt != epoch.Add(time.Hour) {
		t.Errorf("unexpected ExpiresAt: %v", snap[0].ExpiresAt)
	}
}

func TestRecordIsIdempotent(t *testing.T) {
	clock := epoch
	tr := newWithClock(time.Hour, func() time.Time { return clock })
	tr.Record(443)
	clock = epoch.Add(10 * time.Minute)
	tr.Record(443) // should not overwrite
	snap := tr.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(snap))
	}
	if snap[0].OpenedAt != epoch {
		t.Errorf("Record overwrote existing entry; OpenedAt = %v", snap[0].OpenedAt)
	}
}

func TestNotExpiredBeforeMaxAge(t *testing.T) {
	tr := newWithClock(time.Hour, fixedClock(epoch))
	tr.Record(22)
	// advance 30 minutes — still within max age
	tr2 := newWithClock(time.Hour, fixedClock(epoch.Add(30*time.Minute)))
	tr2.entries = tr.entries
	tr2.max = tr.max
	expired := tr2.Expired()
	if len(expired) != 0 {
		t.Errorf("expected no expired ports, got %d", len(expired))
	}
}

func TestExpiredAfterMaxAge(t *testing.T) {
	var now = epoch
	tr := newWithClock(time.Hour, func() time.Time { return now })
	tr.Record(8080)
	now = epoch.Add(2 * time.Hour)
	expired := tr.Expired()
	if len(expired) != 1 {
		t.Fatalf("expected 1 expired entry, got %d", len(expired))
	}
	if !expired[0].Expired {
		t.Error("expected Expired flag to be true")
	}
	if expired[0].Port != 8080 {
		t.Errorf("unexpected port %d", expired[0].Port)
	}
}

func TestRemoveStopsTracking(t *testing.T) {
	tr := newWithClock(time.Hour, fixedClock(epoch))
	tr.Record(3306)
	tr.Remove(3306)
	if len(tr.Snapshot()) != 0 {
		t.Error("expected empty snapshot after Remove")
	}
}

func TestSnapshotMarksExpiredFlag(t *testing.T) {
	var now = epoch
	tr := newWithClock(time.Minute, func() time.Time { return now })
	tr.Record(5432)
	now = epoch.Add(2 * time.Minute)
	snap := tr.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 snapshot entry, got %d", len(snap))
	}
	if !snap[0].Expired {
		t.Error("expected Expired flag set in Snapshot")
	}
}
