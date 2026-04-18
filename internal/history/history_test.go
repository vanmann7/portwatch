package history

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRecordAndEntries(t *testing.T) {
	h := New(10)
	h.Record([]uint16{80, 443}, nil)
	h.Record(nil, []uint16{8080})

	entries := h.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if len(entries[0].Opened) != 2 {
		t.Errorf("expected 2 opened ports in first entry")
	}
	if len(entries[1].Closed) != 1 {
		t.Errorf("expected 1 closed port in second entry")
	}
}

func TestMaxSizeEviction(t *testing.T) {
	h := New(3)
	for i := 0; i < 5; i++ {
		h.Record([]uint16{uint16(i)}, nil)
	}
	if len(h.Entries()) != 3 {
		t.Errorf("expected 3 entries after eviction, got %d", len(h.Entries()))
	}
}

func TestDefaultMaxSize(t *testing.T) {
	h := New(0)
	if h.maxSize != 100 {
		t.Errorf("expected default maxSize 100, got %d", h.maxSize)
	}
}

func TestSaveLoad(t *testing.T) {
	h := New(10)
	h.Record([]uint16{22}, nil)
	h.Record(nil, []uint16{3306})

	dir := t.TempDir()
	path := filepath.Join(dir, "history.json")

	if err := h.Save(path); err != nil {
		t.Fatalf("Save: %v", err)
	}

	h2 := New(10)
	if err := h2.Load(path); err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(h2.Entries()) != 2 {
		t.Errorf("expected 2 entries after load, got %d", len(h2.Entries()))
	}
}

func TestLoadMissingFile(t *testing.T) {
	h := New(10)
	err := h.Load("/nonexistent/path/history.json")
	if !os.IsNotExist(err) {
		t.Errorf("expected not-exist error, got %v", err)
	}
}
