package portclassifier_test

import (
	"strings"
	"testing"

	"github.com/example/portwatch/internal/portclassifier"
)

func TestParseOverridesValid(t *testing.T) {
	input := `
# comment
22 : system
8080: user
`
	overrides, err := portclassifier.ParseOverrides(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if overrides[22] != "system" {
		t.Errorf("port 22: got %q", overrides[22])
	}
	if overrides[8080] != "user" {
		t.Errorf("port 8080: got %q", overrides[8080])
	}
}

func TestParseOverridesBlankAndComments(t *testing.T) {
	input := "\n# just a comment\n\n"
	overrides, err := portclassifier.ParseOverrides(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(overrides) != 0 {
		t.Errorf("expected empty map, got %v", overrides)
	}
}

func TestParseOverridesMissingColon(t *testing.T) {
	_, err := portclassifier.ParseOverrides(strings.NewReader("80 system"))
	if err == nil {
		t.Fatal("expected error for missing colon")
	}
}

func TestParseOverridesInvalidPort(t *testing.T) {
	_, err := portclassifier.ParseOverrides(strings.NewReader("notaport:system"))
	if err == nil {
		t.Fatal("expected error for invalid port")
	}
}

func TestParseOverridesRoundTrip(t *testing.T) {
	input := "443:override\n"
	overrides, err := portclassifier.ParseOverrides(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c, err := portclassifier.New(overrides)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if got := c.Classify(443); got != portclassifier.CategoryOverride {
		t.Errorf("port 443: got %q", got)
	}
}
