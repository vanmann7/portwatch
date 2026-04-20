package rollup_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/rollup"
)

func TestRollupBatchesEvents(t *testing.T) {
	r, ch := rollup.New(50 * time.Millisecond)

	r.Add(rollup.Event{Port: 8080, Action: "opened"})
	r.Add(rollup.Event{Port: 9090, Action: "closed"})

	select {
	case batch := <-ch:
		if len(batch.Events) != 2 {
			t.Fatalf("expected 2 events, got %d", len(batch.Events))
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for batch")
	}

	r.Stop()
}

func TestRollupFlushEmitsImmediately(t *testing.T) {
	r, ch := rollup.New(5 * time.Second) // long window — flush must bypass it

	r.Add(rollup.Event{Port: 443, Action: "opened"})
	r.Flush()

	select {
	case batch := <-ch:
		if len(batch.Events) != 1 {
			t.Fatalf("expected 1 event, got %d", len(batch.Events))
		}
		if batch.Events[0].Port != 443 {
			t.Errorf("expected port 443, got %d", batch.Events[0].Port)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("timed out waiting for flushed batch")
	}

	r.Stop()
}

func TestRollupFlushNoOp(t *testing.T) {
	r, ch := rollup.New(50 * time.Millisecond)
	r.Flush() // nothing pending — should not send

	select {
	case <-ch:
		t.Fatal("expected no batch on empty flush")
	case <-time.After(80 * time.Millisecond):
		// correct — nothing emitted
	}

	r.Stop()
}

func TestRollupStopClosesChannel(t *testing.T) {
	r, ch := rollup.New(50 * time.Millisecond)
	r.Stop()

	_, ok := <-ch
	if ok {
		t.Fatal("expected channel to be closed after Stop")
	}
}

func TestRollupSlidingWindowResets(t *testing.T) {
	r, ch := rollup.New(60 * time.Millisecond)

	r.Add(rollup.Event{Port: 80, Action: "opened"})
	time.Sleep(40 * time.Millisecond)
	r.Add(rollup.Event{Port: 81, Action: "opened"}) // resets timer

	// First batch should contain both events after the window
	select {
	case batch := <-ch:
		if len(batch.Events) != 2 {
			t.Fatalf("expected 2 events in single batch, got %d", len(batch.Events))
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("timed out waiting for rolled-up batch")
	}

	r.Stop()
}
