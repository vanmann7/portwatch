// Package dispatcher provides a generic fan-out event dispatcher.
//
// A Dispatcher reads events from a channel and forwards each event to all
// registered Handler functions concurrently. Errors returned by handlers are
// collected and surfaced through a dedicated error channel so the caller can
// decide how to react without blocking event delivery.
//
// Middleware helpers (WithLogging, WithRecover, Chain) allow composing
// cross-cutting concerns around individual handlers without modifying them.
package dispatcher
