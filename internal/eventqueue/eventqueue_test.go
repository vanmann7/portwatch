package eventqueue_test

import (
	"testing"

	"github.com/user/portwatch/internal/eventqueue"
)

func makeEvent(port int, kind string) eventqueue.Event {
	return eventqueue.Event{Port: port, Kind: kind, Source: "test"}
}

func TestPushAndPop(t *testing.T) {
	q := eventqueue.New(8)
	q.Push(makeEvent(80, "opened"))
	q.Push(makeEvent(443, "opened"))

	e, ok := q.Pop()
	if !ok {
		t.Fatal("expected event, got none")
	}
	if e.Port != 80 {
		t.Fatalf("want port 80, got %d", e.Port)
	}
	if q.Len() != 1 {
		t.Fatalf("want len 1, got %d", q.Len())
	}
}

func TestPopEmptyReturnsFalse(t *testing.T) {
	q := eventqueue.New(4)
	_, ok := q.Pop()
	if ok {
		t.Fatal("expected false on empty queue")
	}
}

func TestEvictsOldestWhenFull(t *testing.T) {
	q := eventqueue.New(3)
	for i := 0; i < 4; i++ {
		q.Push(makeEvent(i, "opened"))
	}
	if q.Len() != 3 {
		t.Fatalf("want len 3, got %d", q.Len())
	}
	// Oldest (port 0) should have been evicted.
	e, _ := q.Pop()
	if e.Port != 1 {
		t.Fatalf("want port 1 (oldest surviving), got %d", e.Port)
	}
}

func TestDrainReturnsAllAndClearsQueue(t *testing.T) {
	q := eventqueue.New(8)
	q.Push(makeEvent(22, "opened"))
	q.Push(makeEvent(22, "closed"))

	events := q.Drain()
	if len(events) != 2 {
		t.Fatalf("want 2 events, got %d", len(events))
	}
	if q.Len() != 0 {
		t.Fatalf("queue should be empty after Drain, got %d", q.Len())
	}
}

func TestDefaultMaxSizeApplied(t *testing.T) {
	q := eventqueue.New(0) // should default to 64
	for i := 0; i < 70; i++ {
		q.Push(makeEvent(i, "opened"))
	}
	if q.Len() != 64 {
		t.Fatalf("want len 64, got %d", q.Len())
	}
}
