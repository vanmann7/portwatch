package suppress_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/suppress"
)

func TestAddAndIsSuppressed(t *testing.T) {
	l := suppress.New()
	l.Add(80, 443)
	if !l.IsSuppressed(80) {
		t.Error("expected 80 to be suppressed")
	}
	if !l.IsSuppressed(443) {
		t.Error("expected 443 to be suppressed")
	}
	if l.IsSuppressed(8080) {
		t.Error("expected 8080 not to be suppressed")
	}
}

func TestRemove(t *testing.T) {
	l := suppress.New()
	l.Add(22)
	l.Remove(22)
	if l.IsSuppressed(22) {
		t.Error("expected 22 to be removed")
	}
}

func TestPortsSnapshot(t *testing.T) {
	l := suppress.New()
	l.Add(1, 2, 3)
	ports := l.Ports()
	if len(ports) != 3 {
		t.Fatalf("expected 3 ports, got %d", len(ports))
	}
}

func TestEmptyListNotSuppressed(t *testing.T) {
	l := suppress.New()
	if l.IsSuppressed(80) {
		t.Error("empty list should not suppress anything")
	}
}

func TestSaveLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "suppress.json")

	l := suppress.New()
	l.Add(22, 80, 443)
	if err := l.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	l2, err := suppress.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	for _, p := range []int{22, 80, 443} {
		if !l2.IsSuppressed(p) {
			t.Errorf("expected port %d to be suppressed after load", p)
		}
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := suppress.Load("/nonexistent/suppress.json")
	if !os.IsNotExist(err) {
		t.Errorf("expected not-exist error, got %v", err)
	}
}
