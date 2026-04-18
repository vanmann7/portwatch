package watch_test

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/watch"
)

func freePort(t *testing.T) (int, net.Listener) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("freePort: %v", err)
	}
	return ln.Addr().(*net.TCPAddr).Port, ln
}

func TestWatcherDetectsNewPort(t *testing.T) {
	port, ln := freePort(t)
	ln.Close() // start with port closed

	var buf strings.Builder
	a := alert.New(&buf)

	cfg := config.Default()
	cfg.StartPort = port
	cfg.EndPort = port
	cfg.Interval = 50 * time.Millisecond

	w, err := watch.New(cfg, a)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	// Open the port after a short delay so the second scan sees it.
	go func() {
		time.Sleep(80 * time.Millisecond)
		ln2, _ := net.Listen("tcp", ln.Addr().String())
		if ln2 != nil {
			defer ln2.Close()
			time.Sleep(200 * time.Millisecond)
		}
	}()

	_ = w.Run(ctx) // returns context.DeadlineExceeded — that's fine

	if !strings.Contains(buf.String(), "OPENED") {
		t.Errorf("expected OPENED alert, got: %q", buf.String())
	}
}

func TestWatcherCancelImmediately(t *testing.T) {
	cfg := config.Default()
	cfg.Interval = 1 * time.Second

	var buf strings.Builder
	a := alert.New(&buf)

	w, err := watch.New(cfg, a)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := w.Run(ctx); err != context.Canceled {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}
