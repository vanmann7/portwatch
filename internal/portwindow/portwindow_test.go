package portwindow

import (
	"testing"
	"time"
)

func TestRecordAndContains(t *testing.T) {
	w := New(5 * time.Second)
	w.Record(80)
	if !w.Contains(80) {
		t.Fatal("expected port 80 to be contained")
	}
}

func TestContainsFalseForUnrecorded(t *testing.T) {
	w := New(5 * time.Second)
	if w.Contains(443) {
		t.Fatal("expected port 443 not to be contained")
	}
}

func TestEvictsAfterWindow(t *testing.T) {
	now := time.Now()
	w := New(1 * time.Second)
	w.now = func() time.Time { return now }
	w.Record(22)

	// advance clock beyond window
	w.now = func() time.Time { return now.Add(2 * time.Second) }
	if w.Contains(22) {
		t.Fatal("expected port 22 to be evicted")
	}
}

func TestRecordUpdatesTimestamp(t *testing.T) {
	now := time.Now()
	w := New(2 * time.Second)
	w.now = func() time.Time { return now }
	w.Record(8080)

	// advance to just before expiry, then re-record to refresh
	w.now = func() time.Time { return now.Add(1500 * time.Millisecond) }
	w.Record(8080)

	// advance past original expiry but within refreshed window
	w.now = func() time.Time { return now.Add(2500 * time.Millisecond) }
	if !w.Contains(8080) {
		t.Fatal("expected refreshed port 8080 to still be present")
	}
}

func TestPortsSnapshot(t *testing.T) {
	w := New(10 * time.Second)
	w.Record(80)
	w.Record(443)
	w.Record(22)
	ports := w.Ports()
	if len(ports) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(ports))
	}
}

func TestReset(t *testing.T) {
	w := New(10 * time.Second)
	w.Record(80)
	w.Record(443)
	w.Reset()
	if len(w.Ports()) != 0 {
		t.Fatal("expected empty window after reset")
	}
}

func TestMultiplePortsIndependent(t *testing.T) {
	now := time.Now()
	w := New(2 * time.Second)
	w.now = func() time.Time { return now }
	w.Record(80)

	w.now = func() time.Time { return now.Add(1 * time.Second) }
	w.Record(443)

	// advance past first port's window
	w.now = func() time.Time { return now.Add(2500 * time.Millisecond) }
	if w.Contains(80) {
		t.Fatal("expected port 80 to be evicted")
	}
	if !w.Contains(443) {
		t.Fatal("expected port 443 to still be present")
	}
}
