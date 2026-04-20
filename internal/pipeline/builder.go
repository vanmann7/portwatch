package pipeline

import (
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/baseline"
	"github.com/user/portwatch/internal/cooldown"
	"github.com/user/portwatch/internal/dedupe"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/resolve"
	"github.com/user/portwatch/internal/rollup"
	"github.com/user/portwatch/internal/suppress"
)

// BuildOptions carries optional overrides for the assembled pipeline.
type BuildOptions struct {
	FilterRules   []string
	RollupWindow  time.Duration
	CooldownTTL   time.Duration
	DedupeWindow  time.Duration
	BaselinePath  string
	SuppressPath  string
}

// Build assembles a fully-wired Pipeline from the provided options and
// returns it along with the rollup channel for downstream consumers.
//
// Callers own calling Pipeline.Stop and Roller.Stop when done.
func Build(opts BuildOptions) (*Pipeline, *rollup.Roller, <-chan rollup.Batch, error) {
	if opts.RollupWindow == 0 {
		opts.RollupWindow = 300 * time.Millisecond
	}
	if opts.CooldownTTL == 0 {
		opts.CooldownTTL = 5 * time.Second
	}
	if opts.DedupeWindow == 0 {
		opts.DedupeWindow = 2 * time.Second
	}

	f, err := filter.New(opts.FilterRules)
	if err != nil {
		return nil, nil, nil, err
	}

	bl := baseline.New(opts.BaselinePath)
	sp := suppress.New(opts.SuppressPath)
	cd := cooldown.New(opts.CooldownTTL)
	dd := dedupe.New(opts.DedupeWindow)
	rv := resolve.New(nil)
	al := alert.New(nil)

	roller, batches := rollup.New(opts.RollupWindow)

	p := New(
		WithFilter(f),
		WithBaseline(bl),
		WithSuppress(sp),
		WithCooldown(cd),
		WithDedupe(dd),
		WithResolver(rv),
		WithAlerter(al),
		WithRoller(roller),
	)

	return p, roller, batches, nil
}
