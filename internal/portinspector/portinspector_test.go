package portinspector_test

import (
	"strings"
	"testing"
	"time"

	"portwatch/internal/portinspector"
)

// --- stubs ---

type stubState struct{ m map[int]string }

func (s *stubState) Get(port int) (string, bool) {
	v, ok := s.m[port]
	return v, ok
}

type stubLabel struct{ m map[int]string }

func (l *stubLabel) Label(port int) string {
	if v, ok := l.m[port]; ok {
		return v
	}
	return "unknown"
}

type stubRank struct{ m map[int]int }

func (r *stubRank) Rank(port int) int { return r.m[port] }

type stubLifetime struct {
	durations map[int]time.Duration
	opened    map[int]time.Time
}

func (lt *stubLifetime) Lifetime(port int) (time.Duration, bool) {
	d, ok := lt.durations[port]
	return d, ok
}
func (lt *stubLifetime) OpenedAt(port int) (time.Time, bool) {
	t, ok := lt.opened[port]
	return t, ok
}

func newInspector() *portinspector.Inspector {
	now := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return portinspector.New(
		&stubState{m: map[int]string{22: "open", 80: "closed"}},
		&stubLabel{m: map[int]string{22: "ssh", 80: "http"}},
		&stubRank{m: map[int]int{22: 10, 80: 8}},
		&stubLifetime{
			durations: map[int]time.Duration{22: 5 * time.Minute},
			opened:    map[int]time.Time{22: now},
		},
	)
}

func TestInspectOpenPort(t *testing.T) {
	ins := newInspector()
	rec, ok := ins.Inspect(22)
	if !ok {
		t.Fatal("expected record for port 22")
	}
	if rec.Port != 22 || rec.State != "open" {
		t.Errorf("unexpected record: %+v", rec)
	}
	if rec.Label != "ssh" {
		t.Errorf("expected label ssh, got %q", rec.Label)
	}
	if rec.Rank != 10 {
		t.Errorf("expected rank 10, got %d", rec.Rank)
	}
	if rec.Lifetime != 5*time.Minute {
		t.Errorf("unexpected lifetime %s", rec.Lifetime)
	}
}

func TestInspectClosedPort(t *testing.T) {
	ins := newInspector()
	rec, ok := ins.Inspect(80)
	if !ok {
		t.Fatal("expected record for port 80")
	}
	if rec.State != "closed" {
		t.Errorf("expected closed, got %s", rec.State)
	}
}

func TestInspectUnknownPortReturnsFalse(t *testing.T) {
	ins := newInspector()
	_, ok := ins.Inspect(9999)
	if ok {
		t.Error("expected false for unknown port")
	}
}

func TestRecordStringContainsPort(t *testing.T) {
	ins := newInspector()
	rec, _ := ins.Inspect(22)
	s := rec.String()
	if !strings.Contains(s, "22") {
		t.Errorf("expected port in string, got %q", s)
	}
	if !strings.Contains(s, "ssh") {
		t.Errorf("expected label in string, got %q", s)
	}
}

func TestRecordStringClosedOmitsLifetime(t *testing.T) {
	ins := newInspector()
	rec, _ := ins.Inspect(80)
	s := rec.String()
	if strings.Contains(s, "lifetime") {
		t.Errorf("closed record should not include lifetime, got %q", s)
	}
}
