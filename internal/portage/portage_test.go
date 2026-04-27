package portage

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) clock {
	return func() time.Time { return t }
}

func TestRecordIsIdempotent(t *testing.T) {
	base := time.Now()
	tr := newWithClock(fixedClock(base))
	tr.Record(80)
	tr.Record(80) // second call must not overwrite

	e, ok := tr.Get(80)
	if !ok {
		t.Fatal("expected entry")
	}
	if !e.OpenedAt.Equal(base) {
		t.Errorf("openedAt changed on second Record")
	}
}

func TestGetUnknownPortReturnsFalse(t *testing.T) {
	tr := New()
	_, ok := tr.Get(9999)
	if ok {
		t.Error("expected false for untracked port")
	}
}

func TestRemoveDeletesEntry(t *testing.T) {
	tr := New()
	tr.Record(443)
	tr.Remove(443)
	_, ok := tr.Get(443)
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestBucketNew(t *testing.T) {
	base := time.Now()
	tr := newWithClock(fixedClock(base))
	tr.Record(8080)

	// advance 2 minutes
	tr.now = fixedClock(base.Add(2 * time.Minute))
	e, _ := tr.Get(8080)
	if e.Bucket != BucketNew {
		t.Errorf("expected new, got %s", e.Bucket)
	}
}

func TestBucketEstablished(t *testing.T) {
	base := time.Now()
	tr := newWithClock(fixedClock(base))
	tr.Record(22)

	tr.now = fixedClock(base.Add(30 * time.Minute))
	e, _ := tr.Get(22)
	if e.Bucket != BucketEstablished {
		t.Errorf("expected established, got %s", e.Bucket)
	}
}

func TestBucketLongRunning(t *testing.T) {
	base := time.Now()
	tr := newWithClock(fixedClock(base))
	tr.Record(3306)

	tr.now = fixedClock(base.Add(2 * time.Hour))
	e, _ := tr.Get(3306)
	if e.Bucket != BucketLongRunning {
		t.Errorf("expected long-running, got %s", e.Bucket)
	}
}

func TestSnapshotReturnsAllPorts(t *testing.T) {
	tr := New()
	tr.Record(80)
	tr.Record(443)
	tr.Record(22)

	snap := tr.Snapshot()
	if len(snap) != 3 {
		t.Errorf("expected 3 entries, got %d", len(snap))
	}
}

func TestSnapshotExcludesRemovedPorts(t *testing.T) {
	tr := New()
	tr.Record(80)
	tr.Record(443)
	tr.Remove(80)

	snap := tr.Snapshot()
	if len(snap) != 1 {
		t.Errorf("expected 1 entry, got %d", len(snap))
	}
	if snap[0].Port != 443 {
		t.Errorf("expected port 443, got %d", snap[0].Port)
	}
}
