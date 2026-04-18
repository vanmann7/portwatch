// Package alert provides alerting functionality for portwatch,
// formatting and emitting notifications when port state changes are detected.
package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/state"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
)

// Alert holds a single alert message.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Message   string
}

// Notifier writes alerts to an output destination.
type Notifier struct {
	out io.Writer
}

// New returns a Notifier that writes to the given writer.
// If w is nil, os.Stdout is used.
func New(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{out: w}
}

// Notify emits alerts for each change in the provided Diff.
// Opened ports are reported at WARN level; closed ports at INFO level.
func (n *Notifier) Notify(diff state.Diff) []Alert {
	var alerts []Alert

	for _, port := range diff.Opened {
		a := Alert{
			Timestamp: time.Now(),
			Level:     LevelWarn,
			Message:   fmt.Sprintf("port %d is now OPEN (unexpected)", port),
		}
		alerts = append(alerts, a)
		fmt.Fprintf(n.out, "[%s] %s %s\n", a.Timestamp.Format(time.RFC3339), a.Level, a.Message)
	}

	for _, port := range diff.Closed {
		a := Alert{
			Timestamp: time.Now(),
			Level:     LevelInfo,
			Message:   fmt.Sprintf("port %d is now CLOSED", port),
		}
		alerts = append(alerts, a)
		fmt.Fprintf(n.out, "[%s] %s %s\n", a.Timestamp.Format(time.RFC3339), a.Level, a.Message)
	}

	return alerts
}
