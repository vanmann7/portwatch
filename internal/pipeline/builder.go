package pipeline

import (
	"io"
	"time"

	"github.com/user/portwatch/internal/aggregator"
	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/dedupe"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/resolve"
	"github.com/user/portwatch/internal/tagger"
)

// BuildOptions carries optional overrides used when constructing the pipeline.
type BuildOptions struct {
	// AlertWriter overrides the default stdout writer for alerts.
	AlertWriter io.Writer
	// AggregateWindow overrides the default 2-second aggregation window.
	AggregateWindow time.Duration
}

// Build constructs a fully wired [Pipeline] from the supplied configuration.
// It wires together the filter, resolver, tagger, deduplicator, aggregator,
// and alert notifier according to cfg.
func Build(cfg *config.Config, opts BuildOptions) (*Pipeline, error) {
	if opts.AggregateWindow == 0 {
		opts.AggregateWindow = 2 * time.Second
	}

	f, err := filter.New(cfg.Rules)
	if err != nil {
		return nil, err
	}

	r := resolve.New(nil)
	tg := tagger.New(nil)
	dd := dedupe.New(5 * time.Second)
	ag := aggregator.New(opts.AggregateWindow)

	var notifier *alert.Notifier
	if opts.AlertWriter != nil {
		notifier = alert.New(opts.AlertWriter)
	} else {
		notifier = alert.New(nil)
	}

	return New(f, r, tg, dd, ag, notifier), nil
}
