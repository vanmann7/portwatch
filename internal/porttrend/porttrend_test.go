package porttrend

import (
	"testing"
	"time"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestStableWithNoEvents(t *testing.T) {
	tr := New(time.Minute)
	if got := tr.Trend(80); got != Stable {
		t.Fatalf("expected Stable, got %s", got)
	}
}

func TestRisingAfterMoreOpens(t *testing.T) {
	tr := New(time.Minute)
	tr.RecordOpen(443)
	tr.RecordOpen(443)
	tr.RecordClose(443)
	if got := tr.Trend(443); got != Rising {
		t.Fatalf("expected Rising, got %s", got)
	}
}

func TestFallingAfterMoreCloses(t *testing.T) {
	tr := New(time.Minute)
	tr.RecordOpen(22)
	tr.RecordClose(22)
	tr.RecordClose(22)
	if got := tr.Trend(22); got != Falling {
		t.Fatalf("expected Falling, got %s", got)
	}
}

func TestStableWhenEqualOpenClose(t *testing.T) {
	tr := New(time.Minute)
	tr.RecordOpen(8080)
	tr.RecordClose(8080)
	if got := tr.Trend(8080); got != Stable {
		t.Fatalf("expected Stable, got %s", got)
	}
}

func TestEvictsOldEvents(t *testing.T) {
	base := time.Now()
	tr := New(time.Minute)

	// inject old opens outside window
	tr.now = fixedClock(base.Add(-2 * time.Minute))
	tr.RecordOpen(9000)
	tr.RecordOpen(9000)

	// inject recent close inside window
	tr.now = fixedClock(base)
	tr.RecordClose(9000)

	if got := tr.Trend(9000); got != Falling {
		t.Fatalf("expected Falling after eviction, got %s", got)
	}
}

func TestResetClearsTrend(t *testing.T) {
	tr := New(time.Minute)
	tr.RecordOpen(3306)
	tr.RecordOpen(3306)
	tr.Reset(3306)
	if got := tr.Trend(3306); got != Stable {
		t.Fatalf("expected Stable after reset, got %s", got)
	}
}

func TestTrendStringRepresentation(t *testing.T) {
	cases := []struct {
		trend Trend
		want  string
	}{
		{Stable, "stable"},
		{Rising, "rising"},
		{Falling, "falling"},
	}
	for _, c := range cases {
		if got := c.trend.String(); got != c.want {
			t.Errorf("Trend(%d).String() = %q, want %q", c.trend, got, c.want)
		}
	}
}
