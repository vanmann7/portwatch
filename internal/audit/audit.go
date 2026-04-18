// Package audit provides structured audit logging for port change events.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     string    `json:"event"`
	Port      int       `json:"port"`
	Protocol  string    `json:"protocol"`
	Note      string    `json:"note,omitempty"`
}

// Logger writes audit entries as newline-delimited JSON.
type Logger struct {
	w io.Writer
}

// New returns a Logger writing to w. Pass nil to use os.Stdout.
func New(w io.Writer) *Logger {
	if w == nil {
		w = os.Stdout
	}
	return &Logger{w: w}
}

// Record writes a single audit entry.
func (l *Logger) Record(event string, port int, protocol, note string) error {
	e := Entry{
		Timestamp: time.Now().UTC(),
		Event:     event,
		Port:      port,
		Protocol:  protocol,
		Note:      note,
	}
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", b)
	return err
}

// Opened is a convenience wrapper for a port-opened event.
func (l *Logger) Opened(port int, protocol string) error {
	return l.Record("opened", port, protocol, "")
}

// Closed is a convenience wrapper for a port-closed event.
func (l *Logger) Closed(port int, protocol string) error {
	return l.Record("closed", port, protocol, "")
}
