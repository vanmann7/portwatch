package portquota_test

import (
	"testing"
	"time"

	"portwatch/internal/portquota"
)

// TestQuotaRealTimeCycle verifies that the quota resets naturally after the
// real window elapses (uses a very short window to keep the test fast).
func TestQuotaRealTimeCycle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping real-time integration test in short mode")
	}

	tr := portquota.New(2, 100*time.Millisecond)

	if !tr.Allow(80) {
		t.Fatal("first alert should be allowed")
	}
	if !tr.Allow(80) {
		t.Fatal("second alert should be allowed")
	}
	if tr.Allow(80) {
		t.Fatal("third alert should be blocked")
	}

	time.Sleep(150 * time.Millisecond)

	if !tr.Allow(80) {
		t.Fatal("alert should be allowed after window reset")
	}
}
