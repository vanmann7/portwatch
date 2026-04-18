package notify_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/notify"
)

func TestStdoutChannelSend(t *testing.T) {
	var buf bytes.Buffer
	ch := notify.NewStdout(&buf)
	if err := ch.Send("ALERT", "port 8080 opened"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ALERT") {
		t.Errorf("expected subject in output, got: %s", out)
	}
	if !strings.Contains(out, "port 8080 opened") {
		t.Errorf("expected body in output, got: %s", out)
	}
}

func TestStdoutDefaultWriter(t *testing.T) {
	ch := notify.NewStdout(nil)
	if ch.Writer == nil {
		t.Error("expected non-nil writer")
	}
}

type errChannel struct{}

func (e *errChannel) Send(_, _ string) error {
	return errors.New("send failed")
}

func TestDispatcherCollectsErrors(t *testing.T) {
	var buf bytes.Buffer
	ok := notify.NewStdout(&buf)
	bad := &errChannel{}
	d := notify.NewDispatcher(ok, bad)
	errs := d.Dispatch("TEST", "body")
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestDispatcherNoErrors(t *testing.T) {
	var buf bytes.Buffer
	ch := notify.NewStdout(&buf)
	d := notify.NewDispatcher(ch)
	errs := d.Dispatch("INFO", "all good")
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}
