// Package envelope wraps arbitrary payloads with routing metadata,
// delivery tracking, and label annotations.
//
// An Envelope carries a typed value T alongside:
//   - a destination identifier for routing
//   - a set of string labels for filtering and tagging
//   - an attempt counter that increments on each retry
//
// Envelopes are immutable after creation; mutation methods return
// new copies so that concurrent consumers each see a consistent view.
package envelope
