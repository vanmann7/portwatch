package state_test

import (
	"os"
	"sort"
	"testing"
	"time"

	"github.com/user/portwatch/internal/state"
)

func snap(ports ...int) state.Snapshot {
	return state.Snapshot{Timestamp: time.Now(), Ports: ports}
}

func TestCompareOpened(t *testing.T) {
	diff := state.Compare(snap(80, 443), snap(80, 443, 8080))
	if len(diff.Opened) != 1 || diff.Opened[0] != 8080 {
		t.Errorf("expected 8080 opened, got %v", diff.Opened)
	}
	if len(diff.Closed) != 0 {
		t.Errorf("expected no closed ports, got %v", diff.Closed)
	}
}

func TestCompareClosed(t *testing.T) {
	diff := state.Compare(snap(80, 443, 8080), snap(80, 443))
	if len(diff.Closed) != 1 || diff.Closed[0] != 8080 {
		t.Errorf("expected 8080 closed, got %v", diff.Closed)
	}
}

func TestCompareNoChange(t *testing.T) {
	diff := state.Compare(snap(80, 443), snap(80, 443))
	if len(diff.Opened) != 0 || len(diff.Closed) != 0 {
		t.Errorf("expected no diff, got %+v", diff)
	}
}

func TestSaveLoad(t *testing.T) {
	f, err := os.CreateTemp("", "portwatch-*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())

	orig := snap(22, 80, 443)
	if err := state.Save(f.Name(), orig); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := state.Load(f.Name())
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	sort.Ints(orig.Ports)
	sort.Ints(loaded.Ports)
	for i, p := range orig.Ports {
		if loaded.Ports[i] != p {
			t.Errorf("port mismatch at %d: want %d got %d", i, p, loaded.Ports[i])
		}
	}
}
