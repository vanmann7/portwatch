package cooldown_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/cooldown"
)

func TestAllowFirstTime(t *testing.T) {
	tr := cooldown.New(5 * time.Second)
	if !tr.Allow(8080) {
		t.Fatal("expected Allow to return true on first call")
	}
}

func TestBlockWithinCooldown(t *testing.T) {
	tr := cooldown.New(5 * time.Second)
	if !tr.Allow(8080) {
		t.Fatal("first Allow should succeed")
	}
	if tr.Allow(8080) {
		t.Fatal("second Allow within cooldown should be blocked")
	}
}

func TestAllowAfterCooldownExpires(t *testing.T) {
	now := time.Now()
	tr := cooldown.New(5 * time.Second)

	// Inject a fake clock via a sub-test helper by wrapping New behaviour.
	// Use a real short interval instead.
	tr2 := cooldown.New(10 * time.Millisecond)
	_ = now
	if !tr2.Allow(9090) {
		t.Fatal("first call should succeed")
	}
	time.Sleep(20 * time.Millisecond)
	if !tr2.Allow(9090) {
		t.Fatal("should allow after interval expires")
	}
}

func TestResetAllowsImmediately(t *testing.T) {
	tr := cooldown.New(1 * time.Hour)
	tr.Allow(443)
	tr.Reset(443)
	if !tr.Allow(443) {
		t.Fatal("expected Allow after Reset to return true")
	}
}

func TestResetAll(t *testing.T) {
	tr := cooldown.New(1 * time.Hour)
	tr.Allow(80)
	tr.Allow(443)
	tr.ResetAll()
	if tr.Count() != 0 {
		t.Fatalf("expected 0 entries after ResetAll, got %d", tr.Count())
	}
}

func TestCountTracksEntries(t *testing.T) {
	tr := cooldown.New(1 * time.Hour)
	tr.Allow(80)
	tr.Allow(443)
	tr.Allow(8080)
	if tr.Count() != 3 {
		t.Fatalf("expected count 3, got %d", tr.Count())
	}
}

func TestIndependentPorts(t *testing.T) {
	tr := cooldown.New(1 * time.Hour)
	tr.Allow(80)
	if !tr.Allow(443) {
		t.Fatal("different port should not be affected by cooldown on port 80")
	}
}
