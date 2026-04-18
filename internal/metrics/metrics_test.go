package metrics_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

func TestRecordScanAccumulates(t *testing.T) {
	tr := metrics.New()
	tr.RecordScan(3, 1, 10*time.Millisecond)
	tr.RecordScan(2, 0, 5*time.Millisecond)

	snap := tr.Snapshot()
	if snap.TotalScans != 2 {
		t.Errorf("expected 2 scans, got %d", snap.TotalScans)
	}
	if snap.OpenedPorts != 5 {
		t.Errorf("expected 5 opened, got %d", snap.OpenedPorts)
	}
	if snap.ClosedPorts != 1 {
		t.Errorf("expected 1 closed, got %d", snap.ClosedPorts)
	}
	if snap.LastDuration != 5*time.Millisecond {
		t.Errorf("unexpected last duration: %v", snap.LastDuration)
	}
}

func TestSnapshotLastScanAt(t *testing.T) {
	tr := metrics.New()
	before := time.Now()
	tr.RecordScan(0, 0, 1*time.Millisecond)
	after := time.Now()

	snap := tr.Snapshot()
	if snap.LastScanAt.Before(before) || snap.LastScanAt.After(after) {
		t.Errorf("LastScanAt %v out of expected range", snap.LastScanAt)
	}
}

func TestReset(t *testing.T) {
	tr := metrics.New()
	tr.RecordScan(5, 3, 20*time.Millisecond)
	tr.Reset()

	snap := tr.Snapshot()
	if snap.TotalScans != 0 || snap.OpenedPorts != 0 || snap.ClosedPorts != 0 {
		t.Errorf("expected zeroed metrics after reset, got %+v", snap)
	}
}

func TestNewTrackerIsEmpty(t *testing.T) {
	tr := metrics.New()
	snap := tr.Snapshot()
	if snap.TotalScans != 0 {
		t.Errorf("new tracker should have zero scans")
	}
	if !snap.LastScanAt.IsZero() {
		t.Errorf("new tracker LastScanAt should be zero")
	}
}
