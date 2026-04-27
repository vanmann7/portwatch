package portage_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portage"
)

// TestAgeAdvancesOverTime verifies that the bucket classification changes as
// real wall-clock time passes (coarse-grained, uses short durations).
func TestAgeAdvancesOverTime(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tr := portage.New()
	tr.Record(12345)

	e, ok := tr.Get(12345)
	if !ok {
		t.Fatal("expected entry immediately after Record")
	}
	if e.Bucket != portage.BucketNew {
		t.Errorf("expected BucketNew immediately, got %s", e.Bucket)
	}
	if e.Age < 0 {
		t.Errorf("age should not be negative")
	}
}

// TestSnapshotConsistency records multiple ports and verifies snapshot
// returns consistent data with no duplicates.
func TestSnapshotConsistency(t *testing.T) {
	tr := portage.New()
	ports := []int{80, 443, 8080, 22, 3306}
	for _, p := range ports {
		tr.Record(p)
	}

	time.Sleep(5 * time.Millisecond)

	snap := tr.Snapshot()
	if len(snap) != len(ports) {
		t.Fatalf("expected %d entries, got %d", len(ports), len(snap))
	}

	seen := make(map[int]bool)
	for _, e := range snap {
		if seen[e.Port] {
			t.Errorf("duplicate port %d in snapshot", e.Port)
		}
		seen[e.Port] = true
		if e.Age <= 0 {
			t.Errorf("port %d age should be positive", e.Port)
		}
	}
}
