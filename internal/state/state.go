// Package state manages port state snapshots and change detection.
package state

import (
	"encoding/json"
	"os"
	"time"
)

// Snapshot holds a recorded set of open ports at a point in time.
type Snapshot struct {
	Timestamp time.Time `json:"timestamp"`
	Ports     []int     `json:"ports"`
}

// Diff represents changes between two snapshots.
type Diff struct {
	Opened []int
	Closed []int
}

// Compare returns the diff between a previous and current snapshot.
func Compare(prev, curr Snapshot) Diff {
	prevSet := toSet(prev.Ports)
	currSet := toSet(curr.Ports)

	var d Diff
	for p := range currSet {
		if !prevSet[p] {
			d.Opened = append(d.Opened, p)
		}
	}
	for p := range prevSet {
		if !currSet[p] {
			d.Closed = append(d.Closed, p)
		}
	}
	return d
}

// Save writes a snapshot to a JSON file.
func Save(path string, s Snapshot) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(s)
}

// Load reads a snapshot from a JSON file.
func Load(path string) (Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		return Snapshot{}, err
	}
	defer f.Close()
	var s Snapshot
	return s, json.NewDecoder(f).Decode(&s)
}

func toSet(ports []int) map[int]bool {
	m := make(map[int]bool, len(ports))
	for _, p := range ports {
		m[p] = true
	}
	return m
}
