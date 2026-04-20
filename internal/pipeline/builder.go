package pipeline

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/dedupe"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/notify"
)

// Config holds the options used by Build to construct a Pipeline.
type Config struct {
	// FilterRules are port/range expressions passed to filter.New.
	// A nil or empty slice allows all ports.
	FilterRules []string

	// DedupeWindow is the time window within which identical events are
	// suppressed. Defaults to 1 minute when zero.
	DedupeWindow time.Duration

	// Writer is the io.Writer used for alert and stdout notifications.
	// Defaults to os.Stdout when nil.
	Writer io.Writer

	// LogPath, when non-empty, adds a file notify channel.
	LogPath string
}

// Build constructs a ready-to-use Pipeline from cfg.
func Build(cfg Config) (*Pipeline, error) {
	if cfg.DedupeWindow == 0 {
		cfg.DedupeWindow = time.Minute
	}
	if cfg.Writer == nil {
		cfg.Writer = os.Stdout
	}

	f, err := filter.New(cfg.FilterRules)
	if err != nil {
		return nil, fmt.Errorf("pipeline: filter: %w", err)
	}

	d := dedupe.New(cfg.DedupeWindow)
	a := alert.New(cfg.Writer)
	m := metrics.New()

	channels := []notify.Channel{notify.NewStdout(cfg.Writer)}
	if cfg.LogPath != "" {
		fc, err := notify.NewFileChannel(cfg.LogPath)
		if err != nil {
			return nil, fmt.Errorf("pipeline: file channel: %w", err)
		}
		channels = append(channels, fc)
	}

	disp := notify.NewDispatcher(channels)
	return New(f, d, a, disp, m), nil
}
