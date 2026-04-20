// Package pipeline wires together the scan, filter, dedupe, and notify
// stages into a single reusable processing pipeline.
package pipeline

import (
	"context"
	"fmt"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/dedupe"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/notify"
	"github.com/user/portwatch/internal/state"
)

// Event represents a single port-change event emitted by the pipeline.
type Event struct {
	Port   int
	Opened bool
}

// Pipeline processes state diffs through filter → dedupe → alert → notify.
type Pipeline struct {
	filter  *filter.Filter
	dedupe  *dedupe.Deduper
	alert   *alert.Alerter
	disp    *notify.Dispatcher
	metrics *metrics.Tracker
}

// New constructs a Pipeline from the provided components.
func New(f *filter.Filter, d *dedupe.Deduper, a *alert.Alerter, disp *notify.Dispatcher, m *metrics.Tracker) *Pipeline {
	return &Pipeline{
		filter:  f,
		dedupe:  d,
		alert:   a,
		disp:    disp,
		metrics: m,
	}
}

// Process takes a state diff, runs it through the pipeline stages, and
// dispatches notifications. It returns any dispatcher errors.
func (p *Pipeline) Process(ctx context.Context, diff state.Diff) error {
	for _, port := range diff.Opened {
		if !p.filter.Allow(port) {
			continue
		}
		key := fmt.Sprintf("opened:%d", port)
		if p.dedupe.IsDuplicate(key) {
			continue
		}
		p.alert.Notify(diff)
		p.metrics.RecordScan(len(diff.Opened), len(diff.Closed))
		return p.disp.Dispatch(ctx)
	}
	for _, port := range diff.Closed {
		if !p.filter.Allow(port) {
			continue
		}
		key := fmt.Sprintf("closed:%d", port)
		if p.dedupe.IsDuplicate(key) {
			continue
		}
		p.alert.Notify(diff)
		p.metrics.RecordScan(len(diff.Opened), len(diff.Closed))
		return p.disp.Dispatch(ctx)
	}
	return nil
}
