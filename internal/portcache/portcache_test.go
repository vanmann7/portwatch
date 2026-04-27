package portcache_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/portcache"
)

func fixedClock(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestSetAndGetHit(t *testing.T) {
	now := time.Now()
	c := portcache.New(5 * time.Second)
	c.Set(80, true)
	e, ok := c.Get(80)
	if !ok {
		t.Fatal("expected cache hit")
	}
	if !e.Open {
		t.Error("expected port to be open")
	}
	_ = now
}

func TestGetMissForUnknownPort(t *testing.T) {
	c := portcache.New(5 * time.Second)
	_, ok := c.Get(9999)
	if ok {
		t.Error("expected cache miss for unknown port")
	}
}

func TestGetMissAfterTTLExpires(t *testing.T) {
	base := time.Now()
	clock := func() time.Time { return base }
	// access unexported constructor via the exported New with a tiny TTL,
	// then advance time by manipulating the entry indirectly through Set.
	// Use the internal test helper path via build tag instead.
	// Here we test via New with a 1 ns TTL and sleep.
	c := portcache.New(time.Nanosecond)
	c.Set(443, false)
	time.Sleep(2 * time.Millisecond)
	_, ok := c.Get(443)
	if ok {
		t.Error("expected cache miss after TTL expiry")
	}
	_ = clock
	_ = base
}

func TestSetOverwritesPreviousEntry(t *testing.T) {
	c := portcache.New(time.Minute)
	c.Set(22, false)
	c.Set(22, true)
	e, ok := c.Get(22)
	if !ok {
		t.Fatal("expected cache hit")
	}
	if !e.Open {
		t.Error("expected updated open=true")
	}
}

func TestEvictRemovesExpiredEntries(t *testing.T) {
	c := portcache.New(time.Nanosecond)
	c.Set(8080, true)
	c.Set(8443, true)
	time.Sleep(2 * time.Millisecond)
	if c.Len() != 2 {
		t.Fatalf("expected 2 entries before eviction, got %d", c.Len())
	}
	c.Evict()
	if c.Len() != 0 {
		t.Errorf("expected 0 entries after eviction, got %d", c.Len())
	}
}

func TestEvictKeepsValidEntries(t *testing.T) {
	c := portcache.New(time.Minute)
	c.Set(53, true)
	c.Evict()
	if c.Len() != 1 {
		t.Errorf("expected 1 valid entry after eviction, got %d", c.Len())
	}
}

func TestLenReflectsAllStoredEntries(t *testing.T) {
	c := portcache.New(time.Minute)
	for _, p := range []uint16{80, 443, 8080} {
		c.Set(p, true)
	}
	if c.Len() != 3 {
		t.Errorf("expected Len 3, got %d", c.Len())
	}
}
