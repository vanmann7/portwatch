package eventlog_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/eventlog"
)

func TestRecordAndEntries(t *testing.T) {
	log := eventlog.New(10)
	log.Record(80, eventlog.EventOpened, "http")
	log.Record(443, eventlog.EventOpened, "https")

	entries := log.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Port != 80 || entries[0].Kind != eventlog.EventOpened {
		t.Errorf("unexpected first entry: %+v", entries[0])
	}
	if entries[1].Port != 443 || entries[1].Service != "https" {
		t.Errorf("unexpected second entry: %+v", entries[1])
	}
}

func TestRingBufferEvictsOldest(t *testing.T) {
	log := eventlog.New(3)
	log.Record(1, eventlog.EventOpened, "")
	log.Record(2, eventlog.EventOpened, "")
	log.Record(3, eventlog.EventOpened, "")
	log.Record(4, eventlog.EventOpened, "") // evicts port 1

	entries := log.Entries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	if entries[0].Port != 2 {
		t.Errorf("expected oldest to be port 2, got %d", entries[0].Port)
	}
	if entries[2].Port != 4 {
		t.Errorf("expected newest to be port 4, got %d", entries[2].Port)
	}
}

func TestDefaultCapacity(t *testing.T) {
	log := eventlog.New(0)
	for i := 0; i < 300; i++ {
		log.Record(i, eventlog.EventOpened, "")
	}
	if log.Len() != eventlog.DefaultCapacity {
		t.Errorf("expected len %d, got %d", eventlog.DefaultCapacity, log.Len())
	}
}

func TestClearResetsLog(t *testing.T) {
	log := eventlog.New(5)
	log.Record(22, eventlog.EventOpened, "ssh")
	log.Clear()

	if log.Len() != 0 {
		t.Errorf("expected 0 entries after clear, got %d", log.Len())
	}
	if log.Entries() != nil {
		t.Error("expected nil entries after clear")
	}
}

func TestTimestampIsRecent(t *testing.T) {
	before := time.Now()
	log := eventlog.New(5)
	log.Record(8080, eventlog.EventClosed, "")
	after := time.Now()

	entries := log.Entries()
	if len(entries) != 1 {
		t.Fatal("expected 1 entry")
	}
	ts := entries[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v not within expected range [%v, %v]", ts, before, after)
	}
}
