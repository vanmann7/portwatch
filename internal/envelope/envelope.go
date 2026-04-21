// Package envelope wraps a port-change event with routing metadata so
// that downstream consumers (webhooks, file channels, stdout) can make
// dispatch decisions without inspecting the raw event payload.
package envelope

import "time"

// Destination identifies where an event should be delivered.
type Destination uint8

const (
	// DestStdout writes the event to standard output.
	DestStdout Destination = 1 << iota
	// DestFile writes the event to the configured log file.
	DestFile
	// DestWebhook posts the event to the configured webhook URL.
	DestWebhook
	// DestAll delivers the event to every registered destination.
	DestAll Destination = DestStdout | DestFile | DestWebhook
)

// Envelope carries an event payload together with routing and tracing
// metadata added by the pipeline before dispatch.
type Envelope[T any] struct {
	// ID is a unique identifier for this delivery attempt.
	ID string
	// CreatedAt is when the envelope was constructed.
	CreatedAt time.Time
	// Dest controls which channels receive this envelope.
	Dest Destination
	// Payload is the original event value.
	Payload T
	// Attempt is the 1-based delivery attempt counter (1 = first try).
	Attempt int
	// Labels contains arbitrary key/value metadata attached by enrichers.
	Labels map[string]string
}

// New constructs an Envelope with CreatedAt set to now and Attempt = 1.
func New[T any](id string, dest Destination, payload T) Envelope[T] {
	return Envelope[T]{
		ID:        id,
		CreatedAt: time.Now(),
		Dest:      dest,
		Payload:   payload,
		Attempt:   1,
		Labels:    make(map[string]string),
	}
}

// WithLabel returns a shallow copy of the envelope with the given label set.
func (e Envelope[T]) WithLabel(key, value string) Envelope[T] {
	next := e
	next.Labels = make(map[string]string, len(e.Labels)+1)
	for k, v := range e.Labels {
		next.Labels[k] = v
	}
	next.Labels[key] = value
	return next
}

// HasDest reports whether all bits in mask are set on the envelope's Dest.
func (e Envelope[T]) HasDest(mask Destination) bool {
	return e.Dest&mask == mask
}

// NextAttempt returns a copy of the envelope with Attempt incremented.
func (e Envelope[T]) NextAttempt() Envelope[T] {
	next := e
	next.Attempt++
	return next
}
