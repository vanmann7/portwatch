// Package correlator groups related port-change events by a shared
// correlation ID derived from the port number and a sliding time window.
//
// Events that affect the same port within the same window bucket share
// a correlation ID, making it easy to detect rapid open/close flapping
// or to attach a single trace identifier to a burst of related alerts.
//
// Usage:
//
//	c := correlator.New(5 * time.Second)
//	id := c.Record(8080, correlator.Opened)
//	events := c.Events(id)
package correlator
