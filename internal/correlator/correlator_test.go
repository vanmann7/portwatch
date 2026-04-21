package correlator

import (
	"testing"
	"time"
)

func TestRecordReturnsConsistentID(t *testing.T) {
	c := New(time.Minute)
	id1 := c.Record(8080, Opened)
	id2 := c.Record(8080, Closed)

	if id1 != id2 {
		t.Errorf("expected same correlation ID within window, got %q and %q", id1, id2)
	}
}

func TestEventsReturnsRecordedEvents(t *testing.T) {
	c := New(time.Minute)
	id := c.Record(443, Opened)
	c.Record(443, Closed)

	events := c.Events(id)
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
	if events[0].Port != 443 || events[0].Kind != Opened {
		t.Errorf("unexpected first event: %+v", events[0])
	}
	if events[1].Kind != Closed {
		t.Errorf("unexpected second event kind: %v", events[1].Kind)
	}
}

func TestDifferentPortsGetDifferentIDs(t *testing.T) {
	c := New(time.Minute)
	id1 := c.Record(80, Opened)
	id2 := c.Record(443, Opened)

	if id1 == id2 {
		t.Error("expected different correlation IDs for different ports")
	}
}

func TestEventsUnknownIDReturnsEmpty(t *testing.T) {
	c := New(time.Minute)
	events := c.Events("nonexistent")
	if len(events) != 0 {
		t.Errorf("expected empty slice, got %d events", len(events))
	}
}

func TestResetClearsAllGroups(t *testing.T) {
	c := New(time.Minute)
	id := c.Record(8080, Opened)
	c.Reset()

	events := c.Events(id)
	if len(events) != 0 {
		t.Errorf("expected empty after reset, got %d events", len(events))
	}
}

func TestFlushRemovesOldBuckets(t *testing.T) {
	// Use a tiny window so events are immediately in a past bucket.
	c := New(time.Nanosecond)
	c.Record(9090, Opened)

	time.Sleep(2 * time.Millisecond)

	flushed := c.Flush()
	if len(flushed) == 0 {
		t.Error("expected at least one flushed group")
	}

	// Verify group is removed from internal state.
	for id := range flushed {
		if evs := c.Events(id); len(evs) != 0 {
			t.Errorf("expected group %q to be removed after flush", id)
		}
	}
}
