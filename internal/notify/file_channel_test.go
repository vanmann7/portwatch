package notify_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/notify"
)

func TestFileChannelWritesEntry(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "portwatch.log")

	ch := notify.NewFileChannel(path)
	if err := ch.Send("OPENED", "port 9090"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read log file: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "OPENED") {
		t.Errorf("expected OPENED in log, got: %s", content)
	}
	if !strings.Contains(content, "port 9090") {
		t.Errorf("expected body in log, got: %s", content)
	}
}

func TestFileChannelAppendsMultiple(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "portwatch.log")
	ch := notify.NewFileChannel(path)

	for i := 0; i < 3; i++ {
		if err := ch.Send("EVENT", "entry"); err != nil {
			t.Fatalf("send error: %v", err)
		}
	}

	data, _ := os.ReadFile(path)
	lines := strings.Count(string(data), "EVENT")
	if lines != 3 {
		t.Errorf("expected 3 log lines, got %d", lines)
	}
}

func TestFileChannelInvalidPath(t *testing.T) {
	ch := notify.NewFileChannel("/nonexistent/dir/portwatch.log")
	if err := ch.Send("X", "y"); err == nil {
		t.Error("expected error for invalid path")
	}
}
