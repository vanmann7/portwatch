package portpulse

import "time"

// Score represents the computed pulse score for a port.
type Score struct {
	Port  int
	Value float64 // hits per hour since first seen
}

// Scorer computes pulse scores from Tracker data.
type Scorer struct {
	tracker *Tracker
	clock   func() time.Time
}

// NewScorer wraps a Tracker with scoring logic.
func NewScorer(tr *Tracker) *Scorer {
	return &Scorer{tracker: tr, clock: time.Now}
}

// ScorePort returns the pulse score for a single port.
// Score is hits-per-hour since first observation; returns 0 if unknown.
func (s *Scorer) ScorePort(port int) Score {
	e, ok := s.tracker.Get(port)
	if !ok {
		return Score{Port: port, Value: 0}
	}
	return Score{Port: port, Value: s.compute(e)}
}

// All returns scores for every tracked port.
func (s *Scorer) All() []Score {
	snap := s.tracker.Snapshot()
	out := make([]Score, len(snap))
	for i, e := range snap {
		out[i] = Score{Port: e.Port, Value: s.compute(e)}
	}
	return out
}

func (s *Scorer) compute(e Entry) float64 {
	now := s.clock()
	dur := now.Sub(e.FirstSeen)
	if dur <= 0 {
		return float64(e.Hits)
	}
	hours := dur.Hours()
	if hours < 1.0/3600 {
		// less than one second — treat as instantaneous
		return float64(e.Hits)
	}
	return float64(e.Hits) / hours
}
