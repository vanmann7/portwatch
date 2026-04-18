package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/state"
)

func TestNotifyOpened(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	diff := state.Diff{Opened: []int{8080, 9090}}
	alerts := n.Notify(diff)

	if len(alerts) != 2 {
		t.Fatalf("expected 2 alerts, got %d", len(alerts))
	}
	for _, a := range alerts {
		if a.Level != alert.LevelWarn {
			t.Errorf("expected WARN level for opened port, got %s", a.Level)
		}
	}
	output := buf.String()
	if !strings.Contains(output, "8080") || !strings.Contains(output, "9090") {
		t.Errorf("output missing expected ports: %s", output)
	}
}

func TestNotifyClosed(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	diff := state.Diff{Closed: []int{22}}
	alerts := n.Notify(diff)

	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Level != alert.LevelInfo {
		t.Errorf("expected INFO level for closed port, got %s", alerts[0].Level)
	}
	if !strings.Contains(buf.String(), "22") {
		t.Errorf("output missing port 22: %s", buf.String())
	}
}

func TestNotifyNoChange(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	alerts := n.Notify(state.Diff{})

	if len(alerts) != 0 {
		t.Errorf("expected no alerts, got %d", len(alerts))
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got: %s", buf.String())
	}
}

func TestNotifyDefaultWriter(t *testing.T) {
	// Should not panic when writer is nil (falls back to stdout).
	n := alert.New(nil)
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
}
