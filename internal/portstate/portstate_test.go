package portstate

import (
	"testing"
	"time"
)

var fixedClock = func() time.Time {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
}

func TestUpdateOpenTransition(t *testing.T) {
	tr := newWithClock(fixedClock)
	trans, changed := tr.Update(80, Open)
	if !changed {
		t.Fatal("expected changed=true")
	}
	if trans.Port != 80 {
		t.Errorf("port: got %d want 80", trans.Port)
	}
	if trans.Prev != Unknown {
		t.Errorf("prev: got %v want Unknown", trans.Prev)
	}
	if trans.Next != Open {
		t.Errorf("next: got %v want Open", trans.Next)
	}
	if !trans.ChangedAt.Equal(fixedClock()) {
		t.Errorf("changedAt: got %v want %v", trans.ChangedAt, fixedClock())
	}
}

func TestUpdateNoChangeWhenSameState(t *testing.T) {
	tr := newWithClock(fixedClock)
	tr.Update(443, Open)
	_, changed := tr.Update(443, Open)
	if changed {
		t.Fatal("expected changed=false for same state")
	}
}

func TestUpdateClosedTransition(t *testing.T) {
	tr := newWithClock(fixedClock)
	tr.Update(22, Open)
	trans, changed := tr.Update(22, Closed)
	if !changed {
		t.Fatal("expected changed=true")
	}
	if trans.Prev != Open || trans.Next != Closed {
		t.Errorf("transition: got %v->%v want Open->Closed", trans.Prev, trans.Next)
	}
}

func TestGetUnknownPort(t *testing.T) {
	tr := New()
	if s := tr.Get(9999); s != Unknown {
		t.Errorf("expected Unknown, got %v", s)
	}
}

func TestGetAfterUpdate(t *testing.T) {
	tr := New()
	tr.Update(8080, Open)
	if s := tr.Get(8080); s != Open {
		t.Errorf("expected Open, got %v", s)
	}
}

func TestSnapshot(t *testing.T) {
	tr := New()
	tr.Update(80, Open)
	tr.Update(443, Open)
	tr.Update(22, Closed)
	snap := tr.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(snap))
	}
	if snap[80] != Open {
		t.Errorf("port 80: expected Open")
	}
	if snap[22] != Closed {
		t.Errorf("port 22: expected Closed")
	}
}

func TestReset(t *testing.T) {
	tr := New()
	tr.Update(80, Open)
	tr.Reset()
	if s := tr.Get(80); s != Unknown {
		t.Errorf("after reset expected Unknown, got %v", s)
	}
	if len(tr.Snapshot()) != 0 {
		t.Error("snapshot should be empty after reset")
	}
}
