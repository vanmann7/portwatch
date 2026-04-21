package eventfilter_test

import (
	"testing"

	"github.com/user/portwatch/internal/eventfilter"
	"github.com/user/portwatch/internal/pipeline"
)

func makeEvent(port int, t pipeline.EventType) pipeline.Event {
	return pipeline.Event{Port: port, Type: t}
}

func TestAllowNoPredicatesPassesAll(t *testing.T) {
	f := eventfilter.New()
	if !f.Allow(makeEvent(80, pipeline.EventOpened)) {
		t.Fatal("expected event to pass with no predicates")
	}
}

func TestMinPortFiltersLow(t *testing.T) {
	f := eventfilter.New(eventfilter.MinPort(1024))
	if f.Allow(makeEvent(80, pipeline.EventOpened)) {
		t.Error("port 80 should be rejected by MinPort(1024)")
	}
	if !f.Allow(makeEvent(8080, pipeline.EventOpened)) {
		t.Error("port 8080 should pass MinPort(1024)")
	}
}

func TestMaxPortFiltersHigh(t *testing.T) {
	f := eventfilter.New(eventfilter.MaxPort(1023))
	if !f.Allow(makeEvent(80, pipeline.EventOpened)) {
		t.Error("port 80 should pass MaxPort(1023)")
	}
	if f.Allow(makeEvent(8080, pipeline.EventOpened)) {
		t.Error("port 8080 should be rejected by MaxPort(1023)")
	}
}

func TestOnlyOpenedFiltersClosedEvents(t *testing.T) {
	f := eventfilter.New(eventfilter.OnlyOpened())
	if f.Allow(makeEvent(443, pipeline.EventClosed)) {
		t.Error("closed event should be rejected by OnlyOpened")
	}
	if !f.Allow(makeEvent(443, pipeline.EventOpened)) {
		t.Error("opened event should pass OnlyOpened")
	}
}

func TestOnlyClosedFiltersOpenedEvents(t *testing.T) {
	f := eventfilter.New(eventfilter.OnlyClosed())
	if f.Allow(makeEvent(22, pipeline.EventOpened)) {
		t.Error("opened event should be rejected by OnlyClosed")
	}
	if !f.Allow(makeEvent(22, pipeline.EventClosed)) {
		t.Error("closed event should pass OnlyClosed")
	}
}

func TestChainedPredicatesAllMustPass(t *testing.T) {
	f := eventfilter.New(
		eventfilter.MinPort(1024),
		eventfilter.MaxPort(9000),
		eventfilter.OnlyOpened(),
	)
	if !f.Allow(makeEvent(8080, pipeline.EventOpened)) {
		t.Error("8080/opened should pass all predicates")
	}
	if f.Allow(makeEvent(80, pipeline.EventOpened)) {
		t.Error("80/opened should fail MinPort")
	}
	if f.Allow(makeEvent(9001, pipeline.EventOpened)) {
		t.Error("9001/opened should fail MaxPort")
	}
	if f.Allow(makeEvent(8080, pipeline.EventClosed)) {
		t.Error("8080/closed should fail OnlyOpened")
	}
}

func TestAddAppendsPredicateDynamically(t *testing.T) {
	f := eventfilter.New(eventfilter.MinPort(1024))
	f.Add(eventfilter.OnlyOpened())
	if f.Allow(makeEvent(8080, pipeline.EventClosed)) {
		t.Error("closed event should be rejected after Add(OnlyOpened)")
	}
	if !f.Allow(makeEvent(8080, pipeline.EventOpened)) {
		t.Error("opened event should still pass after Add")
	}
}
