package tagger_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/tagger"
)

func TestWellKnownPort(t *testing.T) {
	tr := tagger.New(nil)
	got := tr.Tag(22)
	if got != "ssh" {
		t.Fatalf("expected ssh, got %q", got)
	}
}

func TestHTTPSPort(t *testing.T) {
	tr := tagger.New(nil)
	if got := tr.Tag(443); got != "https" {
		t.Fatalf("expected https, got %q", got)
	}
}

func TestUnknownPortFallback(t *testing.T) {
	tr := tagger.New(nil)
	got := tr.Tag(9999)
	if got != "port-9999" {
		t.Fatalf("expected port-9999, got %q", got)
	}
}

func TestOverrideTakesPrecedence(t *testing.T) {
	overrides := map[int]string{80: "my-web"}
	tr := tagger.New(overrides)
	if got := tr.Tag(80); got != "my-web" {
		t.Fatalf("expected my-web, got %q", got)
	}
}

func TestOverrideUnknownPort(t *testing.T) {
	overrides := map[int]string{12345: "custom-svc"}
	tr := tagger.New(overrides)
	if got := tr.Tag(12345); got != "custom-svc" {
		t.Fatalf("expected custom-svc, got %q", got)
	}
}

func TestTagAll(t *testing.T) {
	tr := tagger.New(nil)
	ports := []int{22, 80, 9999}
	tags := tr.TagAll(ports)

	expected := map[int]string{
		22:   "ssh",
		80:   "http",
		9999: "port-9999",
	}
	for port, want := range expected {
		if got := tags[port]; got != want {
			t.Errorf("port %d: expected %q, got %q", port, want, got)
		}
	}
}

func TestTagAllEmpty(t *testing.T) {
	tr := tagger.New(nil)
	result := tr.TagAll([]int{})
	if len(result) != 0 {
		t.Fatalf("expected empty map, got %v", result)
	}
}

func TestNilOverridesDoesNotPanic(t *testing.T) {
	tr := tagger.New(nil)
	if tr == nil {
		t.Fatal("expected non-nil Tagger")
	}
	_ = tr.Tag(443)
}
