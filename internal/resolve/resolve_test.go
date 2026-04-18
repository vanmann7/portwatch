package resolve_test

import (
	"testing"

	"github.com/user/portwatch/internal/resolve"
)

func TestKnownPort(t *testing.T) {
	r := resolve.New(nil)
	if got := r.Name(22); got != "ssh" {
		t.Fatalf("expected ssh, got %q", got)
	}
}

func TestHTTPSPort(t *testing.T) {
	r := resolve.New(nil)
	if got := r.Name(443); got != "https" {
		t.Fatalf("expected https, got %q", got)
	}
}

func TestUnknownPortFallback(t *testing.T) {
	r := resolve.New(nil)
	got := r.Name(19999)
	if got == "" {
		t.Fatal("expected non-empty name for unknown port")
	}
}

func TestOverrideTakesPrecedence(t *testing.T) {
	overrides := map[uint16]string{80: "my-app"}
	r := resolve.New(overrides)
	if got := r.Name(80); got != "my-app" {
		t.Fatalf("expected my-app, got %q", got)
	}
}

func TestCacheReturnsSameValue(t *testing.T) {
	r := resolve.New(nil)
	first := r.Name(22)
	second := r.Name(22)
	if first != second {
		t.Fatalf("cache inconsistency: %q vs %q", first, second)
	}
}

func TestResetClearsCache(t *testing.T) {
	r := resolve.New(nil)
	_ = r.Name(22)
	r.Reset()
	// After reset the lookup should still succeed (rebuilds from table).
	if got := r.Name(22); got != "ssh" {
		t.Fatalf("expected ssh after reset, got %q", got)
	}
}

func TestMultiplePorts(t *testing.T) {
	r := resolve.New(nil)
	cases := map[uint16]string{
		21:   "ftp",
		25:   "smtp",
		3306: "mysql",
		6379: "redis",
	}
	for port, want := range cases {
		if got := r.Name(port); got != want {
			t.Errorf("port %d: expected %q, got %q", port, want, got)
		}
	}
}
