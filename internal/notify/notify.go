package notify

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Channel represents a notification delivery channel.
type Channel interface {
	Send(subject, body string) error
}

// StdoutChannel writes notifications to an io.Writer.
type StdoutChannel struct {
	Writer io.Writer
}

// NewStdout creates a StdoutChannel writing to os.Stdout by default.
func NewStdout(w io.Writer) *StdoutChannel {
	if w == nil {
		w = os.Stdout
	}
	return &StdoutChannel{Writer: w}
}

// Send writes a formatted notification message.
func (s *StdoutChannel) Send(subject, body string) error {
	_, err := fmt.Fprintf(s.Writer, "[%s] %s: %s\n", time.Now().Format(time.RFC3339), subject, body)
	return err
}

// Dispatcher fans out notifications to multiple channels.
type Dispatcher struct {
	channels []Channel
}

// NewDispatcher creates a Dispatcher with the given channels.
func NewDispatcher(channels ...Channel) *Dispatcher {
	return &Dispatcher{channels: channels}
}

// Dispatch sends subject/body to all registered channels, collecting errors.
func (d *Dispatcher) Dispatch(subject, body string) []error {
	var errs []error
	for _, ch := range d.channels {
		if err := ch.Send(subject, body); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
