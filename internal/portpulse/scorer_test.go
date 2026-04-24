package portpulse

import (
	"testing"
	"time"
)

func TestScorePortUnknownReturnsZero(t *testing.T) {
	tr := New()
	sc := NewScorer(tr)
	s := sc.ScorePort(9999)
	if s.Value != 0 {
		t.Fatalf("expected 0 for unknown port, got %f", s.Value)
	}
}

func TestScorePortHitsPerHour(t *testing.T) {
	base := time.Unix(0, 0)
	tr := New()
	tr.clock = fixedClock(base)
	tr.Record(80)
	tr.Record(80)

	sc := NewScorer(tr)
	// advance clock by 2 hours
	sc.clock = fixedClock(base.Add(2 * time.Hour))

	s := sc.ScorePort(80)
	// 2 hits over 2 hours = 1.0 hits/hour
	if s.Value != 1.0 {
		t.Fatalf("expected 1.0 hits/hour, got %f", s.Value)
	}
}

func TestScorePortInstantaneous(t *testing.T) {
	base := time.Unix(1000, 0)
	tr := New()
	tr.clock = fixedClock(base)
	tr.Record(443)
	tr.Record(443)
	tr.Record(443)

	sc := NewScorer(tr)
	sc.clock = fixedClock(base) // same instant

	s := sc.ScorePort(443)
	// duration <= 0, returns raw hit count
	if s.Value != 3.0 {
		t.Fatalf("expected 3.0 for instantaneous, got %f", s.Value)
	}
}

func TestAllReturnsScoresForAllPorts(t *testing.T) {
	base := time.Unix(0, 0)
	tr := New()
	tr.clock = fixedClock(base)
	tr.Record(22)
	tr.Record(22)
	tr.Record(8080)

	sc := NewScorer(tr)
	sc.clock = fixedClock(base.Add(1 * time.Hour))

	scores := sc.All()
	if len(scores) != 2 {
		t.Fatalf("expected 2 scores, got %d", len(scores))
	}

	byPort := make(map[int]float64)
	for _, s := range scores {
		byPort[s.Port] = s.Value
	}
	if byPort[22] != 2.0 {
		t.Fatalf("expected 2.0 for port 22, got %f", byPort[22])
	}
	if byPort[8080] != 1.0 {
		t.Fatalf("expected 1.0 for port 8080, got %f", byPort[8080])
	}
}
