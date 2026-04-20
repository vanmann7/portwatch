// Package enricher attaches additional metadata to port events before
// they are dispatched to notification channels or written to the audit log.
//
// Enrichment includes the resolved service name, any user-defined tags,
// the hostname of the machine being monitored, and a UTC timestamp.
package enricher

import (
	"os"
	"time"
)

// Event represents a port-change event with enriched metadata.
type Event struct {
	// Port is the TCP port number that changed state.
	Port int

	// State is either "opened" or "closed".
	State string

	// Service is the well-known service name for the port, or an empty
	// string when the port is not recognised.
	Service string

	// Tags contains any user-defined labels associated with this port.
	Tags []string

	// Host is the hostname of the machine being monitored.
	Host string

	// OccurredAt is the UTC time at which the change was detected.
	OccurredAt time.Time
}

// Resolver maps a port number to a human-readable service name.
type Resolver interface {
	Resolve(port int) string
}

// Tagger maps a port number to a slice of user-defined tag strings.
type Tagger interface {
	Tags(port int) []string
}

// Enricher attaches metadata to raw port-change notifications.
type Enricher struct {
	resolver Resolver
	tagger   Tagger
	hostname string
}

// New returns an Enricher that uses the provided Resolver and Tagger.
// The hostname is read once at construction time; if the lookup fails the
// empty string is stored and no error is returned.
func New(r Resolver, t Tagger) *Enricher {
	host, _ := os.Hostname()
	return &Enricher{
		resolver: r,
		tagger:   t,
		hostname: host,
	}
}

// Enrich constructs an Event for the given port and state string.
// state must be either "opened" or "closed"; any other value is stored
// verbatim without validation.
func (e *Enricher) Enrich(port int, state string) Event {
	var service string
	if e.resolver != nil {
		service = e.resolver.Resolve(port)
	}

	var tags []string
	if e.tagger != nil {
		tags = e.tagger.Tags(port)
	}
	if tags == nil {
		tags = []string{}
	}

	return Event{
		Port:       port,
		State:      state,
		Service:    service,
		Tags:       tags,
		Host:       e.hostname,
		OccurredAt: time.Now().UTC(),
	}
}
