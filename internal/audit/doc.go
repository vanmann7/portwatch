// Package audit provides structured, append-only audit logging for portwatch.
//
// Each port change event (opened, closed, suppressed) is recorded as a
// newline-delimited JSON entry containing a timestamp, event type, port
// number, protocol, and an optional note.
//
// Use New for in-memory or custom writers and NewFileLogger for persistent
// on-disk audit trails.
package audit
