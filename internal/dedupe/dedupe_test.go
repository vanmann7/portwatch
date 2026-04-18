package dedupe

import (
	"testing"
	"time"
)

func TestNotDuplicateFirstTime(t *testing.T) {
	d := New(5 * time.Second)
	if d.IsDuplicate(8080, "opened") {
		t.Fatal("first occurrence should not be a duplicate")
	}
}

func TestDuplicateWithinWindow(t *testing.T) {
	d := New(5 * time.Second)
	d.IsDuplicate(8080, "opened")
	if !d.IsDuplicate(8080, "opened") {
		t.Fatal("second occurrence within window should be duplicate")
	}
}

func TestDifferentEventNotDuplicate(t *testing.T) {
	d := New(5 * time.Second)
	d.IsDuplicate(8080, "opened")
	if d.IsDuplicate(8080, "closed") {
		t.Fatal("different event type should not be duplicate")
	}
}

func TestEvictedAfterWindow(t *testing.T) {
	now := time.Now()
	d := New(1 * time.Second)
	d.now = func() time.Time { return now }

	d.IsDuplicate(9090, "opened")

	// advance time past window
	d.now = func() time.Time { return now.Add(2 * time.Second) }

	if d.IsDuplicate(9090, "opened") {
		t.Fatal("entry should have been evicted after window")
	}
}

func TestResetClearsEntries(t *testing.T) {
	d := New(10 * time.Second)
	d.IsDuplicate(1234, "opened")
	d.Reset()
	if d.IsDuplicate(1234, "opened") {
		t.Fatal("after reset, entry should not be duplicate")
	}
}

func TestDifferentPortNotDuplicate(t *testing.T) {
	d := New(10 * time.Second)
	d.IsDuplicate(80, "opened")
	if d.IsDuplicate(443, "opened") {
		t.Fatal("different port should not be duplicate")
	}
}
