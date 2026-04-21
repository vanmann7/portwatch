// Package formatter provides structured text and JSON formatting
// for port change events before they are dispatched to notifiers.
package formatter

import (
	"encoding/json"
	"fmt"
	"time"
)

// Format describes the output format for an event.
type Format int

const (
	FormatText Format = iota
	FormatJSON
)

// Event represents a port change event to be formatted.
type Event struct {
	Port      int
	Proto     string
	Action    string // "opened" or "closed"
	Service   string
	Timestamp time.Time
}

// Formatter converts Events into formatted strings.
type Formatter struct {
	fmt Format
}

// New returns a new Formatter for the given Format.
func New(f Format) *Formatter {
	return &Formatter{fmt: f}
}

// Format converts an Event to a string in the configured format.
func (f *Formatter) Format(e Event) (string, error) {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now().UTC()
	}
	switch f.fmt {
	case FormatJSON:
		return formatJSON(e)
	default:
		return formatText(e), nil
	}
}

func formatText(e Event) string {
	service := e.Service
	if service == "" {
		service = "unknown"
	}
	return fmt.Sprintf("%s  port %d/%s (%s) %s",
		e.Timestamp.Format(time.RFC3339),
		e.Port,
		e.Proto,
		service,
		e.Action,
	)
}

func formatJSON(e Event) (string, error) {
	payload := struct {
		Timestamp string `json:"timestamp"`
		Port      int    `json:"port"`
		Proto     string `json:"proto"`
		Action    string `json:"action"`
		Service   string `json:"service,omitempty"`
	}{
		Timestamp: e.Timestamp.Format(time.RFC3339),
		Port:      e.Port,
		Proto:     e.Proto,
		Action:    e.Action,
		Service:   e.Service,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("formatter: marshal: %w", err)
	}
	return string(b), nil
}
