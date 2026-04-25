package portlifetime

import (
	"testing"
	"time"
)

func fixedClock(base time.Time) func() time.Time {
	return func() time.Time { return base }
}

func TestOpenCreatesEntry(t *testing.T) {
	base := time.Now()
	tr := newWithClock(fixedClock(base))
	tr.Open(8080)
	e := tr.Snapshot()
	if _, ok := e[8080]; !ok {
		t.Fatal("expected port 8080 to be tracked")
	}
	if !e[8080].OpenedAt.Equal(base) {
		t.Errorf("OpenedAt: got %v, want %v", e[8080].OpenedAt, base)
	}
}

func TestOpenIsIdempotent(t *testing.T) {
	base := time.Now()
	later := base.Add(5 * time.Second)
	clock := base
	tr := newWithClock(func() time.Time { return clock })
	tr.Open(443)
	clock = later
	tr.Open(443) // second call should update LastSeen but not OpenedAt
	e := tr.Snapshot()
	if !e[443].OpenedAt.Equal(base) {
		t.Errorf("OpenedAt changed: got %v, want %v", e[443].OpenedAt, base)
	}
	if !e[443].LastSeen.Equal(later) {
		t.Errorf("LastSeen: got %v, want %v", e[443].LastSeen, later)
	}
}

func TestLifetimeReturnsElapsed(t *testing.T) {
	base := time.Now()
	clock := base
	tr := newWithClock(func() time.Time { return clock })
	tr.Open(22)
	clock = base.Add(10 * time.Second)
	d, ok := tr.Lifetime(22)
	if !ok {
		t.Fatal("expected port 22 to be tracked")
	}
	if d != 10*time.Second {
		t.Errorf("Lifetime: got %v, want 10s", d)
	}
}

func TestLifetimeUnknownPort(t *testing.T) {
	tr := New()
	_, ok := tr.Lifetime(9999)
	if ok {
		t.Error("expected untracked port to return ok=false")
	}
}

func TestCloseRemovesPort(t *testing.T) {
	tr := New()
	tr.Open(80)
	tr.Close(80)
	_, ok := tr.Lifetime(80)
	if ok {
		t.Error("expected port 80 to be removed after Close")
	}
}

func TestResetClearsAll(t *testing.T) {
	tr := New()
	tr.Open(80)
	tr.Open(443)
	tr.Reset()
	if len(tr.Snapshot()) != 0 {
		t.Error("expected empty snapshot after Reset")
	}
}

func TestSnapshotIsCopy(t *testing.T) {
	tr := New()
	tr.Open(8080)
	snap := tr.Snapshot()
	delete(snap, 8080)
	if _, ok := tr.Snapshot()[8080]; !ok {
		t.Error("Snapshot modification should not affect tracker")
	}
}
