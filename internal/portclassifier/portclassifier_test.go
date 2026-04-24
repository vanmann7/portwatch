package portclassifier_test

import (
	"testing"

	"github.com/example/portwatch/internal/portclassifier"
)

func TestSystemPort(t *testing.T) {
	c, err := portclassifier.New(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := c.Classify(22); got != portclassifier.CategorySystem {
		t.Errorf("port 22: got %q, want %q", got, portclassifier.CategorySystem)
	}
}

func TestUserPort(t *testing.T) {
	c, _ := portclassifier.New(nil)
	if got := c.Classify(8080); got != portclassifier.CategoryUser {
		t.Errorf("port 8080: got %q, want %q", got, portclassifier.CategoryUser)
	}
}

func TestDynamicPort(t *testing.T) {
	c, _ := portclassifier.New(nil)
	if got := c.Classify(60000); got != portclassifier.CategoryDynamic {
		t.Errorf("port 60000: got %q, want %q", got, portclassifier.CategoryDynamic)
	}
}

func TestOverrideTakesPrecedence(t *testing.T) {
	c, err := portclassifier.New(map[uint16]string{22: "override"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := c.Classify(22); got != portclassifier.CategoryOverride {
		t.Errorf("port 22: got %q, want %q", got, portclassifier.CategoryOverride)
	}
}

func TestInvalidOverrideReturnsError(t *testing.T) {
	_, err := portclassifier.New(map[uint16]string{80: "bogus"})
	if err == nil {
		t.Fatal("expected error for invalid category, got nil")
	}
}

func TestClassifyAll(t *testing.T) {
	c, _ := portclassifier.New(nil)
	result := c.ClassifyAll([]uint16{80, 8080, 55000})
	if result[80] != portclassifier.CategorySystem {
		t.Errorf("port 80: got %q", result[80])
	}
	if result[8080] != portclassifier.CategoryUser {
		t.Errorf("port 8080: got %q", result[8080])
	}
	if result[55000] != portclassifier.CategoryDynamic {
		t.Errorf("port 55000: got %q", result[55000])
	}
}

func TestBoundaryPort1023(t *testing.T) {
	c, _ := portclassifier.New(nil)
	if got := c.Classify(1023); got != portclassifier.CategorySystem {
		t.Errorf("port 1023: got %q", got)
	}
}

func TestBoundaryPort1024(t *testing.T) {
	c, _ := portclassifier.New(nil)
	if got := c.Classify(1024); got != portclassifier.CategoryUser {
		t.Errorf("port 1024: got %q", got)
	}
}
