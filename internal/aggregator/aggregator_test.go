package aggregator_test

import (
	"sort"
	"testing"
	"time"

	"github.com/user/portwatch/internal/aggregator"
)

func TestAggregatorBatchesEvents(t *testing.T) {
	a := aggregator.New(50 * time.Millisecond)
	defer a.Stop()

	a.Add(aggregator.Event{Port: 8080, Action: "opened"})
	a.Add(aggregator.Event{Port: 9090, Action: "opened"})
	a.Add(aggregator.Event{Port: 3000, Action: "closed"})

	select {
	case s := <-a.Out():
		sort.Ints(s.Opened)
		if len(s.Opened) != 2 || s.Opened[0] != 8080 || s.Opened[1] != 9090 {
			t.Fatalf("unexpected opened: %v", s.Opened)
		}
		if len(s.Closed) != 1 || s.Closed[0] != 3000 {
			t.Fatalf("unexpected closed: %v", s.Closed)
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for summary")
	}
}

func TestAggregatorFlushEmitsImmediately(t *testing.T) {
	a := aggregator.New(10 * time.Second)
	defer a.Stop()

	a.Add(aggregator.Event{Port: 443, Action: "opened"})
	a.Flush()

	select {
	case s := <-a.Out():
		if len(s.Opened) != 1 || s.Opened[0] != 443 {
			t.Fatalf("unexpected summary: %+v", s)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out after flush")
	}
}

func TestAggregatorOpenCancelsClosed(t *testing.T) {
	a := aggregator.New(10 * time.Second)
	defer a.Stop()

	a.Add(aggregator.Event{Port: 80, Action: "closed"})
	a.Add(aggregator.Event{Port: 80, Action: "opened"}) // cancels the closed
	a.Flush()

	s := <-a.Out()
	if len(s.Closed) != 0 {
		t.Fatalf("expected no closed ports, got %v", s.Closed)
	}
	if len(s.Opened) != 1 || s.Opened[0] != 80 {
		t.Fatalf("expected port 80 opened, got %v", s.Opened)
	}
}

func TestAggregatorStopClosesChannel(t *testing.T) {
	a := aggregator.New(50 * time.Millisecond)
	a.Stop()
	_, ok := <-a.Out()
	if ok {
		t.Fatal("expected channel to be closed")
	}
}

func TestAggregatorFlushNoOp(t *testing.T) {
	a := aggregator.New(50 * time.Millisecond)
	defer a.Stop()
	a.Flush() // no events — should not block or send
	select {
	case s := <-a.Out():
		t.Fatalf("unexpected summary on empty flush: %+v", s)
	case <-time.After(80 * time.Millisecond):
		// correct: nothing emitted
	}
}
