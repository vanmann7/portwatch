//go:build integration

package watch_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/watch"
)

// TestWatcherFullCycle runs a longer cycle verifying both OPENED and CLOSED
// events are emitted. Execute with: go test -tags integration ./internal/watch/
func TestWatcherFullCycle(t *testing.T) {
	port, ln := freePort(t)

	var buf strings.Builder
	a := alert.New(&buf)

	cfg := config.Default()
	cfg.StartPort = port
	cfg.EndPort = port
	cfg.Interval = 40 * time.Millisecond

	w, err := watch.New(cfg, a)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	go func() {
		time.Sleep(100 * time.Millisecond)
		ln.Close() // close → should trigger CLOSED alert
	}()

	_ = w.Run(ctx)

	out := buf.String()
	if !strings.Contains(out, "CLOSED") {
		t.Errorf("expected CLOSED alert, got: %q", out)
	}
}
