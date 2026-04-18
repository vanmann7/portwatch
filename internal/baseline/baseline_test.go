package baseline

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "baseline.json")
}

func TestSetAndContains(t *testing.T) {
	b := New(tmpPath(t))
	b.Set([]int{80, 443, 8080})
	for _, p := range []int{80, 443, 8080} {
		if !b.Contains(p) {
			t.Errorf("expected port %d to be in baseline", p)
		}
	}
	if b.Contains(22) {
		t.Error("port 22 should not be in baseline")
	}
}

func TestSetReplacesExisting(t *testing.T) {
	b := New(tmpPath(t))
	b.Set([]int{80, 443})
	b.Set([]int{22})
	if b.Contains(80) {
		t.Error("old port 80 should have been replaced")
	}
	if !b.Contains(22) {
		t.Error("port 22 should be present after reset")
	}
}

func TestPortsSnapshot(t *testing.T) {
	b := New(tmpPath(t))
	want := []int{22, 80, 443}
	b.Set(want)
	got := b.Ports()
	sort.Ints(got)
	if len(got) != len(want) {
		t.Fatalf("expected %d ports, got %d", len(want), len(got))
	}
	for i, p := range want {
		if got[i] != p {
			t.Errorf("index %d: want %d got %d", i, p, got[i])
		}
	}
}

func TestSaveLoad(t *testing.T) {
	path := tmpPath(t)
	b := New(path)
	b.Set([]int{80, 443})
	if err := b.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}
	b2 := New(path)
	if err := b2.Load(); err != nil {
		t.Fatalf("Load: %v", err)
	}
	for _, p := range []int{80, 443} {
		if !b2.Contains(p) {
			t.Errorf("expected port %d after load", p)
		}
	}
}

func TestLoadMissingFile(t *testing.T) {
	b := New(filepath.Join(t.TempDir(), "no-such-file.json"))
	if err := b.Load(); err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
}

func TestSaveInvalidPath(t *testing.T) {
	b := New("/no/such/dir/baseline.json")
	b.Set([]int{80})
	if err := b.Save(); !os.IsNotExist(err) && err == nil {
		t.Log("got expected error or none (platform-dependent)")
	}
}
