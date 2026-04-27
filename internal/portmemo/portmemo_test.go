package portmemo

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestRecordAndGet(t *testing.T) {
	now := time.Now()
	m := newWithClock(fixedClock(now))
	m.Record(8080, "opened")

	e, ok := m.Get(8080)
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Port != 8080 {
		t.Errorf("port: got %d, want 8080", e.Port)
	}
	if e.Event != "opened" {
		t.Errorf("event: got %q, want \"opened\"", e.Event)
	}
	if !e.RecordedAt.Equal(now) {
		t.Errorf("timestamp: got %v, want %v", e.RecordedAt, now)
	}
}

func TestGetUnknownPortReturnsFalse(t *testing.T) {
	m := New()
	_, ok := m.Get(9999)
	if ok {
		t.Error("expected false for unknown port")
	}
}

func TestRecordOverwritesPreviousEntry(t *testing.T) {
	m := New()
	m.Record(22, "opened")
	m.Record(22, "closed")

	e, ok := m.Get(22)
	if !ok {
		t.Fatal("expected entry")
	}
	if e.Event != "closed" {
		t.Errorf("event: got %q, want \"closed\"", e.Event)
	}
}

func TestForgetRemovesEntry(t *testing.T) {
	m := New()
	m.Record(443, "opened")
	m.Forget(443)

	_, ok := m.Get(443)
	if ok {
		t.Error("expected entry to be removed after Forget")
	}
}

func TestSnapshotReturnsCopy(t *testing.T) {
	m := New()
	m.Record(80, "opened")
	m.Record(443, "opened")

	snap := m.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("snapshot len: got %d, want 2", len(snap))
	}

	// Mutating the snapshot must not affect the memo.
	delete(snap, 80)
	_, ok := m.Get(80)
	if !ok {
		t.Error("memo entry should not be affected by snapshot mutation")
	}
}

func TestClearRemovesAllEntries(t *testing.T) {
	m := New()
	m.Record(80, "opened")
	m.Record(22, "opened")
	m.Clear()

	if len(m.Snapshot()) != 0 {
		t.Error("expected memo to be empty after Clear")
	}
}
