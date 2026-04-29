package portreport_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portreport"
)

// --- stub implementations ---

type stubLabeler struct{ m map[int]string }

func (s *stubLabeler) Label(p int) string {
	if v, ok := s.m[p]; ok {
		return v
	}
	return "unknown"
}

type stubSeverity struct{ m map[int]string }

func (s *stubSeverity) Level(p int) string {
	if v, ok := s.m[p]; ok {
		return v
	}
	return "info"
}

type stubPresence struct {
	firstSeen map[int]time.Time
	uptime    map[int]time.Duration
}

func (s *stubPresence) FirstSeen(p int) (time.Time, bool) {
	v, ok := s.firstSeen[p]
	return v, ok
}
func (s *stubPresence) Uptime(p int) (time.Duration, bool) {
	v, ok := s.uptime[p]
	return v, ok
}

func newBuilder() *portreport.Builder {
	return portreport.New(
		&stubLabeler{m: map[int]string{22: "ssh", 80: "http"}},
		&stubSeverity{m: map[int]string{22: "warning", 80: "info"}},
		&stubPresence{
			firstSeen: map[int]time.Time{22: time.Unix(1000, 0)},
			uptime:    map[int]time.Duration{22: 5 * time.Minute},
		},
	)
}

func TestBuildEntriesAreSorted(t *testing.T) {
	b := newBuilder()
	r := b.Build([]int{80, 22})
	if len(r.Entries) != 2 {
		t.Fatalf("want 2 entries, got %d", len(r.Entries))
	}
	if r.Entries[0].Port != 22 || r.Entries[1].Port != 80 {
		t.Errorf("entries not sorted: %v", r.Entries)
	}
}

func TestBuildLabelAndSeverity(t *testing.T) {
	b := newBuilder()
	r := b.Build([]int{22})
	e := r.Entries[0]
	if e.Label != "ssh" {
		t.Errorf("want label ssh, got %q", e.Label)
	}
	if e.Severity != "warning" {
		t.Errorf("want severity warning, got %q", e.Severity)
	}
}

func TestBuildPresencePopulated(t *testing.T) {
	b := newBuilder()
	r := b.Build([]int{22})
	e := r.Entries[0]
	if e.FirstSeen.IsZero() {
		t.Error("expected FirstSeen to be set")
	}
	if e.Uptime != 5*time.Minute {
		t.Errorf("want 5m uptime, got %v", e.Uptime)
	}
}

func TestBuildUnknownPresenceIsZero(t *testing.T) {
	b := newBuilder()
	r := b.Build([]int{80})
	e := r.Entries[0]
	if !e.FirstSeen.IsZero() {
		t.Error("expected zero FirstSeen for unknown port")
	}
	if e.Uptime != 0 {
		t.Errorf("expected zero uptime, got %v", e.Uptime)
	}
}

func TestBuildEmptyPorts(t *testing.T) {
	b := newBuilder()
	r := b.Build(nil)
	if len(r.Entries) != 0 {
		t.Errorf("want 0 entries, got %d", len(r.Entries))
	}
	if r.GeneratedAt.IsZero() {
		t.Error("GeneratedAt should not be zero")
	}
}

func TestSummaryNoPorts(t *testing.T) {
	r := portreport.Report{GeneratedAt: time.Now()}
	if s := r.Summary(); s == "" {
		t.Error("expected non-empty summary")
	}
}

func TestSummaryWithPorts(t *testing.T) {
	b := newBuilder()
	r := b.Build([]int{22, 80})
	s := r.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
