package portwindow_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portwindow"
)

// TestWindowConcurrentAccess verifies that concurrent Record and Contains
// calls do not race or panic.
func TestWindowConcurrentAccess(t *testing.T) {
	w := portwindow.New(5 * time.Second)
	done := make(chan struct{})

	go func() {
		for i := 0; i < 200; i++ {
			w.Record(i % 100)
		}
		close(done)
	}()

	for i := 0; i < 200; i++ {
		w.Contains(i % 100)
	}
	<-done
}

// TestWindowExpiresRealTime exercises eviction using real wall-clock time.
func TestWindowExpiresRealTime(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping real-time test in short mode")
	}
	w := portwindow.New(100 * time.Millisecond)
	w.Record(9000)
	if !w.Contains(9000) {
		t.Fatal("expected port 9000 immediately after record")
	}
	time.Sleep(150 * time.Millisecond)
	if w.Contains(9000) {
		t.Fatal("expected port 9000 to be evicted after window expired")
	}
}
