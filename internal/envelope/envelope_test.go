package envelope_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/envelope"
)

func TestNewEnvelopeDefaults(t *testing.T) {
	e := envelope.New("abc", envelope.DestStdout, 42)

	if e.ID != "abc" {
		t.Fatalf("expected ID abc, got %s", e.ID)
	}
	if e.Attempt != 1 {
		t.Fatalf("expected Attempt 1, got %d", e.Attempt)
	}
	if e.Payload != 42 {
		t.Fatalf("expected payload 42, got %d", e.Payload)
	}
	if time.Since(e.CreatedAt) > time.Second {
		t.Fatal("CreatedAt is too old")
	}
	if e.Labels == nil {
		t.Fatal("Labels map should be initialised")
	}
}

func TestHasDest(t *testing.T) {
	e := envelope.New("x", envelope.DestAll, "msg")

	if !e.HasDest(envelope.DestStdout) {
		t.Error("expected DestStdout to be set")
	}
	if !e.HasDest(envelope.DestFile) {
		t.Error("expected DestFile to be set")
	}
	if !e.HasDest(envelope.DestWebhook) {
		t.Error("expected DestWebhook to be set")
	}

	only := envelope.New("y", envelope.DestFile, "msg")
	if only.HasDest(envelope.DestStdout) {
		t.Error("DestStdout should not be set")
	}
}

func TestWithLabel(t *testing.T) {
	e := envelope.New("l", envelope.DestStdout, "payload")
	e2 := e.WithLabel("severity", "critical")

	if e2.Labels["severity"] != "critical" {
		t.Fatalf("expected label severity=critical, got %q", e2.Labels["severity"])
	}
	// original must not be modified
	if _, ok := e.Labels["severity"]; ok {
		t.Fatal("WithLabel must not mutate the original envelope")
	}
}

func TestWithLabelPreservesExisting(t *testing.T) {
	e := envelope.New("m", envelope.DestStdout, 0)
	e = e.WithLabel("a", "1")
	e = e.WithLabel("b", "2")

	if e.Labels["a"] != "1" || e.Labels["b"] != "2" {
		t.Fatalf("labels not preserved: %v", e.Labels)
	}
}

func TestNextAttempt(t *testing.T) {
	e := envelope.New("r", envelope.DestWebhook, "data")
	e2 := e.NextAttempt()
	e3 := e2.NextAttempt()

	if e.Attempt != 1 {
		t.Fatalf("original Attempt should remain 1, got %d", e.Attempt)
	}
	if e2.Attempt != 2 {
		t.Fatalf("expected Attempt 2, got %d", e2.Attempt)
	}
	if e3.Attempt != 3 {
		t.Fatalf("expected Attempt 3, got %d", e3.Attempt)
	}
}
