package debounce_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/debounce"
)

func TestDebouncerCallsAfterWait(t *testing.T) {
	var count int32
	d := debounce.New(50*time.Millisecond, func() {
		atomic.AddInt32(&count, 1)
	})
	d.Trigger()
	time.Sleep(100 * time.Millisecond)
	if got := atomic.LoadInt32(&count); got != 1 {
		t.Fatalf("expected 1 call, got %d", got)
	}
}

func TestDebouncerResetsOnRapidTriggers(t *testing.T) {
	var count int32
	d := debounce.New(60*time.Millisecond, func() {
		atomic.AddInt32(&count, 1)
	})
	for i := 0; i < 5; i++ {
		d.Trigger()
		time.Sleep(20 * time.Millisecond)
	}
	time.Sleep(120 * time.Millisecond)
	if got := atomic.LoadInt32(&count); got != 1 {
		t.Fatalf("expected 1 call after debounce, got %d", got)
	}
}

func TestDebouncerFlush(t *testing.T) {
	var count int32
	d := debounce.New(500*time.Millisecond, func() {
		atomic.AddInt32(&count, 1)
	})
	d.Trigger()
	flushed := d.Flush()
	if !flushed {
		t.Fatal("expected Flush to return true")
	}
	time.Sleep(30 * time.Millisecond)
	if got := atomic.LoadInt32(&count); got != 1 {
		t.Fatalf("expected 1 call after flush, got %d", got)
	}
}

func TestDebouncerStop(t *testing.T) {
	var count int32
	d := debounce.New(50*time.Millisecond, func() {
		atomic.AddInt32(&count, 1)
	})
	d.Trigger()
	d.Stop()
	time.Sleep(100 * time.Millisecond)
	if got := atomic.LoadInt32(&count); got != 0 {
		t.Fatalf("expected 0 calls after stop, got %d", got)
	}
}

func TestFlushWithNoPending(t *testing.T) {
	d := debounce.New(50*time.Millisecond, func() {})
	if d.Flush() {
		t.Fatal("expected Flush to return false when no timer pending")
	}
}
