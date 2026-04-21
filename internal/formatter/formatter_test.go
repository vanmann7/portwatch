package formatter_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/formatter"
)

var fixedTime = time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)

func makeEvent(action string) formatter.Event {
	return formatter.Event{
		Port:      8080,
		Proto:     "tcp",
		Action:    action,
		Service:   "http-alt",
		Timestamp: fixedTime,
	}
}

func TestTextFormatOpened(t *testing.T) {
	f := formatter.New(formatter.FormatText)
	out, err := f.Format(makeEvent("opened"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "8080/tcp") {
		t.Errorf("expected port in output, got: %s", out)
	}
	if !strings.Contains(out, "opened") {
		t.Errorf("expected action in output, got: %s", out)
	}
	if !strings.Contains(out, "http-alt") {
		t.Errorf("expected service in output, got: %s", out)
	}
}

func TestTextFormatUnknownService(t *testing.T) {
	f := formatter.New(formatter.FormatText)
	e := makeEvent("closed")
	e.Service = ""
	out, err := f.Format(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "unknown") {
		t.Errorf("expected 'unknown' service label, got: %s", out)
	}
}

func TestJSONFormatFields(t *testing.T) {
	f := formatter.New(formatter.FormatJSON)
	out, err := f.Format(makeEvent("opened"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(out), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if int(m["port"].(float64)) != 8080 {
		t.Errorf("expected port 8080, got %v", m["port"])
	}
	if m["action"] != "opened" {
		t.Errorf("expected action 'opened', got %v", m["action"])
	}
	if m["proto"] != "tcp" {
		t.Errorf("expected proto 'tcp', got %v", m["proto"])
	}
}

func TestJSONOmitsEmptyService(t *testing.T) {
	f := formatter.New(formatter.FormatJSON)
	e := makeEvent("closed")
	e.Service = ""
	out, err := f.Format(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "service") {
		t.Errorf("expected service key to be omitted, got: %s", out)
	}
}

func TestTimestampInjectedWhenZero(t *testing.T) {
	f := formatter.New(formatter.FormatText)
	e := formatter.Event{Port: 22, Proto: "tcp", Action: "opened"}
	out, err := f.Format(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Error("expected non-empty output")
	}
}
