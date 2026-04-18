package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/audit"
)

func TestOpenedWritesEntry(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	if err := l.Opened(8080, "tcp"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var e audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if e.Event != "opened" {
		t.Errorf("event = %q, want opened", e.Event)
	}
	if e.Port != 8080 {
		t.Errorf("port = %d, want 8080", e.Port)
	}
	if e.Protocol != "tcp" {
		t.Errorf("protocol = %q, want tcp", e.Protocol)
	}
}

func TestClosedWritesEntry(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	if err := l.Closed(443, "tcp"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var e audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if e.Event != "closed" {
		t.Errorf("event = %q, want closed", e.Event)
	}
}

func TestRecordWithNote(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	if err := l.Record("suppressed", 22, "tcp", "in baseline"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "in baseline") {
		t.Errorf("expected note in output, got: %s", buf.String())
	}
}

func TestNewNilWriterUsesStdout(t *testing.T) {
	// Just ensure no panic when w is nil.
	l := audit.New(nil)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestMultipleEntriesNewlines(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	_ = l.Opened(80, "tcp")
	_ = l.Closed(80, "tcp")
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}
