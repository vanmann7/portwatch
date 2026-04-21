package portname_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/portname"
)

func TestBuiltinHTTP(t *testing.T) {
	r := portname.New(nil)
	if got := r.Lookup(80); got != "http" {
		t.Fatalf("expected http, got %q", got)
	}
}

func TestBuiltinSSH(t *testing.T) {
	r := portname.New(nil)
	if got := r.Lookup(22); got != "ssh" {
		t.Fatalf("expected ssh, got %q", got)
	}
}

func TestUnknownPortFallback(t *testing.T) {
	r := portname.New(nil)
	if got := r.Lookup(9999); got != "unknown" {
		t.Fatalf("expected unknown, got %q", got)
	}
}

func TestOverrideTakesPrecedence(t *testing.T) {
	r := portname.New(map[int]string{80: "my-service"})
	if got := r.Lookup(80); got != "my-service" {
		t.Fatalf("expected my-service, got %q", got)
	}
}

func TestOverrideNewPort(t *testing.T) {
	r := portname.New(map[int]string{12345: "custom-app"})
	if got := r.Lookup(12345); got != "custom-app" {
		t.Fatalf("expected custom-app, got %q", got)
	}
}

func TestIsKnownBuiltin(t *testing.T) {
	r := portname.New(nil)
	if !r.IsKnown(443) {
		t.Fatal("expected port 443 to be known")
	}
}

func TestIsKnownUnknown(t *testing.T) {
	r := portname.New(nil)
	if r.IsKnown(19999) {
		t.Fatal("expected port 19999 to be unknown")
	}
}

func TestIsKnownViaOverride(t *testing.T) {
	r := portname.New(map[int]string{19999: "internal-api"})
	if !r.IsKnown(19999) {
		t.Fatal("expected port 19999 to be known via override")
	}
}

func TestNilOverrideDoesNotPanic(t *testing.T) {
	r := portname.New(nil)
	_ = r.Lookup(22)
}
