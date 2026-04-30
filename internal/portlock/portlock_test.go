package portlock_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/portlock"
)

func TestLockAndIsLocked(t *testing.T) {
	l := portlock.New("")
	l.Lock(8080, "test reason")
	if !l.IsLocked(8080) {
		t.Fatal("expected port 8080 to be locked")
	}
}

func TestIsLockedFalseForUnknown(t *testing.T) {
	l := portlock.New("")
	if l.IsLocked(9999) {
		t.Fatal("expected port 9999 to be unlocked")
	}
}

func TestUnlockRemovesPort(t *testing.T) {
	l := portlock.New("")
	l.Lock(443, "https")
	if !l.Unlock(443) {
		t.Fatal("expected Unlock to return true")
	}
	if l.IsLocked(443) {
		t.Fatal("expected port 443 to be unlocked after Unlock")
	}
}

func TestUnlockReturnsFalseWhenNotLocked(t *testing.T) {
	l := portlock.New("")
	if l.Unlock(22) {
		t.Fatal("expected Unlock to return false for unknown port")
	}
}

func TestReasonReturnsAssociatedText(t *testing.T) {
	l := portlock.New("")
	l.Lock(22, "ssh access")
	r, ok := l.Reason(22)
	if !ok {
		t.Fatal("expected reason to be found")
	}
	if r != "ssh access" {
		t.Fatalf("expected 'ssh access', got %q", r)
	}
}

func TestReasonReturnsFalseForUnknown(t *testing.T) {
	l := portlock.New("")
	_, ok := l.Reason(1234)
	if ok {
		t.Fatal("expected reason not found for unknown port")
	}
}

func TestSnapshotReturnsCopy(t *testing.T) {
	l := portlock.New("")
	l.Lock(80, "http")
	l.Lock(443, "https")
	s := l.Snapshot()
	if len(s) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(s))
	}
	// Mutating the snapshot must not affect the locker.
	delete(s, 80)
	if !l.IsLocked(80) {
		t.Fatal("snapshot mutation affected the locker")
	}
}

func TestSaveLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "locks.json")

	a := portlock.New(path)
	a.Lock(22, "ssh")
	a.Lock(3306, "mysql")
	if err := a.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	b := portlock.New(path)
	if err := b.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	for _, port := range []int{22, 3306} {
		if !b.IsLocked(port) {
			t.Errorf("expected port %d to be locked after Load", port)
		}
	}
}

func TestLoadMissingFileIsNoop(t *testing.T) {
	l := portlock.New(filepath.Join(t.TempDir(), "missing.json"))
	if err := l.Load(); err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
}

func TestSaveNoPathIsNoop(t *testing.T) {
	l := portlock.New("")
	l.Lock(80, "http")
	if err := l.Save(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSaveWritesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "locks.json")
	l := portlock.New(path)
	l.Lock(8080, "dev")
	if err := l.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}
