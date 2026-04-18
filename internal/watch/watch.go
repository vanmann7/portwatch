// Package watch provides a periodic port-scanning loop that compares
// successive snapshots and fires alerts on any change.
package watch

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

// Watcher runs the monitoring loop.
type Watcher struct {
	cfg     *config.Config
	scanner *scanner.Scanner
	alerter *alert.Alerter
}

// New creates a Watcher from the provided configuration.
func New(cfg *config.Config, a *alert.Alerter) (*Watcher, error) {
	s, err := scanner.New(cfg.StartPort, cfg.EndPort)
	if err != nil {
		return nil, err
	}
	return &Watcher{cfg: cfg, scanner: s, alerter: a}, nil
}

// Run starts the watch loop, blocking until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.cfg.Interval)
	defer ticker.Stop()

	prev, err := w.scanner.Scan()
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			curr, err := w.scanner.Scan()
			if err != nil {
				log.Printf("scan error: %v", err)
				continue
			}
			diff := state.Compare(prev, curr)
			if err := w.alerter.Notify(diff); err != nil {
				log.Printf("alert error: %v", err)
			}
			prev = curr
		}
	}
}
